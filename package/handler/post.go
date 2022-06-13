package handler

import (
	"fmt"
	"mediumuz/model"
	// "mediumuz/util/logrus"

	// "mediumuz/util/error"

	// "mediumuz/util/logrus"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResponseSuccessForCreatePost struct {
	Message string
	PostId  int
}

// @Summary Get Comments
// @Tags Post
// @Description get commits
// @ID get-comments
// @Accept  json
// @Produce  json
// @Param        offset   query  int     false "Offset "
// @Param        limit   query  int     false "Limit "
// @Param        postID   query  int     true "postID "
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /post/get-comments [GET]
func (handler *Handler) getComments(ctx *gin.Context) {
	logrus := handler.logrus
	var pagination model.Pagination
	postIDQuery := ctx.Query("postID")
	if postIDQuery == "" {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "Offset is empty", logrus)
		return
	}

	postID, err := strconv.Atoi(postIDQuery)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	offsetQuery := ctx.DefaultQuery("offset", "0")
	if offsetQuery == "" {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "Offset is empty", logrus)
		return
	}

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	limitQuery := ctx.DefaultQuery("limit", "10")

	if limitQuery == "" {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "Limit is empty", logrus)
		return
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	pagination.Offset = offset
	pagination.Limit = limit
	result, err := handler.services.GetCommentsPost(postID, pagination, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseSuccess{Message: "Post Commits", Data: result})

}

// @Summary Commit  Post
// @Tags Post
// @Description Commit post by user
// @ID commit-post-id
// @Accept  json
// @Produce  json
// @Param input body model.CommentPost true "commit info"
// @Success 200 {object} model.ResponseSuccess
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/post/comment [POST]
//@Security ApiKeyAuth
func (handler *Handler) commentPost(ctx *gin.Context) {
	logrus := handler.logrus

	var input model.CommentPost
	err := ctx.BindJSON(&input)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	userID, err := getUserId(ctx, logrus)
	if err != nil {
		if err != nil {
			NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
			return
		}
		return
	}

	input.ReaderID = userID
	commentID, err := handler.services.CommentPost(input, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	if commentID == 0 {
		if err != nil {
			NewHandlerErrorResponse(ctx, http.StatusBadRequest, "DO NOT WORK", logrus)
			return
		}
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{

		"msg": "successfully comments posted",
		"id":  commentID,
	})

}

// @Summary View  Post By ID
// @Tags Post
// @Description View post by id
// @ID view-post-id
// @Accept  json
// @Produce  json
// @Param        id   query  int     true "Param ID"
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/post/view [GET]
//@Security ApiKeyAuth
func (handler *Handler) viewPost(ctx *gin.Context) {
	logrus := handler.logrus
	userID, err := getUserId(ctx, logrus)
	if err != nil {
		if err != nil {
			NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
			return
		}
		return
	}

	paramID := ctx.Query("id")

	if paramID == "" {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "Param is empty", logrus)
		return
	}

	postID, err := strconv.Atoi(paramID)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	result, err := handler.services.ViewPost(userID, postID, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	if result == 0 {
		if err != nil {
			NewHandlerErrorResponse(ctx, http.StatusBadRequest, "DO NOT WORK", logrus)
			return
		}

	}
	ctx.JSON(http.StatusOK, map[string]interface{}{

		"msg": "successfully getted",
		"id":  result,
	})
}

// @Summary ClickLike
// @Tags Post
// @Description update like-count
// @ID counting-like
// @Accept  json
// @Produce  json
// @Param        id   path  int     true "Param ID"
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/post/count-like/{id} [get]
// @Security ApiKeyAuth
func (handler *Handler) ClickLike(ctx *gin.Context) {
	logrus := handler.logrus
	userId, err := getUserId(ctx, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	paramID := ctx.Param("id")
	if paramID == "" {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "Param is empty", logrus)
		return
	}
	id, err := strconv.Atoi(paramID)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	err = handler.services.Post.ClickLike(userId, id, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{

		"clickingLike": "clicked",
	})
}

// @Summary UpdatePost
// @Tags Post
// @Description update post
// @ID update-post
// @Accept  json
// @Produce  json
// @Param        id   path  int     true "Param ID"
// @Param input body model.Post true "account info"
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/post/update/{id} [post]
// @Security ApiKeyAuth
func (handler *Handler) UpdatePost(ctx *gin.Context) {
	logrus := handler.logrus
	var input model.Post
	paramID := ctx.Param("id")
	if paramID == "" {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "Param is empty", logrus)
		return
	}
	id, err := strconv.Atoi(paramID)
	err = ctx.BindJSON(&input)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	err = handler.services.Post.UpdatePost(id, input, logrus)
	if err != nil {

		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "not updated", logrus)
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "updated",
		"data":    input,
	})

}

