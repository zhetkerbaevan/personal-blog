package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zhetkerbaevan/personal-blog/internal/service/post"
	"github.com/zhetkerbaevan/personal-blog/internal/service/user"
	"github.com/zhetkerbaevan/personal-blog/internal/store"
)

type APIServer struct {
	db *sql.DB
	address string
}

func NewAPIServer(db *sql.DB, address string) *APIServer { 
	//create new instance of APIServer struct
	return &APIServer{
		db : db,
		address: address,
	}
}

func (s *APIServer) Run() error {
	//create new router
	router := mux.NewRouter()
	
	//create subrouter (version of api)
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := store.NewUserStore(s.db)
	userService := user.NewHandler(userStore)
	userService.RegisterRoutes(subrouter)

	postStore := store.NewPostStore(s.db)
	postService := post.NewHandler(postStore, userStore)
	postService.RegisterRoutes(subrouter)

	log.Println("Listening on", s.address)

	return http.ListenAndServe(s.address, router)
}