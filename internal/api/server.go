package api

import (
	"fmt"
	db "github.com/alexispell/minibank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store // db interactions
	router *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing
func NewServer(store *db.Store) *Server {
	router := gin.Default()

	server := &Server{
		store:  store,
		router: router,
	}

	// add routes to router
	server.router.POST("accounts", server.createAccount)
	server.router.GET("accounts/:id", server.getAccountById)
	server.router.GET("accounts", server.listAccounts)

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error, description ...string) gin.H {
	res := gin.H{
		"error": err.Error(),
	}
	if len(description) != 0 && description[0] != "" {
		fmt.Println("Err Descriptions:", description[0])
		res["description"] = string(description[0])
	}
	return res
}
