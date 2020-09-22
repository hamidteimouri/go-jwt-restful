package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {

	/* to create response message */
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}

}

func JsonResponse(w http.ResponseWriter, statusCode int, data map[string]interface{}) {
	/* adding status code, then encode and send the response */
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	/* used when we have to give error response */
	if err != nil {
		JsonResponse(w, statusCode, Message(false, err.Error()))
		return
	}
	JsonResponse(w, http.StatusBadRequest, nil)
}
