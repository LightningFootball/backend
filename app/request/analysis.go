package request

// GetSubmissionsBasicAnalysisRequest
// 获取提交基础分析。
// 基础分析为简约版本的分析，获取指定问题、指定用户的历次提交的时间维度分析。
type GetSubmissionsBasicAnalysisRequest struct {
	ProblemId uint `json:"problem_id" form:"problem_id" query:"problem_id"`
	UserId    uint `json:"user_id" form:"user_id" query:"user_id"`
}
