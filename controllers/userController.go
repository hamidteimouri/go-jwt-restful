package controllers

import (
	"encoding/json"
	"errors"
	"github.com/hamidteimouri/go-jwt-restful/models"
	"github.com/hamidteimouri/go-jwt-restful/utils"
	"net/http"
)

var CreateUser = func(writer http.ResponseWriter, request *http.Request) {
	/* decode request body to user struct */
	user := &models.Token{}
	err := json.NewDecoder(request.Body).Decode(user)
	if err != nil {
		utils.ERROR(writer, http.StatusUnprocessableEntity, errors.New(err.Error()))
		return
	}

	/* calling create user function */
	//response , err := user.Create()
}
