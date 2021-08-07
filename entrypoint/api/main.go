package main

import (
	"os"

	"gke-go-recruiting-server/di"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(os.Getenv("GO_ENV")); err != nil {
		panic(err)
	}

	di.ResolveAPIHandler()()
}
