package main

import (
	"automation-engine/internal/azbus"
	"automation-engine/internal/domain/model"
	"automation-engine/internal/dto"
	"automation-engine/internal/repository"
	"automation-engine/internal/service"
	"automation-engine/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	myConfig "github.com/go-sql-driver/mysql"
)

func main() {
	utils.LoadEnvVariables()

	ctx := context.Background()

	connStr := utils.GetEnv("SERVICE_BUS_CONNECTION_STRING", "")

	// New azure service bus client
	client, err := azservicebus.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		log.Fatalf("Failed to create Service Bus client: %v", err)
	}
	defer client.Close(ctx)

	sender, err := azbus.NewSender(ctx, client, "automate_queue")
	if err != nil {
		log.Fatalf("")
	}

	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DB")

	// เชื่อมต่อ Database (GORM)
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

	txManager := repository.NewTransactionManager(db)
	automationRepo := repository.NewAutomationRepository(db)
	automationActionRepo := repository.NewAutomationActionRepository(db)
	automationConditionGroupRepo := repository.NewAutomationConditionGroupRepository(db)
	automationConditionRepo := repository.NewAutomationConditionRepository(db)
	automationTargetRepo := repository.NewAutomationTargetRepository(db)
	automationExecutionRepo := repository.NewAutomationExecutionRepository(db)

	runService := service.NewRunService(
		txManager,
		automationRepo,
		automationActionRepo,
		automationConditionGroupRepo,
		automationConditionRepo,
		automationTargetRepo,
		automationExecutionRepo,
	)

	c := cron.New()

	// ตั้ง Cron ทำงานทุก 1 นาที
	c.AddFunc("* * * * *", func() {
		go runWorker(ctx, time.Now(), runService, sender)
	})

	c.Start()
	log.Println("Scheduler Started... Press Ctrl+C to exit")

	select {}
}

func runWorker(ctx context.Context, runTime time.Time, runService service.RunService, sender *azbus.Sender) {
	workerName := "Worker-at-" + runTime.Format("15:04:05")

	log.Printf("[%s] Worker cycle started", workerName)

	for {
		// 1. ดึงงานและจองสถานะเป็น PROCESSING (ผ่าน Skip Locked)
		tasks, err := runService.FetchAndLockTasks(ctx, runTime, 20)
		if err != nil {
			log.Printf("[%s] Error fetching tasks: %v", workerName, err)
			// ถ้า Error ให้หยุดพักสักครู่แล้วค่อยวนลูปใหม่
			time.Sleep(5 * time.Second)
			continue
		}

		// 2. ถ้า Query แล้วได้ 0 row ให้หยุด Loop
		if len(tasks) == 0 {
			log.Printf("[%s] No more tasks found. Worker cycle finished.", workerName)
			break
		}

		log.Printf("[%s] Picked up %d tasks", workerName, len(tasks))

		// 3. เริ่มประมวลผลงาน
		updatedTasks := make([]*model.RunAutomation, 0, len(tasks))
		for _, task := range tasks {
			log.Printf("[%s] Executing Automation ID: %s", workerName, task.AutomationID)

			var nextRunTime time.Time

			switch task.Frequency {
			case "daily":
				nextRunTime, err = utils.CalculateDailyNextRun(
					time.Now(),
					task.StartDate,
					time.Local,
				)
				if err != nil {
					continue
				}
			case "weekly":
				nextRunTime, err = utils.CalculateWeeklyNextRun(
					time.Now(),
					task.StartDate,
					task.DayOfWeek,
					time.Local,
				)
				if err != nil {
					continue
				}
			case "monthly":
				nextRunTime, err = utils.CalculateMonthlyNextRun(
					time.Now(),
					task.StartDate,
					int(task.DayOfMonth),
					time.Local,
				)
				if err != nil {
					continue
				}
			default:
				continue
			}

			task.Status = "PENDING"
			task.NextRunTime = nextRunTime
			task.LastUpd = time.Now()

			updatedTasks = append(updatedTasks, task)

			message := dto.MessageServiceBus{
				AutomationID: task.AutomationID,
			}
			body, err := json.Marshal(message)
			if err != nil {
				fmt.Println(err)
				continue
			}
			sender.SendMessage(ctx, task.InstanceServerChannelID, body)
		}

		// 4. อัปเดตครั้งเดียว (1 Database Round-trip)
		if err := runService.BulkUpdateNextRun(ctx, updatedTasks); err != nil {
			log.Printf("[%s] Bulk Update Error: %v", workerName, err)
		}
	}
}
