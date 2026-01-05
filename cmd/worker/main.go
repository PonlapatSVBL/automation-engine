package main

import (
	"automation-engine/internal/azbus"
	"automation-engine/internal/service"
	"automation-engine/internal/utils"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	myConfig "github.com/go-sql-driver/mysql"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DB")

	// 1. เชื่อมต่อ Database (GORM)
	cfg := myConfig.Config{
		User:   dbUser,
		Passwd: dbPass,
		Net:    "tcp",
		Addr:   dbHost + ":3306",
		DBName: dbName,
		Params: map[string]string{
			"charset":              "utf8mb4",
			"allowNativePasswords": "true",
		},
		ParseTime: true,
		Loc:       time.Local,
	}
	dsn := cfg.FormatDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	runService := service.NewRunService(db)

	// Create context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create WaitGroup
	var wg sync.WaitGroup

	// Load env
	utils.LoadEnvVariables()
	connStr := utils.GetEnv("SERVICE_BUS_CONNECTION_STRING", "")

	// New azure service bus client
	client, err := azservicebus.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		log.Fatalf("Failed to create Service Bus client: %v", err)
	}
	defer client.Close(ctx)

	// Create session receiver
	receiver1Opts := azbus.SessionReceiverOptions{
		SessionPool: utils.GetEnvAsInt("SESSION_POOL", 20),
		BatchSize:   utils.GetEnvAsInt("BATCH_SIZE", 5),
		ProcessPool: utils.GetEnvAsInt("PROCESS_POOL", 1),
	}
	receiver1 := azbus.NewSessionReceiver(ctx, &wg, client, "automate_queue", runService, &receiver1Opts)

	// Run session receiver
	go receiver1.RunDispatcher()

	// Wait for OS interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Received interrupt signal. Initiating graceful shutdown...")

	// Cancel context to notify goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()

	log.Println("All services shut down completely.")
}
