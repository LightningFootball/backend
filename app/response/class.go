package response

import "github.com/LightningFootball/backend/app/response/resource"

type CreateClassResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type GetClassResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Class `json:"class"`
	} `json:"data"`
}

type GetClassResponseForAdmin struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type UpdateClassResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type AddStudentsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type DeleteStudentsResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type RefreshInviteCodeResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.ClassDetail `json:"class"`
	} `json:"data"`
}

type JoinClassResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
	Data    struct {
		*resource.Class `json:"class"`
	} `json:"data"`
}
