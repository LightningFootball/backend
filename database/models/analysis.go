package models

import "time"

type SubmissionBasicAnalysis struct {
	UserID uint `json:"user_id"`
	//model.User
	User         *User  `json:"user"`
	ProblemID    uint   `json:"problem_id"`
	ProblemSetID uint   `json:"problem_set_id"` // 0 means not in problem set
	LanguageName string `json:"language_name"`

	//totalSubmissionCount为总提交次数
	//firstSubmissionTime为首次提交时间,
	//firstPassTime为首次通过时间
	//lastSubmissionTime为最后提交时间
	//totalWorkTime为从第一次提交到最后一次提交时间的总时间
	//highestScore为最高分
	TotalSubmissionCount int           `json:"total_submission_count"`
	FirstSubmissionTime  time.Time     `json:"first_submission_time"`
	FirstPassTime        time.Time     `json:"first_pass_time"`
	LastSubmissionTime   time.Time     `json:"last_submission_time"`
	TotalWorkTime        time.Duration `json:"total_work_time"`

	HighestScore uint `json:"highest_score"`
}

type ProblemSetSpecificProblemAnalysis struct {
	UserID uint `json:"user_id"`
	//model.User
	User *User `json:"user"`

	TotalSubmissionCount int           `json:"total_submission_count"`
	FirstSubmissionTime  time.Time     `json:"first_submission_time"`
	FirstPassTime        time.Time     `json:"first_pass_time"`
	LastSubmissionTime   time.Time     `json:"last_submission_time"`
	TotalWorkTime        time.Duration `json:"total_work_time"`

	HighestScore uint `json:"highest_score"`

	Submissions []Submission `json:"submissions"`
}
