package handlers

type ArticleRequest struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type ArticleResponse struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

type RespError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewRespError(err error, code int) RespError {
	return RespError{
		Message: err.Error(),
		Code:    code}

}
