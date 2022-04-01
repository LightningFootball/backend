package controller

import (
	"github.com/EduOJ/backend/app/request"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/app/response/resource"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

// GetSubmissionsBasicAnalysis
// 获取基础提交分析，用于指定用户的指定题目的提交集的基础分析
func GetSubmissionsBasicAnalysis(c echo.Context) error {
	req := request.GetSubmissionsBasicAnalysisRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}

	query := base.DB.Model(&models.Submission{}).Preload("User").Preload("Problem").
		Where("problem_set_id = 0").Order("id DESC") // Force order by id desc.

	if req.ProblemId != 0 {
		query = query.Where("problem_id = ?", req.ProblemId)
	}
	if req.UserId != 0 {
		query = query.Where("user_id = ?", req.UserId)
	}

	total64 := int64(0)
	err := query.Count(&total64).Error
	if err != nil {
		err = errors.Wrap(err, "could not query count of objects")
	}

	var submissions []models.Submission
	err = query.Find(&submissions).Error
	if err != nil {
		err = errors.Wrap(err, "could not query objects")
	}

	//totalSubmissionCount为总提交次数
	//firstSubmissionTime为首次提交时间,
	//firstPassTime为首次通过时间
	//lastSubmissionTime为最后提交时间
	//totalWorkTime为从第一次提交到最后一次提交时间的总时间
	var analysis models.SubmissionBasicAnalysis
	if total64 > 0 {
		analysis.UserID = submissions[0].UserID
		analysis.User = submissions[0].User
		analysis.ProblemID = submissions[0].ProblemID
		analysis.ProblemSetID = submissions[0].ProblemSetID
		analysis.LanguageName = submissions[0].LanguageName

		analysis.TotalSubmissionCount = len(submissions)

		t := time.Now()
		for i := 0; i < len(submissions); i++ {
			if submissions[i].CreatedAt.Before(t) && submissions[i].Score == 100 {
				t = submissions[i].CreatedAt
			}
		}
		analysis.FirstPassTime = t

		t = time.Now()
		for i := 0; i < len(submissions); i++ {
			if submissions[i].CreatedAt.Before(t) {
				t = submissions[i].CreatedAt
			}
		}
		analysis.FirstSubmissionTime = t

		//此时t=analysis.FirstSubmissionTime
		for i := 0; i < len(submissions); i++ {
			if submissions[i].CreatedAt.After(t) {
				t = submissions[i].CreatedAt
			}
		}
		analysis.LastSubmissionTime = t

		analysis.TotalWorkTime = time.Duration(analysis.FirstSubmissionTime.Sub(analysis.FirstSubmissionTime).Seconds())
	} else {
		return c.JSON(http.StatusOK, response.GetSubmissionsBasicAnalysisResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				Analysis resource.Analysis `json:"analysis"`
			}{
				Analysis: resource.Analysis{
					UserID: analysis.UserID,
					//此处需要对User进行转换，analysis.User为model.User类型，response.User为resource.User类型
					User:                 nil,
					ProblemID:            analysis.ProblemID,
					ProblemSetID:         analysis.ProblemSetID,
					LanguageName:         analysis.LanguageName,
					TotalSubmissionCount: 0,
					FirstSubmissionTime:  analysis.FirstSubmissionTime,
					FirstPassTime:        analysis.FirstPassTime,
					LastSubmissionTime:   analysis.LastSubmissionTime,
					TotalWorkTime:        analysis.TotalWorkTime,
				}},
		})
	}

	return c.JSON(http.StatusOK, response.GetSubmissionsBasicAnalysisResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Analysis resource.Analysis `json:"analysis"`
		}{
			Analysis: resource.Analysis{
				UserID: analysis.UserID,
				//此处需要对User进行转换，analysis.User为model.User类型，response.User为resource.User类型
				User: &resource.User{
					ID:       analysis.User.ID,
					Username: analysis.User.Username,
					Nickname: analysis.User.Nickname,
					Email:    analysis.User.Email,
				},
				ProblemID:            analysis.ProblemID,
				ProblemSetID:         analysis.ProblemSetID,
				LanguageName:         analysis.LanguageName,
				TotalSubmissionCount: analysis.TotalSubmissionCount,
				FirstSubmissionTime:  analysis.FirstSubmissionTime,
				FirstPassTime:        analysis.FirstPassTime,
				LastSubmissionTime:   analysis.LastSubmissionTime,
				TotalWorkTime:        analysis.TotalWorkTime,
			}},
	})
}
