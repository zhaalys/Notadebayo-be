package task

type CreateTaskRequest struct {
	UserID string `json:"userId" validate:"required"`
	Title  string `json:"title" validate:"required"`
	Desc   string `json:"desc"` // description bersifat optional
	Label  string `json:"label" validate:"required"`
}

type EditTaskRequest struct {
	Title string `json:"title" validate:"required"`
	Desc  string `json:"desc"` // description bersifat optional
	Label string `json:"label" validate:"required"`
}