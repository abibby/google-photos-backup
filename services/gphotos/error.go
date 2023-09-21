package gphotos

type GError struct {
	Err struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *GError) Error() string {
	return e.Err.Message
}
