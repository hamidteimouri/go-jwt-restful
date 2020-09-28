package controllers

import (
	"encoding/json"
	"github.com/hamidteimouri/go-jwt-restful/models"
	"github.com/hamidteimouri/go-jwt-restful/utils"
	"io/ioutil"
	"net/http"
)

var CreateUser = func(writer http.ResponseWriter, r *http.Request) {
	/* decoding request body into User struct */
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	/* calling create user func */
	response, err := user.Create()
	if err != nil {
		/* return error response if any error or malformed request path */
		utils.ERROR(writer, http.StatusUnauthorized, err)
		return
	}

	/* sending OK and the user information as response */
	utils.JsonResponse(writer, http.StatusCreated, response)
}

var SignInUser = func(writer http.ResponseWriter, request *http.Request) {
	/* decode request's body into user struct */
	user := &models.User{}
	err := json.NewDecoder(request.Body).Decode(user)
	if err != nil {
		utils.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	/* calling SignIn function in user struct by passing only email and password */
	resp, err := models.SignIn(user.Email, user.Password)

	if err != nil {
		utils.ERROR(writer, http.StatusUnauthorized, err)
		return
	}

	/* sending OK and the user information as response */
	utils.JsonResponse(writer, http.StatusOK, resp)
}

var DeleteUser = func(writer http.ResponseWriter, request *http.Request) {
	/* getting userId from the request context */
	user := request.Context().Value("user").(uint)

	/* calling delete from models module and passing user id */
	resp, err := models.DeleteUser(user)

	if err != nil {
		utils.ERROR(writer, http.StatusBadRequest, err)
		return
	}

	/* sending OK and the message as response */
	utils.JsonResponse(writer, http.StatusOK, resp)
}

var UpdateUserPassword = func(writer http.ResponseWriter, request *http.Request) {
	/* getting the userid from the request context */
	id := request.Context().Value("user").(uint)

	/* decoding the request body and make is as map of interface */
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	err = json.Unmarshal(body, &keyVal)

	if err != nil {
		utils.ERROR(writer, http.StatusUnprocessableEntity, err)
		return
	}

	/* set new password and check of error */
	resp, err := models.UpdatePassword(id, keyVal["new password"])
	if err != nil {
		utils.ERROR(writer, http.StatusBadRequest, err)
	}

	/* sending OK and the message as response */
	utils.JsonResponse(writer, http.StatusOK, resp)
}
