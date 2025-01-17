package response

import (
	"github.com/LightningFootball/backend/app/response/resource"
)

type GetProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Problem `json:"problem"`
	} `json:"data"`
}

type GetProblemsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Problems []resource.ProblemSummary `json:"problems"`
		Total    int                       `json:"total"`
		Count    int                       `json:"count"`
		Offset   int                       `json:"offset"`
		Prev     *string                   `json:"prev"`
		Next     *string                   `json:"next"`
	} `json:"data"`
}

type CreateProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemForAdmin `json:"problem"`
	} `json:"data"`
}

type GetProblemResponseForAdmin struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemForAdmin `json:"problem"`
	} `json:"data"`
}

type GetProblemsResponseForAdmin struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		Problems []resource.ProblemSummaryForAdmin `json:"problems"`
		Total    int                               `json:"total"`
		Count    int                               `json:"count"`
		Offset   int                               `json:"offset"`
		Prev     *string                           `json:"prev"`
		Next     *string                           `json:"next"`
	} `json:"data"`
}

type UpdateProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ProblemForAdmin `json:"problem"`
	} `json:"data"`
}

type CreateTestCaseResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.TestCaseForAdmin `json:"test_case"`
	} `json:"data"`
}

type UpdateTestCaseResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.TestCaseForAdmin `json:"test_case"`
	} `json:"data"`
}

type GetRandomProblemResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Problem `json:"problem"`
	} `json:"data"`
}
