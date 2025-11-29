package responses

type APIErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	HTTP    int    `json:"http"`
}
