package schemas

type Response struct {
	Ok    bool `json:"ok"`
	Data  *any `json:"data"`
	Error *any `json:"error"`
}

func NewResponse(data any, err any) *Response {
	return &Response{
		Ok:    err == nil,
		Data:  &data,
		Error: &err,
	}
}
