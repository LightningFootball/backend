package request

// GetProblemSetProblemUserAnalysisRequest
// 获取指定问题集、制定问题、指定用户的历次提交的时间维度分析。
// 班级-问题集-问题-用户
type GetProblemSetProblemUserAnalysisRequest struct {
	ProblemSetId int64 `json:"problem_set_id" form:"problem_set_id" query:"problem_set_id"` // 问题集ID
	ProblemId    uint  `json:"problem_id" form:"problem_id" query:"problem_id"`             // 问题ID
	UserId       uint  `json:"user_id" form:"user_id" query:"user_id"`                      // 用户ID
}

// GetProblemSetProblemAnalysisRequest
// 获取指定问题集、制定问题、全部用户的历次提交的分析。
type GetProblemSetProblemAnalysisRequest struct {
	ProblemSetId int64 `json:"problem_set_id" form:"problem_set_id" query:"problem_set_id"` // 问题集ID
	ProblemId    uint  `json:"problem_id" form:"problem_id" query:"problem_id"`             // 问题ID
}
