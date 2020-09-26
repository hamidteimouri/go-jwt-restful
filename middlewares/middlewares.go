package middlewares

import "net/http"

func SetJsonMiddleware(next http.HandlerFunc) http.HandlerFunc {

	/* to set response as json response */
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		next(writer, request)
	}

}
