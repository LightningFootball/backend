package resource

import (
	"github.com/EduOJ/backend/database/models"

	"time"
)

type BasicAnalysis struct {
	UserID       uint   `json:"user_id"`
	User         *User  `json:"user"`
	ProblemID    uint   `json:"problem_id"`
	ProblemSetID uint   `json:"problem_set_id"` // 0 means not in problem set
	LanguageName string `json:"language_name"`

	TotalSubmissionCount int           `json:"total_submission_count"`
	FirstSubmissionTime  time.Time     `json:"first_submission_time"`
	FirstPassTime        time.Time     `json:"first_pass_time"`
	LastSubmissionTime   time.Time     `json:"last_submission_time"`
	TotalWorkTime        time.Duration `json:"total_work_time"`
	HighestScore         uint          `json:"highest_score"`
}

type ProblemSetSpecificProblemAnalysis struct {
	UserID uint  `json:"user_id"`
	User   *User `json:"user"`

	TotalSubmissionCount int           `json:"total_submission_count"`
	FirstSubmissionTime  time.Time     `json:"first_submission_time"`
	FirstPassTime        time.Time     `json:"first_pass_time"`
	LastSubmissionTime   time.Time     `json:"last_submission_time"`
	TotalWorkTime        time.Duration `json:"total_work_time"`

	HighestScore uint `json:"highest_score"`

	Submissions []Submission `json:"submissions"`
}

func GetProblemSetSpecificProblemAnalysis(model []models.ProblemSetSpecificProblemAnalysis) []ProblemSetSpecificProblemAnalysis {
	var result []ProblemSetSpecificProblemAnalysis
	for i, _ := range model {
		result = append(result, ProblemSetSpecificProblemAnalysis{
			UserID:               model[i].UserID,
			User:                 GetUser(model[i].User),
			TotalSubmissionCount: model[i].TotalSubmissionCount,
			FirstSubmissionTime:  model[i].FirstSubmissionTime,
			FirstPassTime:        model[i].FirstPassTime,
			LastSubmissionTime:   model[i].LastSubmissionTime,
			TotalWorkTime:        model[i].TotalWorkTime,
			HighestScore:         model[i].HighestScore,
			Submissions:          GetSubmissionSlice(model[i].Submissions),
		})
	}
	return result
	//analysis.UserID = model.UserID
	//analysis.User = GetUser(model.User)
	//
	//analysis.TotalSubmissionCount = model.TotalSubmissionCount
	//analysis.FirstSubmissionTime = model.FirstSubmissionTime
	//analysis.FirstPassTime = model.FirstPassTime
	//analysis.LastSubmissionTime = model.LastSubmissionTime
	//analysis.TotalWorkTime = model.TotalWorkTime
	//analysis.HighestScore = model.HighestScore
	//
	//analysis.Submissions=GetSubmissionSlice(model.Submissions)
}
