package controller

import (
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/app/response/resource"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

// GetProblemSetProblemUserAnalysis
// 获取基础提交分析，用于指定用户的指定题目的提交集的基础分析
func GetProblemSetProblemUserAnalysis(c echo.Context) error {
	var submissions []models.Submission
	var analysis models.ProblemSetProblemUserAnalysis

	query := base.DB.Model(&models.Submission{}).Preload("User").Preload("Problem").
		Where("problem_set_id=?", c.Param("problem_set_id")).Order("id DESC") // Force order by id desc.
	query = query.Where("problem_id=?", c.Param("problem_id"))
	query = query.Where("user_id=?", c.Param("user_id"))

	total64 := int64(0)
	err := query.Count(&total64).Error
	if err != nil {
		err = errors.Wrap(err, "could not query count of objects")
	}
	if total64 == 0 {
		return c.JSON(http.StatusOK, response.GetProblemSetProblemUserAnalysisResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				Analysis resource.ProblemSetProblemUserAnalysisResource `json:"analysis"`
			}{
				Analysis: resource.ProblemSetProblemUserAnalysisResource{
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
					HighestScore:         analysis.HighestScore,
				}},
		})
	}
	if total64 > 0 {
		err = query.Find(&submissions).Error
		if err != nil {
			err = errors.Wrap(err, "could not query objects")
		}

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

		analysis.TotalWorkTime = time.Duration(analysis.LastSubmissionTime.Sub(analysis.FirstSubmissionTime).Seconds())

		for i := 0; i < len(submissions); i++ {
			if submissions[i].Score > analysis.HighestScore {
				analysis.HighestScore = submissions[i].Score
			}
		}
	}

	return c.JSON(http.StatusOK, response.GetProblemSetProblemUserAnalysisResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Analysis resource.ProblemSetProblemUserAnalysisResource `json:"analysis"`
		}{
			Analysis: resource.ProblemSetProblemUserAnalysisResource{
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
				HighestScore:         analysis.HighestScore,
			}},
	})
}

// GetProblemSetProblemAnalysis
// 获取指定问题集指定题目各学生总的分析，形式为一个结构体数组，每结构体包含一个学生在该问题集的该题目下的提交情况
func GetProblemSetProblemAnalysis(c echo.Context) error {
	var users []models.User
	var submissions []models.Submission
	var generalAnalysis []models.ProblemSetProblemAnalysis

	//先找到所有学生，生成以每个学生为一个结构体的结构体切片
	//1.先找到班级id
	var problemSet models.ProblemSet
	queryClassID := base.DB.Model(&models.ProblemSet{}).Where("id = ?", c.Param("problem_set_id"))
	err := queryClassID.First(&problemSet).Error
	if err != nil {
		err = errors.Wrap(err, "could not query classID")
	}
	classID := problemSet.ClassID
	//2.再找到班级结构体，取出其下的所有学生
	var classStruct models.Class
	queryUsers := base.DB.Model(&models.Class{}).Preload("Students").Where("id=?", classID)
	err = queryUsers.Find(&classStruct).Error
	if err != nil {
		err = errors.Wrap(err, "could not query users")
	}
	users = classStruct.GetStudents()

	//再找到所有学生的所有提交
	query := base.DB.Model(&models.Submission{}).Preload("User").Preload("Problem").
		Where("problem_set_id = ?", c.Param("problem_set_id")).Order("id DESC") // Force order by id desc.
	query = query.Where("problem_id = ?", c.Param("problem_id"))

	total64 := int64(0)
	err = query.Count(&total64).Error
	if err != nil {
		err = errors.Wrap(err, "could not query count of objects")
	}
	//if total64 == 0 {
	//
	//}

	err = query.Find(&submissions).Error
	if err != nil {
		err = errors.Wrap(err, "could not query objects")
	}

	//最后将所有提交填入
	for i := 0; i < len(users); i++ {
		generalAnalysis = append(generalAnalysis, models.ProblemSetProblemAnalysis{})
		generalAnalysis[i].UserID = users[i].ID
		generalAnalysis[i].User = &users[i]
	}

	for i, _ := range generalAnalysis {
		for j, _ := range submissions {
			if generalAnalysis[i].UserID == submissions[j].UserID {
				generalAnalysis[i].Submissions = append(generalAnalysis[i].Submissions, submissions[j])
			}
		}
	}

	for i, _ := range generalAnalysis {
		generalAnalysis[i].TotalSubmissionCount = len(generalAnalysis[i].Submissions)

		t := time.Now()
		for j := 0; j < len(generalAnalysis[i].Submissions); j++ {
			if generalAnalysis[i].Submissions[j].CreatedAt.Before(t) && generalAnalysis[i].Submissions[j].Score == 100 {
				t = generalAnalysis[i].Submissions[j].CreatedAt
			}
		}
		generalAnalysis[i].FirstPassTime = t

		t = time.Now()
		for j := 0; j < len(generalAnalysis[i].Submissions); j++ {
			if generalAnalysis[i].Submissions[j].CreatedAt.Before(t) {
				t = submissions[i].CreatedAt
			}
		}
		generalAnalysis[i].FirstSubmissionTime = t

		//此时t=analysis.FirstSubmissionTime
		for j := 0; j < len(generalAnalysis[i].Submissions); j++ {
			if generalAnalysis[i].Submissions[j].CreatedAt.After(t) {
				t = generalAnalysis[i].Submissions[j].CreatedAt
			}
		}
		generalAnalysis[i].LastSubmissionTime = t

		generalAnalysis[i].TotalWorkTime = time.Duration(generalAnalysis[i].LastSubmissionTime.Sub(generalAnalysis[i].FirstSubmissionTime).Seconds())

		for j := 0; j < len(generalAnalysis[i].Submissions); j++ {
			if generalAnalysis[i].Submissions[j].Score > generalAnalysis[i].HighestScore {
				generalAnalysis[i].HighestScore = generalAnalysis[i].Submissions[j].Score
			}
		}
	}

	return c.JSON(http.StatusOK, response.GetProblemSetProblemAnalysisResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Analysis []resource.ProblemSetProblemAnalysisResource `json:"analysis"`
		}{
			Analysis: resource.GetProblemSetProblemAnalysisResource(generalAnalysis),
		},
	})
}
