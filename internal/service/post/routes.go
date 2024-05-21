package post

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/zhetkerbaevan/personal-blog/internal/auth"
	"github.com/zhetkerbaevan/personal-blog/internal/models"
	"github.com/zhetkerbaevan/personal-blog/internal/utils"
)

type Handler struct {
	postStore models.PostStoreInterface
	userStore models.UserStoreInterface
}

func NewHandler(postStore models.PostStoreInterface, userStore models.UserStoreInterface) *Handler {
	return &Handler{postStore: postStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/posts", auth.WithJWTAuth(h.handleGetPosts, h.userStore)).Methods("GET")
	router.HandleFunc("/posts", auth.WithJWTAuth(h.handleCreatePost, h.userStore)).Methods("POST")
	router.HandleFunc("/posts/{id}", auth.WithJWTAuth(h.handleUpdatePost , h.userStore)).Methods("PUT")
	router.HandleFunc("/posts/{id}", auth.WithJWTAuth(h.handleDeletePost , h.userStore)).Methods("DELETE")
}

func (h *Handler) handleGetPosts(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIdFromContext(r.Context())
	posts, err := h.postStore.GetPostsByUserIds(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	utils.WriteJSON(w, http.StatusOK, posts)
}

func (h *Handler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIdFromContext(r.Context())
	var postPayload models.CreatePost
	//extract data from request
	if err := utils.ParseJSON(r, &postPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	//validate token
	if err := utils.Validator.Struct(postPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("INVALID PAYLOAD %v", errors))
		return
	}

	//create new post
	err := h.postStore.CreatePost(models.CreatePost{
		Title : postPayload.Title,
		Description: postPayload.Description,
		UserId: userId,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	//get postId from URL parametres
	params := mux.Vars(r)
	id := params["id"]

	postId, _ := strconv.Atoi(id) //convert it to int

	//check if post exists
	post, err := h.postStore.GetPostById(postId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	//extract data from request
	var postToUpdate models.UpdatePost
	if err := utils.ParseJSON(r, &postToUpdate); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	//update post with new info
	if postToUpdate.Title != "" {
		post.Title = postToUpdate.Title
	}
	if postToUpdate.Description != "" {
		post.Description = postToUpdate.Description
	}
	
	//update in db
	err = h.postStore.UpdatePost(post.Id, *post)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	//get postId from URL parametres
	params := mux.Vars(r)
	id := params["id"]

	postId, _ := strconv.Atoi(id)

	//check if post exists
	post, err := h.postStore.GetPostById(postId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}

	//delete post
	err = h.postStore.DeletePost(post.Id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