// @Summary Get  Posts
// @Tags Post
// @Description Get  Posts
// @ID get-posts
// @Accept  json
// @Produce  json
// @Param        offset   query  int     false "Offset "
// @Param        limit   query  int     false "Limit "
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /post/get-all [GET]
func (handler *Handler) GetAllPosts(ctx *gin.Context) {
	logrus := handler.logrus
	var pagination model.Pagination
	offsetQuery := ctx.DefaultQuery("offset", "0")
	if offsetQuery == "" {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "Offset is empty", logrus)
		return
	}

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	limitQuery := ctx.DefaultQuery("limit", "10")

	if limitQuery == "" {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "Limit is empty", logrus)
		return
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	pagination.Offset = offset
	pagination.Limit = limit
	result, err := handler.services.GetAllPosts(pagination, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseSuccess{Message: "User Posts Body", Data: result})

}

// @Summary Create Post
// @Tags Post
// @Description create post
// @ID create-post
// @Accept  json
// @Produce  json
// @Param input body model.Post true "post info"
// @Success 200 {object} ResponseSuccessForCreatePost
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/post/create [post]
//@Security ApiKeyAuth
func (handler *Handler) createPost(ctx *gin.Context) {
	logrus := handler.logrus
	var input model.Post

	err := ctx.BindJSON(&input)
	if err != nil {

		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	userId, err := getUserId(ctx, logrus)
	if err != nil {

		return
	}

	postId, err := handler.services.CreatePost(userId, input, logrus)
	if err != nil {

		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	ctx.JSON(http.StatusOK, ResponseSuccessForCreatePost{PostId: postId, Message: "DONE"})
}

// @Summary Get  Post By ID
// @Tags Post
// @Description get post by id
// @ID get-post-id
// @Accept  json
// @Produce  json
// @Param        id   path  int     true "Param ID"
// @Success 200 {object} model.ResponseSuccess
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /post/get/{id} [GET]
func (handler *Handler) getPostID(ctx *gin.Context) {
	logrus := handler.logrus
	paramID := ctx.Param("id")
	if paramID == "" {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "Param is empty", logrus)
		return
	}
	id, err := strconv.Atoi(paramID)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	check, err := handler.services.Post.CheckPostId(id, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	if check == 0 {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "ID Not Fount", logrus)
		return
	}
	resp, err := handler.services.GetPostById(id, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}
	ctx.JSON(http.StatusOK, model.ResponseSuccess{Data: resp, Message: "DONE"})

}

// @Summary Delete  Post By ID
// @Tags Post
// @Description delete post by id
// @ID delete-post-id
// @Accept  json
// @Produce  json
// @Param        id   path  int     true "Param ID"
// @Success 200 {object} model.ResponseSuccess
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/post/delete/{id} [DELETE]
//@Security ApiKeyAuth
func (handler *Handler) PostDelete(ctx *gin.Context) {
	logrus := handler.logrus
	userId, err := getUserId(ctx, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "invalid id param", logrus)
		return
	}
	err = handler.services.Post.PostDelete(userId, id)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":     id,
		"status": "deleted",
	})

}

// @Summary Upload Post Image
// @Description Upload Post Image
// @ID uploadImgPost
// @Tags   Post
// @Accept       json
// @Produce      json
// @Produce application/octet-stream
// @Produce image/png
// @Produce image/jpeg
// @Produce image/jpg
// @Param        id   path  int     true "Param ID"
// @Param file formData file true "file"
// @Accept multipart/form-data
// @Success      200   {object} model.ResponseSuccess
// @Failure 400,404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router   /api/post/upload/{id} [PATCH]
//@Security ApiKeyAuth
func (handler *Handler) uploadPostImg(ctx *gin.Context) {
	fmt.Println("yes")
	logrus := handler.logrus
	userId, err := getUserId(ctx, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusInternalServerError, err.Error(), logrus)
		return
	}
	fmt.Println("yes")
	ctx.Request.ParseMultipartForm(10 << 20)
	file, header, err := ctx.Request.FormFile("file")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	postId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, "invalid id param", logrus)
		return
	}
	post, err := handler.services.Post.GetPostById(postId, logrus)
	filePath, err := handler.services.Post.UploadPostImage(file, header, post, logrus)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = handler.services.Post.UpdatePostImage(userId, postId, post, filePath, logrus)
	if err != nil {
		NewHandlerErrorResponse(ctx, http.StatusBadRequest, err.Error(), logrus)
		return
	}

	ctx.JSON(http.StatusOK, model.ResponseSuccess{Message: "Uploaded", Data: filePath})
}
