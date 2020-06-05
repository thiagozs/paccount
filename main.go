package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"paccount/database"
	"paccount/pkg/account"
	"paccount/pkg/oprtype"
	"paccount/pkg/transaction"
	"paccount/server"
	"path/filepath"

	_ "paccount/docs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	port    = flag.String("p", ":8080", "default port for api")
	debug   = flag.Bool("d", false, "log debug is on or off")
	verApp  = flag.Bool("version", false, "prints current program version")
	version = "1.0.0"
)

// @title Codding challenge
// @version 2.0
// @description Simple documentation of API.
// @termsOfService https://thiagozs.com/terms/

// @contact.name API Support
// @contact.url https://thiagozs.com
// @contact.email thiago.zilli@gmail.com

// @license.name Reserved Commons
// @license.url https://thiagozs.com/license

// @host localhost:8080
// @schemes http
// @BasePath /
func main() {

	flag.Parse()

	if *verApp {
		fmt.Println("Version : ", version)
		os.Exit(0)
	}

	log.Printf("Start Server on port %s...\n", *port)

	d, err := gorm.Open("sqlite3", filepath.Base("database.db"))
	if err != nil {
		log.Printf("Error on start data base, got: %s\n ", err.Error())
		os.Exit(1)
	}

	// options server...
	opts := func(s *server.Server) {
		s.Port = *port
		s.Debug = *debug
		s.Models = append(s.Models,
			&account.Entity{},
			&oprtype.Entity{},
			&transaction.Entity{},
		)
		s.DB = database.NewGormRepo(d)
	}

	s := server.New(opts)
	s.MigrationDB()
	s.PublicRoutes()

	if err := s.Run(); err != nil {
		log.Printf("Error on run server, got : %s\n", err.Error())
		os.Exit(1)
	}
}
