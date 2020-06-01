package server

import (
	"log"
	"paccount/database"
	"paccount/models"
	"time"

	"github.com/gin-gonic/gin"
)

type Option func(sr *Server)

type Server struct {
	Engine  *gin.Engine
	Models  []interface{}
	Port    string
	Debug   bool
	DB      database.IGormRepo
	OprType map[uint64]string
}

// New start a new service
func New(opts ...Option) *Server {
	s := Server{}

	// get all options need
	for _, option := range opts {
		option(&s)
	}

	// Set are debug or not
	gin.SetMode(gin.ReleaseMode)
	if s.Debug {
		gin.SetMode(gin.DebugMode)
	}

	// load gin framework
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(Cors())

	// Set Gin engine
	s.Engine = r
	s.OprType = map[uint64]string{}

	return &s
}

// StartDB start all process for connection database and models
func (s *Server) MigrationDB() {
	// Running the migrations
	s.DB.GetDB().AutoMigrate(s.Models...)
}

func (s *Server) Run() error {
	// forward runner
	return s.Engine.Run(s.Port)
}

// Seeds fillup database with some registrys
func (s *Server) Seeds() {
	s.seedOprType()
}

func (s *Server) seedOprType() {
	desc := []string{
		"COMPRA A VISTA",
		"COMPRA PARCELADA",
		"SAQUE",
		"PAGAMENTO",
	}
	for i, v := range desc {
		oprt := models.OprType{Description: v}
		var model models.OprType
		_ = s.DB.FindOne(oprt, &model)

		if (models.OprType{}) == model {
			log.Printf("Seed :: OperationType add '%s'\n", v)
			time := int32(time.Now().Unix())
			oprt.CreatedAt = time
			oprt.UpdatedAt = time
			if err := s.DB.Create(&oprt); err != nil {
				panic(err)
			}
			s.OprType[uint64(i)] = v
		} else {
			s.OprType[model.ID] = model.Description
		}
	}
}
