package api

import (
	db "github.com/emmyvera/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our bankings system.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer returns a new HTTP server and setup our routeing
func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	// REgister the validator functions
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// add toutes to router for account
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:ID", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.DELETE("/accounts/:ID", server.delAccount)

	// add route for transfer
	router.POST("/transfers", server.createTransfer)

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
