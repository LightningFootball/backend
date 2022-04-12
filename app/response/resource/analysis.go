package resource

import (
	"github.com/EduOJ/backend/database/models"

	"time"
)

type ProblemSetProblemUserAnalysisResource struct {
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

type ProblemSetProblemAnalysisResource struct {
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

func GetProblemSetProblemAnalysisResource(model []models.ProblemSetProblemAnalysis) []ProblemSetProblemAnalysisResource {
	var result []ProblemSetProblemAnalysisResource
	for i, _ := range model {
		result = append(result, ProblemSetProblemAnalysisResource{
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
}
