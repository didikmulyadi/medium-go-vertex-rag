package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariable struct {
	AppName string
	Port    string
	Mode    string

	// Google
	GoogleProjectID string
	GoogleLocation  string
	VertexModel     string
}

func GetEnv() *EnvVariable {

	err := godotenv.Load()
	if err != nil {
		// Instead of fatal, just log a warning since env vars might be set through Docker
		log.Println("Warning: .env file not found, will use environment variables")
	}

	port := "8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	log.Println("RUN ON PORT: " + port)

	return &EnvVariable{
		Port:    port,
		Mode:    os.Getenv("GIN_MODE"),
		AppName: os.Getenv("APP_NAME"),

		// Google
		VertexModel:     os.Getenv("VERTEX_MODEL"),
		GoogleProjectID: os.Getenv("GOOGLE_PROJECT_ID"),
		GoogleLocation:  os.Getenv("GOOGLE_LOCATION"),
	}
}
