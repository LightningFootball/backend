package controller

import (
	"context"
	"fmt"
	"github.com/LightningFootball/backend/app/request"
	"github.com/LightningFootball/backend/app/response"
	"github.com/LightningFootball/backend/app/response/resource"
	"github.com/LightningFootball/backend/base"
	"github.com/LightningFootball/backend/base/utils"
	"github.com/LightningFootball/backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"time"
)

func ProblemSetCreateSubmission(c echo.Context) error {
	req := request.ProblemSetCreateSubmissionRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	problemSet := c.Get("problem_set").(*models.ProblemSet)
	findErr := c.Get("find_problem_set_error")
	if findErr != nil {
		if errors.Is(findErr.(error), gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_SET_NOT_FOUND", nil))
		}
		panic(errors.Wrap(findErr.(error), "could not find problem set for getting submission"))
	}
	user := c.Get("user").(models.User)
	var students []models.User
	class := models.Class{}
	if err := base.DB.First(&class, problemSet.ClassID).Error; err != nil {
		panic(errors.Wrap(err, "could not find class for problem set creating submissions"))
	}
	if err := base.DB.Model(&class).Association("Students").Find(&students, "id = ?", user.ID); err != nil {
		panic(errors.Wrap(err, "could not check if student in class for problem set creating submissions"))
	}
	if len(students) == 0 {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}

	var problems []models.Problem
	if err := base.DB.Model(&problemSet).Association("Problems").Find(&problems, "id = ?", c.Param("problem_id")); err != nil {
		panic(errors.Wrap(err, "could not find problem"))
	}
	if len(problems) == 0 {
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	}
	problems[0].LoadTestCases()

	if !utils.Contain(req.Language, problems[0].LanguageAllowed) {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_LANGUAGE", nil))
	}

	language := models.Language{}

	if err := base.DB.First(&language, "name = ?", req.Language).Error; err != nil {
		panic(errors.Wrap(err, "could not find language"))
	}

	file, err := c.FormFile("code")
	if err != nil && err != http.ErrMissingFile && err.Error() != "request Content-Type isn't multipart/form-data" {
		panic(errors.Wrap(err, "could not read file"))
	}
	if file == nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_FILE", nil))
	}

	ext := filepath.Ext(file.Filename)

	if ext != "" {
		ext = ext[1:]
	}

	priority := models.PriorityDefault + 8

	submission := models.Submission{
		UserID:       user.ID,
		User:         &user,
		ProblemID:    problems[0].ID,
		Problem:      &problems[0],
		ProblemSetID: problemSet.ID,
		LanguageName: language.Name,
		Language:     &language,
		FileName:     file.Filename,
		Priority:     priority,
		Judged:       false,
		Score:        0,
		Status:       "PENDING",
		Runs:         make([]models.Run, len(problems[0].TestCases)),
	}
	for i, testCase := range problems[0].TestCases {
		submission.Runs[i] = models.Run{
			UserID:             user.ID,
			ProblemID:          problems[0].ID,
			ProblemSetID:       problemSet.ID,
			TestCaseID:         testCase.ID,
			Sample:             testCase.Sample,
			Priority:           priority,
			Judged:             false,
			Status:             "PENDING",
			MemoryUsed:         0,
			TimeUsed:           0,
			OutputStrippedHash: "",
		}
	}
	utils.PanicIfDBError(base.DB.Create(&submission), "could not create submission and runs")

	utils.MustPutObject(file, c.Request().Context(), "submissions", fmt.Sprintf("%d/code", submission.ID))

	if !inTest {
		base.Redis.Publish(context.Background(), "runs", nil)
	}

	return c.JSON(http.StatusCreated, response.ProblemSetCreateSubmissionResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.SubmissionDetail `json:"submission"`
		}{
			resource.GetSubmissionDetail(&submission),
		},
	})
}

func ProblemSetGetSubmission(c echo.Context) error {
	startedAt := time.Now()
	poll := false
	if c.QueryParam("poll") == "1" {
		poll = true
	}

	problemSet := c.Get("problem_set")
	if problemSet != nil {
		err := c.Get("find_problem_set_error")
		if err != nil {
			if errors.Is(err.(error), gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_SET_NOT_FOUND", nil))
			}
			panic(errors.Wrap(err.(error), "could not find problem set for getting submission"))
		}
	}

	user := c.Get("user").(models.User)
	submission := models.Submission{}
	if err := base.DB.Preload("Problem").Preload("User").
		First(&submission, "problem_set_id = ? and id = ?", c.Param("problem_set_id"), c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find problem"))
		}
	}

	// If problem set is empty here, the user is considered to have permission read_answers(because of
	// the short-circuit in middleware HasPermission), and the submission info is returned directly
	if user.ID != submission.UserID && problemSet != nil {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}
	if !(submission.UpdatedAt.Before(startedAt) && poll) {
		submission.LoadRuns()
		return c.JSON(http.StatusOK, response.ProblemSetGetSubmissionResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.SubmissionDetail `json:"submission"`
			}{
				resource.GetSubmissionDetail(&submission),
			},
		})
	}
	timeoutChan := time.After(viper.GetDuration("polling_timeout"))
	timeout := false
	sub := base.Redis.Subscribe(c.Request().Context(), fmt.Sprintf("submission_update:%d", submission.ID))
	defer sub.Close()
	for {
		select {
		case <-sub.Channel():
			if err := base.DB.Preload("Problem").Preload("User").First(&submission, c.Param("id")).Error; err != nil {
				panic(errors.Wrap(err, "could not find submission"))
			}
			break
		case <-c.Request().Context().Done():
			// context cancelled
			return nil
		case <-timeoutChan:
			timeout = true
		}
		if timeout {
			break
		}
		if submission.UpdatedAt.After(startedAt) {
			break
		}
	}
	submission.LoadRuns()
	return c.JSON(http.StatusOK, response.ProblemSetGetSubmissionResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.SubmissionDetail `json:"submission"`
		}{
			resource.GetSubmissionDetail(&submission),
		},
	})
}

