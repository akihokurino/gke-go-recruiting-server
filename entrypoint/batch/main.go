package main

import (
	"context"
	"flag"
	"os"

	"gke-go-sample/di"
	pb "gke-go-sample/proto/go/pb"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(os.Getenv("GO_ENV")); err != nil {
		panic(err)
	}

	flag.Parse()
	args := flag.Args()

	taskName := args[0]

	ctx := context.Background()
	handler := di.ResolveBatchHandler()

	switch taskName {
	case "reindex-search":
		handler(ctx, pb.BatchTask_ReIndexSearch)
	case "proceed-work-status":
		handler(ctx, pb.BatchTask_ProceedWorkStatus)
	case "reorder-work":
		handler(ctx, pb.BatchTask_ReOrderWork)
	}
}
