package handler

import (
	// "fmt"
	"mediumuz/model"
	"mediumuz/package/repository"
	"mediumuz/package/service"
	"mediumuz/util/error"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary SignUp
// @Tags Auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body model.User true "account info"
// @Success 200 {object} model.ResponseSignUp
// @Failure 400,404 {object} error.errorResponse
// @Failure 409 {object} error.errorResponse
// @Failure 500 {object} error.errorResponse
// @Failure default {object} error.errorResponse
// @Router /auth/sign-up [post]
// @Security ApiKeyAuth
func (handler *Handler) signUp(ctx *gin.Context) {
	logrus := handler.logrus
	var input model.User
	err := ctx.BindJSON(&input)
	if err != nil {
		error.NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	// logrus.Info("signUp data send for  create user to service")
	id, err := handler.services.Authorization.CreateUser(input, logrus)
	if err != nil {
		error.NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	token, err := handler.services.Authorization.GenerateToken(input.Email, input.Password, logrus)
	if err != nil {
		error.NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	verificationCode, err := service.SendCodeToEmail(input.Email, input.UserName, logrus)
	repository.SetVerificationCode(verificationCode)
	repository.GetVerificationCode()

	if err != nil {
		error.NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseSignUp{Id: id, Token: token})
}

// @Summary SignIn
// @Tags Auth
// @Description login account
// @ID login-account
// @Accept  json
// @Produce  json
// @Param input body model.SignInInput true "accountin infor"
// @Success 200 {object} model.ResponseSignIn
// @Failure 400,404 {object} error.errorResponse
// @Failure 409 {object} error.errorResponse
// @Failure 500 {object} error.errorResponse
// @Failure default {object} error.errorResponse
// @Router /auth/sign-in [post]
//@Security ApiKeyAuth
func (handler *Handler) signIn(ctx *gin.Context) {
	logrus := handler.logrus
	var input model.SignInInput

	if err := ctx.BindJSON(&input); err != nil {
		error.NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	token, err := handler.services.Authorization.GenerateToken(input.Email, input.Password, logrus)
	if err != nil {
		error.NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseSignIn{Token: token, Email: input.Email})
}

// @Summary Verify email
// @Tags Auth
// @Description verify account
// @ID verify-account
// @Accept  json
// @Produce  json
// @Param input body model.VerificationCode true "account info"
// @Success 200 {object} model.ResponseSuccess
// @Failure 400,404 {object} error.errorResponse
// @Failure 409 {object} error.errorResponse
// @Failure 500 {object} error.errorResponse
// @Failure default {object} error.errorResponse
// @Router /auth/verify [post]
//@Security ApiKeyAuth
func (handler *Handler) ConfirmEmail(ctx *gin.Context) {
	logrus := handler.logrus
	var CodeInput model.VerificationCode
	userId, err := getUserId(ctx, logrus)
	if err != nil {
		error.NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	// Id := repository.GetAccountId()
	err = ctx.BindJSON(&CodeInput)
	if err != nil {
		error.NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	code := repository.GetVerificationCode()
	if code == CodeInput.Code {
		ctx.JSON(http.StatusOK, model.ResponseSuccess{Message: "correct verificationCode", Data: true})
		effectedRowsNum, err := handler.services.Authorization.VerifyEmail(userId, logrus)
		if effectedRowsNum == 0 {
			error.NewHandlerErrorResponse(ctx, http.StatusBadRequest, "User not found", logrus)
			return
		}
		ctx.JSON(http.StatusOK, model.ResponseSuccess{Message: "DONE", Data: effectedRowsNum})
		if err != nil {
			error.NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
			return
		}
	} else {
		error.NewHandlerErrorResponse(ctx, http.StatusBadRequest, "incorrect code", logrus)
	}

}
