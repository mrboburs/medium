package handler

import (
	"fmt"
	"mediumuz/model"
	// "mediumuz/util/error"
	"strconv"

	// "mediumuz/package/repository"

	// "mediumuz/package/repository"
	// "mediumuz/util/error"
	"net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get AllUsers
// @Tags Profile
// @Description get  users
// @ID get-users
// @Accept  json
// @Produce  json
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/account/getUsers [GET]
//@Security ApiKeyAuth
func (handler *Handler) GetAllUsers(ctx *gin.Context) {
	logrus := handler.logrus

	users, err := handler.services.User.GetAllUsers(logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": users,
	})
}

// @Summary Delete User
// @Tags Profile
// @Description delete user
// @ID delete-user
// @Accept  json
// @Produce  json
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/account/delete [DELETE]
//@Security ApiKeyAuth
func (h *Handler) DeleteUser(ctx *gin.Context) {

	logrus := h.logrus
	userId, err := getUserId(ctx, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}

	h.services.User.DeleteUser(strconv.Itoa(userId), logrus)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":      userId,
		"message": "deleted",
	})
}

// @Summary UpdateProfile
// @Tags Profile
// @Description update profile
// @ID update-profile
// @Accept  json
// @Produce  json
// @Param input body model.UserUpdate true "account info"
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/account/update [post]
// @Security ApiKeyAuth
func (handler *Handler) UpdateProfile(ctx *gin.Context) {
	logrus := handler.logrus
	var input model.UserUpdate
	userId, err := getUserId(ctx, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	err = ctx.BindJSON(&input)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	effectedRowsNum, err := handler.services.User.UpdateProfile(userId, input, logrus)
	if effectedRowsNum == 0 {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "User not found", logrus)
		return
	} else if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	} else {

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "updated",
			"data":    input,
		})
	}

}

// @Summary Upload Account Image
// @Description Upload Account Image
// @ID upload-image
// @Tags   Profile
// @Accept       json
// @Produce      json
// @Produce application/octet-stream
// @Produce image/png
// @Produce image/jpeg
// @Produce image/jpg
// @Param file formData file true "file"
// @Accept multipart/form-data
// @Success      200   {object}      model.ResponseSuccess
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router   /api/account/upload-image [PATCH]
//@Security ApiKeyAuth
func (handler *Handler) uploadAccountImage(ctx *gin.Context) {
	logrus := handler.logrus
	userId, err := getUserId(ctx, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	// userId := id
	//  strconv.Itoa(id)
	ctx.Request.ParseMultipartForm(10 << 20)
	file, header, err := ctx.Request.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	// user, err := handler.services.GetUserData(strconv.Itoa(userId), logrus)
	filePath, err := handler.services.UploadAccountImage(file, header, logrus)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	effectedRowsNum, err := handler.services.UpdateAccountImage(userId, filePath, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	if effectedRowsNum == 0 {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "User not found", logrus)
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseSuccess{Message: "Uploaded", Data: filePath})
}

// @Summary GetUserById
// @Description GetUserById
// @ID get-account
// @Tags   Profile
// @Accept       json
// @Produce      json
// @Param id query int false "id"
// @Success      200   {object}      model.ResponseSuccess
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/account/get [GET]
//@Security ApiKeyAuth
func (handler *Handler) GetUserById(ctx *gin.Context) {
	logrus := handler.logrus
	id := ctx.Query("id")
	authID, err := getUserId(ctx, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	userId := strconv.Itoa(authID)

	if id == "" {
		id = userId
	}

	user, err := handler.services.GetUserById(id, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})

}

func (handler *Handler) searchUser(ctx *gin.Context) {

}
