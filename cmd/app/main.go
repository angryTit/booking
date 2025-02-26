package main

import (
	"applicationDesignTest/internal/api"
	"applicationDesignTest/internal/repo"
	"applicationDesignTest/internal/services"
	"applicationDesignTest/pkg/logger"
)

func main() {
	log := logger.NewZapLogger()

	orderRepo := repo.NewInMemoryOrder()
	availRepo := repo.NewinMemoryAvailability()

	availRepo.InitTestData()

	orderService := services.NewOrderService(orderRepo)
	roomService := services.NewRoomService(availRepo)
	bookingManager := services.NewBookingManager(orderService, roomService)

	server := api.NewServer(bookingManager, log)

	log.Info("Starting server on port 8080...")
	if err := server.Run("8080"); err != nil {
		log.Fatal("Failed to run server: %v", err)
	}
}
