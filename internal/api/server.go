package api

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	Anuskh "github.com/nilesh0729/Transactly/internal/db/Result"
	"github.com/nilesh0729/Transactly/internal/token"
	"github.com/nilesh0729/Transactly/internal/util"
)

type Server struct {
	config     util.Config
	store      Anuskh.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store Anuskh.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create Token maker : %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.SetupRouter()

	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // For development only
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.POST("/user", server.CreateUser)

	router.POST("/user/login", server.LoginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.CreateAccount)

	authRoutes.GET("/accounts/:id", server.GetAccount)

	authRoutes.GET("/accounts", server.ListAccount)

	authRoutes.POST("/transfers", server.CreateTransfer)
	authRoutes.GET("/transfers", server.ListTransfer)
	authRoutes.GET("/accounts/:id/entries", server.ListEntry)

	server.router = router

}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
