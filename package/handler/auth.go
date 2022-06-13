package handler

import (
	// "fmt"
	"fmt"
	"mediumuz/model"
	"mediumuz/package/repository"
	"mediumuz/package/service"

	// "mediumuz/util/err0r"
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
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
// @Security ApiKeyAuth
func (handler *Handler) signUp(ctx *gin.Context) {
	logrus := handler.logrus
	var input model.User
	err := ctx.BindJSON(&input)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	// logrus.Info("signUp data send for  create user to service")
	id, err := handler.services.Authorization.CreateUser(input, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	token, err := handler.services.Authorization.GenerateToken(input.Email, input.Password, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	verificationCode, err := service.SendCodeToEmail(input.Email, input.FirstName, logrus)
	repository.SetVerificationCode(verificationCode)
	repository.GetVerificationCode()

	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
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
// @Param input body model.SignInInput true "credentials"
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
//@Security ApiKeyAuth
func (handler *Handler) signIn(ctx *gin.Context) {
	fmt.Println(333333)
	logrus := handler.logrus
	var input model.SignInInput

	if err := ctx.BindJSON(&input); err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	token, err := handler.services.Authorization.GenerateToken(input.Email, input.Password, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	id, err := handler.services.Authorization.GetUserID(input.Email, logrus)
	if err != nil {
		logrus.Errorf("errrrr", err)
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
	}
	// error.NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"id":    id,
	})
	// ctx.JSON(http.StatusOK, model.ResponseSignIn{Token: token, Id: user.Id})
}

// @Summary Verify account by email
// @Tags Auth
// @Description verify account
// @ID verify-account
// @Accept  json
// @Produce  json
// @Param input body model.VerificationCode true "verification code"
// @Success 200 {object} model.ResponseSuccess
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/verify [post]
//@Security ApiKeyAuth
func (handler *Handler) ConfirmByEmail(ctx *gin.Context) {
	logrus := handler.logrus
	var CodeInput model.VerificationCode
	userId, err := getUserId(ctx, logrus)
	fmt.Println(userId, "id")
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	// Id := repository.GetAccountId()
	err = ctx.BindJSON(&CodeInput)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	code := repository.GetVerificationCode()
	fmt.Println(code, "cooddd")
	if code == CodeInput.Code {
		ctx.JSON(http.StatusOK, model.ResponseSuccess{Message: "correct verificationCode", Data: true})
		effectedRowsNum, err := handler.services.Authorization.IsVerified(userId, logrus)
		if effectedRowsNum == 0 {
			NewHandlerErrorResponse(ctx, http.StatusBadRequest, "User not found", logrus)
			return
		}
		ctx.JSON(http.StatusOK, model.ResponseSuccess{Message: "DONE", Data: effectedRowsNum})
		if err != nil {
			NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
			return
		}
	} else {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "incorrect code", logrus)
	}

}
