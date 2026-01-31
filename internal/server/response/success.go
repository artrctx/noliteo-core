package response

type SuccessResponse[t any] struct {
	Message string `json:"message"`
	Data    t      `json:"data"`
}
