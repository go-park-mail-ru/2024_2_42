package main

import (
	"fmt"
	"log"
	"net"
	"pinset/internal/app/db"
	"pinset/mailer-service/delivery"
	"pinset/mailer-service/mailer"
	"pinset/mailer-service/repository"
	"pinset/mailer-service/usecase"
	"pinset/pkg/logger"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	db := db.InitDB(logger)

	userRepo := repository.NewUserRepositoryController(db, logger)

	messageRepo := repository.NewMessageRepositoryController(db, logger)

	messageUsecase := usecase.NewMessageUsecase(messageRepo, userRepo)

	chatManager := delivery.NewMessageDeliveryController(messageUsecase)

	server := grpc.NewServer()

	mailer.RegisterChatServiceServer(server, chatManager)

	fmt.Println("starting server at :50051")

	server.Serve(lis)
}
