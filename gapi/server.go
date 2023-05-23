package gapi

import (
	"fmt"

	db "github.com/emmyvera/simplebank/db/sqlc"
	"github.com/emmyvera/simplebank/pb"
	"github.com/emmyvera/simplebank/token"
	"github.com/emmyvera/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
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

	return server, nil
}
