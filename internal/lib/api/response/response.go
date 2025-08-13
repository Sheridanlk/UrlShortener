package response

type Response struct {
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
}

const (
	SatusOK = "OK"
	SatusError = "Error"
)

func OK() Response {
	return Response{
		Status: SatusOK,
	}
}
	
func Error(msg string) Response {
	return Response{
		Status: SatusError,
		Error: msg,
	}
}