package response

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

const (
	SUCCESS = iota
	WRONGTOKEN
	BADREQUEST
	NOTFOUND
	INTERNALERROR
)