func ProblemSetGetSubmissions(c echo.Context) error {
	req := request.ProblemSetGetSubmissionsRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	problemSet := c.Get("problem_set")
	if problemSet != nil {
		err := c.Get("find_problem_set_error")
		if err != nil {
			if errors.Is(err.(error), gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_SET_NOT_FOUND", nil))
			}
			panic(errors.Wrap(err.(error), "could not find problem set for getting submission"))
		}
		user := c.Get("user").(models.User)
		var students []models.User
		class := models.Class{}
		if err := base.DB.First(&class, problemSet.(*models.ProblemSet).ClassID).Error; err != nil {
			panic(errors.Wrap(err, "could not find class for problem set getting submissions"))
		}
		if err := base.DB.Model(&class).Association("Students").Find(&students, "id = ?", user.ID); err != nil {
			panic(errors.Wrap(err, "could not check if student in class for problem set getting submissions"))
		}
		if len(students) == 0 {
			return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
		}
	}
	// If problem set is empty here, the user is considered to have permission read_answers(because of
	// the short-circuit in middleware HasPermission), and the submission info is returned directly

	query := base.DB.Model(&models.Submission{}).Preload("User").Preload("Problem").
		Where("problem_set_id = ?", c.Param("problem_set_id")).Order("id DESC") // Force order by id desc.

	if req.ProblemId != 0 {
		query = query.Where("problem_id = ?", req.ProblemId)
	}
	if req.UserId != 0 {
		query = query.Where("user_id = ?", req.UserId)
	}

	var submissions []models.Submission
	total, prevUrl, nextUrl, err := utils.Paginator(query, req.Limit, req.Offset, c.Request().URL, &submissions)
	if err != nil {
		if herr, ok := err.(utils.HttpError); ok {
			return herr.Response(c)
		}
		panic(err)
	}
	return c.JSON(http.StatusOK, response.ProblemSetGetSubmissionsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Submissions []resource.Submission `json:"submissions"`
			Total       int                   `json:"total"`
			Count       int                   `json:"count"`
			Offset      int                   `json:"offset"`
			Prev        *string               `json:"prev"`
			Next        *string               `json:"next"`
		}{
			Submissions: resource.GetSubmissionSlice(submissions),
			Total:       total,
			Count:       len(submissions),
			Offset:      req.Offset,
			Prev:        prevUrl,
			Next:        nextUrl,
		},
	})
}

func ProblemSetGetSubmissionCode(c echo.Context) error {
	problemSet := c.Get("problem_set")
	if problemSet != nil {
		err := c.Get("find_problem_set_error")
		if err != nil {
			if errors.Is(err.(error), gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_SET_NOT_FOUND", nil))
			}
			panic(errors.Wrap(err.(error), "could not find problem set for getting submission"))
		}
	}

	user := c.Get("user").(models.User)
	submission := models.Submission{}
	if err := base.DB.First(&submission, "problem_set_id = ? and id = ?",
		c.Param("problem_set_id"), c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find submission for getting submission code"))
		}
	}

	// If problem set is empty here, the user is considered to have permission read_answers(because of
	// the short-circuit in middleware HasPermission), and the submission info is returned directly
	if user.ID != submission.UserID && problemSet != nil {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}

	presignedUrl, err := utils.GetPresignedURL("submissions", fmt.Sprintf("%d/code", submission.ID), submission.FileName)
	if err != nil {
		panic(errors.Wrap(err, "could not get presigned url"))
	}
	return c.Redirect(http.StatusFound, presignedUrl)
}

