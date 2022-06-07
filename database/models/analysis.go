package models

import "time"

type ProblemSetProblemUserAnalysis struct {
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

type ProblemSetProblemAnalysis struct {
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

// 作业报告
// todo:后续可以把report几个struct改成map，应该会提高report那部分数据填充的效率
type ProblemSetStudentReport struct {
	UserID   uint            `json:"user_id"`
	UserName string          `json:"user_name"`
	NickName string          `json:"nick_name"`
	Problem  []SimpleProblem `json:"simple_problem"`
}

type SimpleProblem struct {
	ProblemID    uint `json:"problem_id"`
	HighestScore uint `json:"highest_score"`
}
