package request

type CreateSubmissionRequest struct {
	Language string `json:"language" form:"language" query:"language" validate:"required"`
	// code(required)
}

type GetSubmissionRequest struct {
}

type GetSubmissionsRequest struct {
	ProblemId uint `json:"problem" form:"problem" query:"problem"`
	UserId    uint `json:"user" form:"user" query:"user"`

	Limit  int `json:"limit" form:"limit" query:"limit" validate:"max=100,min=0"`
	Offset int `json:"offset" form:"offset" query:"offset" validate:"min=0"`
}