func ProblemSetGetRunCompilerOutput(c echo.Context) error {
	problemSet := c.Get("problem_set")
	if problemSet != nil {
		err := c.Get("find_problem_set_error")
		if err != nil {
			if errors.Is(err.(error), gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_SET_NOT_FOUND", nil))
			}
			panic(errors.Wrap(err.(error), "could not find problem set for getting submission"))
		}
	}

	user := c.Get("user").(models.User)
	run := models.Run{}
	if err := base.DB.Preload("Submission").First(&run, "problem_set_id = ? and submission_id = ? and id = ?",
		c.Param("problem_set_id"), c.Param("submission_id"), c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find submission for getting submission code"))
		}
	}

	// If problem set is empty here, the user is considered to have permission read_answers(because of
	// the short-circuit in middleware HasPermission), and the submission info is returned directly
	if user.ID != run.UserID && problemSet != nil {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}

	if run.Status == "PENDING" || run.Status == "JUDGEMENT_FAILED" || run.Status == "NO_COMMENT" {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("JUDGEMENT_UNFINISHED", nil))
	}

	presignedUrl, err := utils.GetPresignedURL("submissions", fmt.Sprintf("%d/run/%d/compiler_output", run.Submission.ID, run.ID),
		"compiler_output.txt")
	if err != nil {
		panic(errors.Wrap(err, "could not get presigned url"))
	}
	return c.Redirect(http.StatusFound, presignedUrl)
}

func ProblemSetGetRunOutput(c echo.Context) error {
	problemSet := c.Get("problem_set")
	if problemSet != nil {
		err := c.Get("find_problem_set_error")
		if err != nil {
			if errors.Is(err.(error), gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_SET_NOT_FOUND", nil))
			}
			panic(errors.Wrap(err.(error), "could not find problem set for getting submission"))
		}
	}

	user := c.Get("user").(models.User)
	run := models.Run{}
	if err := base.DB.Preload("TestCase").Preload("Submission").
		First(&run, "problem_set_id = ? and submission_id = ? and id = ?",
			c.Param("problem_set_id"), c.Param("submission_id"), c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find submission for getting submission code"))
		}
	}

	// If problem set is empty here, the user is considered to have permission read_answers(because of
	// the short-circuit in middleware HasPermission), and the submission info is returned directly
	if (user.ID != run.UserID || !run.Sample) && problemSet != nil {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}

	if run.Status == "PENDING" || run.Status == "JUDGEMENT_FAILED" || run.Status == "NO_COMMENT" {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("JUDGEMENT_UNFINISHED", nil))
	}

	presignedUrl, err := utils.GetPresignedURL("submissions", fmt.Sprintf("%d/run/%d/output", run.Submission.ID, run.ID),
		run.TestCase.OutputFileName)
	if err != nil {
		panic(errors.Wrap(err, "could not get presigned url"))
	}
	return c.Redirect(http.StatusFound, presignedUrl)
}

func ProblemSetGetRunInput(c echo.Context) error {
	problemSet := c.Get("problem_set")
	if problemSet != nil {
		err := c.Get("find_problem_set_error")
		if err != nil {
			if errors.Is(err.(error), gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_SET_NOT_FOUND", nil))
			}
			panic(errors.Wrap(err.(error), "could not find problem set for getting submission"))
		}
	}

	user := c.Get("user").(models.User)
	run := models.Run{}
	if err := base.DB.Preload("Submission").Preload("TestCase").
		First(&run, "problem_set_id = ? and submission_id = ? and id = ?",
			c.Param("problem_set_id"), c.Param("submission_id"), c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find submission for getting submission code"))
		}
	}

	// If problem set is empty here, the user is considered to have permission read_answers(because of
	// the short-circuit in middleware HasPermission), and the submission info is returned directly
	if (user.ID != run.UserID || !run.Sample) && problemSet != nil {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}

	presignedUrl, err := utils.GetPresignedURL("problems", fmt.Sprintf("%d/input/%d.in", run.ProblemID, run.TestCaseID),
		run.TestCase.InputFileName)
	if err != nil {
		panic(errors.Wrap(err, "could not get presigned url"))
	}
	return c.Redirect(http.StatusFound, presignedUrl)
}

func ProblemSetGetRunComparerOutput(c echo.Context) error {
	problemSet := c.Get("problem_set")
	if problemSet != nil {
		err := c.Get("find_problem_set_error")
		if err != nil {
			if errors.Is(err.(error), gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, response.ErrorResp("PROBLEM_SET_NOT_FOUND", nil))
			}
			panic(errors.Wrap(err.(error), "could not find problem set for getting submission"))
		}
	}

	var ss []models.Submission
	base.DB.Preload("Runs").Find(&ss)

	user := c.Get("user").(models.User)
	run := models.Run{}
	if err := base.DB.Preload("Submission").First(&run, "problem_set_id = ? and submission_id = ? and id = ?",
		c.Param("problem_set_id"), c.Param("submission_id"), c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find submission for getting submission code"))
		}
	}

	// If problem set is empty here, the user is considered to have permission read_answers(because of
	// the short-circuit in middleware HasPermission), and the submission info is returned directly
	if (user.ID != run.UserID || !run.Sample) && problemSet != nil {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}

	if run.Status == "PENDING" || run.Status == "JUDGEMENT_FAILED" || run.Status == "NO_COMMENT" {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("JUDGEMENT_UNFINISHED", nil))
	}

	presignedUrl, err := utils.GetPresignedURL("submissions", fmt.Sprintf("%d/run/%d/comparer_output", run.Submission.ID, run.ID),
		"comparer_output.txt")
	if err != nil {
		panic(errors.Wrap(err, "could not get presigned url"))
	}
	return c.Redirect(http.StatusFound, presignedUrl)
}
