package env

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/subosito/gotenv"
)

var (
	DATABASE_PORT     = 0
	API_PORT          = 0
	DATABASE_PASSWORD = ""
	DATABASE_HOST     = ""
	DATABASE_USER     = ""
	DATABASE_NAME     = ""
	DATABASE_DRIVE    = ""
	ENDOPOINT_API     = ""
	CONECT            = ""
)

func LoadEnv(ctx context.Context) error {
	err := gotenv.Load("env/aplication.env")
	if err != nil {
		log.Fatal("erro ao carregar env")
		return err
	}

	ENDOPOINT_API = os.Getenv("ENDOPOINT_API")
	DATABASE_DRIVE = os.Getenv("DATABASE_DRIVER")
	DATABASE_HOST = os.Getenv("DATABASE_HOST")
	DATABASE_NAME = os.Getenv("DATABASE_NAME")
	DATABASE_USER = os.Getenv("DATABASE_USER")
	DATABASE_PASSWORD = os.Getenv("DATABASE_PASSWORD")

	DATABASE_PORT, err = strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		log.Fatal("erro ao converter variavel db port")
		return err
	}

	API_PORT, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		log.Fatal("erro ao converter variavel API port")
		return err
	}

	CONECT = fmt.Sprintf("host=%s port=%d user=%s password=%v dbname=%s sslmode=disable", DATABASE_HOST, DATABASE_PORT, DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME)

	return nil
}
