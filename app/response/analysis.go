package response

import "github.com/EduOJ/backend/app/response/resource"

type GetSubmissionsBasicAnalysisResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		BasicAnalysisResponse resource.BasicAnalysis `json:"analysis"`
	} `json:"data"`
}

type GetProblemSetSpecificProblemAnalysis struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		ProblemSetSpecificProblemAnalysisResponse []resource.ProblemSetSpecificProblemAnalysis `json:"problem_set_specific_problem_analysis_response"`
	} `json:"data"`
}
