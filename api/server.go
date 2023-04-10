package api

import (
	db "github.com/emmyvera/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our bankings system.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer returns a new HTTP server and setup our routeing
func NewServer(store *db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	// add toutes to router
	router.POST("/accounts", server.createAccount)

	server.router = router
	return server
}

// Start runs the server on a specific address
func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

// Handles errors messages
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
