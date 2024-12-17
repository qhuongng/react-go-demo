package controllers

import (
	httpcommon "chi-mysql-boilerplate/internal/domain/http_common"
	"chi-mysql-boilerplate/internal/domain/models"
	"chi-mysql-boilerplate/internal/services"
	"chi-mysql-boilerplate/internal/utils/helpers"
	"chi-mysql-boilerplate/internal/utils/validators"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postService *services.PostService
	validator   *validators.Validator
}

func NewPostHandler(db *sql.DB, validator *validators.Validator) *PostHandler {
	return &PostHandler{postService: services.NewPostService(db), validator: validator}
}

// POST /posts
func (handler *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userId := GetUserIdFromContext(w, r)
	var req models.PostRequest
	if err := handler.validator.BindJSONAndValidate(w, r, &req); err != nil {
		// error is already handled in the validator
		return
	}

	newPost, err := handler.postService.Create(userId, req.Content)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}

	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&newPost))
}

// GET /posts
func (handler *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := handler.postService.GetAll()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}

	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&posts))
}

// GET /posts/by-user/{id}
func (handler *PostHandler) GetPostsByUserId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")

	// check if the id is of the correct format
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Field:   "userId",
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			}))
		return
	}

	posts, err := handler.postService.GetByUserId(uint64(userId))
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}

	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&posts))
}

// GET /posts/{id}
func (handler *PostHandler) GetPostById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// check if the id is of the correct format
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Field:   "id",
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			}))
		return
	}

	post, err := handler.postService.GetById(uint64(postId))
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}

	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&post))
}

// PUT /posts/{id}
func (handler *PostHandler) UpdatePostById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "Missing id parameter",
				Field:   "id",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			}))
		return
	}

	// check if the id is of the correct format
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Field:   "id",
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			}))
		return
	}

	var req models.PostRequest
	if err = handler.validator.BindJSONAndValidate(w, r, &req); err != nil {
		// error is already handled in the validator
		return
	}

	// check if the one sending the request is the author of the post
	if !handler.IsPostAuthor(w, r, uint64(postId)) {
		return
	}

	if err = handler.postService.UpdateById(uint64(postId), req.Content); err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}

	message := "Post updated successfully"
	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&message))
}

// DELETE /posts/{id}
func (handler *PostHandler) DeletePostById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Missing id parameter",
			Field:   "id",
			Code:    httpcommon.ErrorResponseCode.InvalidRequest,
		}))
		return
	}

	// check if the id is of the correct format
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Field:   "id",
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			}))
		return
	}

	// check if the one sending the request is the author of the post
	if !handler.IsPostAuthor(w, r, uint64(postId)) {
		return
	}

	if err = handler.postService.DeleteById(uint64(postId)); err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}

	message := "Post deleted successfully"
	helpers.WriteJSON(w, http.StatusOK, httpcommon.NewSuccessResponse(&message))
}

// helper function to grab the user ID from a HTTP context
func GetUserIdFromContext(w http.ResponseWriter, r *http.Request) uint64 {
	id, ok := r.Context().Value(httpcommon.ContextKeyConstants.UserId).(uint64)
	if !ok {
		helpers.MessageLogs.ErrorLog.Println("Failed to retrieve user ID from context")
		helpers.WriteJSON(w, http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.BadCredentials,
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			}))
		return 0
	}

	return id
}

// helper function to determine if the one updating or deleting the post is its author
func (handler *PostHandler) IsPostAuthor(w http.ResponseWriter, r *http.Request, postId uint64) bool {
	userId := GetUserIdFromContext(w, r)

	// fetch the post for the author's ID
	post, err := handler.postService.GetById(uint64(postId))
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return false
	}

	if post.UserID != userId {
		helpers.WriteJSON(w, http.StatusForbidden, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: httpcommon.ErrorMessage.InvalidRequest,
				Code:    httpcommon.ErrorResponseCode.Unauthorized,
			}))
		return false
	}

	return true
}
