package api

import (
	"database/sql"

	"github.com/gorilla/mux"
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
	subrouter := router.PathPrefix("/v1/api").Subrouter()
	
}