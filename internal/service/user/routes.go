package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/zhetkerbaevan/personal-blog/internal/auth"
	"github.com/zhetkerbaevan/personal-blog/internal/models"
	"github.com/zhetkerbaevan/personal-blog/internal/utils"
)

type Handler struct {
	userStore models.UserStore
}

func NewHandler(userStore models.UserStore) *Handler {
	return &Handler{userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload models.RegisterUser
	//extract data from request
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	} 

	//validate payload
	if err := utils.Validator.Struct(payload); err != nil {
		//get validation error
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("INVALID PAYLOAD %v", errors))
		return
	}

	//check if user with email exists
	_, err := h.userStore.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("INVALID EMAIL"))
		return 
	}

	//if does not exist, we create new user start with hashing the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 	
	}

	//then we create new user
	err = h.userStore.CreateUser(models.User{
		Email: payload.Email,
		Password: hashedPassword,
		Name: payload.Name,
		Surname: payload.Surname,
		Age: payload.Age,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return 	
	}

	//succesfully created
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	
}