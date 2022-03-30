package main

import (
	"flag"
	"hitokoto-go/cmd"
	"hitokoto-go/global"
	"hitokoto-go/handlers/public"
	"hitokoto-go/inits"
	"log"
)

var (
	isImportMode bool
	isExportMode bool
	dataDir      string
)

func init() {
	flag.BoolVar(&isImportMode, "import", false, "Import sentences into database")
	flag.BoolVar(&isExportMode, "export", false, "Export sentences from database")
	flag.StringVar(&dataDir, "data-dir", "./sentences-bundle/", "Data directory of sentences bundle")
}

func main() {

	prepare()

	flag.Parse()

	// Check mutually exclusive flags

	if isImportMode && isExportMode {
		log.Fatal("Import and export mode are mutually exclusive")
		return
	}

	// Check commands

	if isImportMode {
		cmd.Import(dataDir)
		return
	} else if isExportMode {
		cmd.Export(dataDir)
		return
	}

	// Default (serv) mode
	serv()

}

func prepare() {

	// Initialize config
	if err := inits.Config(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Config initialized")

	// Connect to Database
	if err := inits.DB(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Database connected")

}

func serv() {

	log.Println("Starting application...")

	// Connect to Redis
	if err := inits.Redis(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Redis connected")

	// Initialize Meta
	if err := inits.Meta(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Meta initialized with", global.Meta.AllCount, "sentences")

	// Initialize Random Seeds
	if err := inits.RandomSeeds(); err != nil {
		log.Fatalln(err.Error())
	}

	// Initialize routes
	r := inits.Routes()
	log.Println("Routes initialized")

	// Start app server
	log.Println(public.RandGetOne()) // Startup message
	if err := r.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
