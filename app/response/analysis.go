package response

import "github.com/EduOJ/backend/app/response/resource"

type GetProblemSetProblemUserAnalysisResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Analysis resource.ProblemSetProblemUserAnalysisResource `json:"analysis"`
	} `json:"data"`
}

type GetProblemSetProblemAnalysisResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Analysis []resource.ProblemSetProblemAnalysisResource `json:"analysis"`
	} `json:"data"`
}
