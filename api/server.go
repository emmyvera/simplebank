package api

import (
	"fmt"

	db "github.com/emmyvera/simplebank/db/sqlc"
	"github.com/emmyvera/simplebank/token"
	"github.com/emmyvera/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our bankings system.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer returns a new HTTP server and setup our routeing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot create token : %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// REgister the validator functions
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpServer()

	return server, nil
}

func (server *Server) setUpServer() {
	router := gin.Default()

	// add toutes to router for account
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:ID", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.DELETE("/accounts/:ID", server.delAccount)

	// add route for transfer
	router.POST("/transfers", server.createTransfer)

	// add route for user
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	server.router = router
}

// Start runs the server on a specific address
func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

// Handles errors messages
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
