package common

type Response struct {
	Method  string
	Code    int
	Message string
	Data    string
}

type RequestData struct {
	Version string
	Method  string
	Params  string
}

const (
    HTTP_GET = iota
    HTTP_POST 
    HTTPS_GET 
    HTTPS_POST 
    RPC 
)
