package response

import "github.com/EduOJ/backend/app/response/resource"

type GetSubmissionsBasicAnalysisResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Analysis resource.Analysis `json:"analysis"`
	} `json:"data"`
}
