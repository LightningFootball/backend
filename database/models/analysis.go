package models

import "time"

type SubmissionBasicAnalysis struct {
	UserID uint `json:"user_id"`
	//model.User
	User         *User  `json:"user"`
	ProblemID    uint   `json:"problem_id"`
	ProblemSetID uint   `json:"problem_set_id"` // 0 means not in problem set
	LanguageName string `json:"language_name"`

	TotalSubmissionCount int           `json:"total_submission_count"`
	FirstSubmissionTime  time.Time     `json:"first_submission_time"`
	FirstPassTime        time.Time     `json:"first_pass_time"`
	LastSubmissionTime   time.Time     `json:"last_submission_time"`
	TotalWorkTime        time.Duration `json:"total_work_time"`
}
