package response

import (
	"DotConnect/utils/json"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, status int, message string, err error) error {
	response := map[string]any{}
	errDetails := any(nil)
	if err != nil {
		errDetails = err.Error()
	}

	if status < 500 {
		response["message"] = message
		response["details"] = errDetails
		response["data"] = map[string]any{}
	} else {
		response["success"] = false
		response["message"] = message
		response["error"] = errDetails
		response["data"] = map[string]any{}
	}

	return json.ConvertTOJSON(w, status, response)
}

func SucessResponse(w http.ResponseWriter,status int,message string,data any) error{

	response:=map[string]any{
		"Success": true,
		"Message": message,
	}
	if data != nil {
		response["data"] = data
	}
	if status == http.StatusCreated {
		response["error"] = map[string]any{}
	}

	return json.ConvertTOJSON(w,status,response)
}