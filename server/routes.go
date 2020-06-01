package server

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) PublicRoutes() {

	s.Engine.NoRoute(s.NoRoute)

	s.Engine.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.Engine.GET("/", s.Welcome)
	s.Engine.GET("/ping", s.Ping)

	acc := s.Engine.Group("/accounts")
	{
		acc.POST("", s.CreateAccount)
		acc.GET("/:id", s.FindAccount)
	}

	tx := s.Engine.Group("/transactions")
	{
		tx.POST("", s.CreateTx)
		tx.GET("/account/:id", s.FindAllTxsByAccount)
	}
}
