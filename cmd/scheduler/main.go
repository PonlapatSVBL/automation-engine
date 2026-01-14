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
	logService := service.NewLogService(
		txManager,
		automationExecutionRepo,
	)

	c := cron.New()

	// ตั้ง Cron ทำงานทุก 1 นาที
	c.AddFunc("* * * * *", func() {
		go runWorker(ctx, time.Now(), runService, logService, sender)
	})

	c.Start()
	log.Println("Scheduler Started... Press Ctrl+C to exit")

	select {}
}

func runWorker(ctx context.Context, runTime time.Time, runService service.RunService, logService service.LogService, sender *azbus.Sender) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Hour)
	defer cancel()

	workerID := runTime.Format("15:04:05")

	for {
		if err := ctx.Err(); err != nil {
			log.Printf("[Worker-%s] ⏱ Timeout/Cancelled: %v", workerID, err)
			return
		}

		// 1. Fetch & Lock
		tasks, err := runService.FetchAndLockTasks(ctx, runTime, 100)
		if err != nil {
			log.Printf("[Worker-%s] ❌ Error fetching tasks: %v", workerID, err)
			time.Sleep(5 * time.Second)
			return
		}

		if len(tasks) == 0 {
			log.Printf("[Worker-%s] ✅ No more tasks.", workerID)
			break
		}

		log.Printf("[Worker-%s] 📥 Picked up %d tasks", workerID, len(tasks))

		var successTasks []*model.RunAutomation

		for _, task := range tasks {
			// 2. คำนวณเวลาถัดไป
			nextRun, err := calculateNextRun(task)
			if err != nil {
				log.Printf("[Worker-%s] ⚠️ Skip Automation [%s]: %v", workerID, task.AutomationID, err)
				continue
			}

			// 3. เตรียม Message (DTO)
			msgPayload := dto.MessageServiceBus{
				LogID:        logService.GenerateLogID(),
				AutomationID: task.AutomationID,
				TriggeredAt:  time.Now(),
			}

			body, _ := json.Marshal(msgPayload)

			// 4. ส่งเข้า Service Bus (ถ้าพังจะไม่ update DB เพื่อให้รอบหน้ามาทำใหม่)
			err = sender.SendMessage(ctx, task.InstanceServerChannelID, body)
			if err != nil {
				log.Printf("[Worker-%s] ‼️ Failed to dispatch [%s] to bus: %v", workerID, task.AutomationID, err)
				continue
			}

			// เตรียมข้อมูลเพื่อ Update DB
			task.Status = "PENDING"
			task.NextRunTime = nextRun
			task.LastUpd = time.Now()
			successTasks = append(successTasks, task)

			log.Printf("[Worker-%s] 📤 Dispatched: AutomationID=%s | LogID=%s", workerID, task.AutomationID, msgPayload.LogID)
		}

		// 5. Bulk Update เฉพาะรายการที่ส่ง Bus สำเร็จ
		if len(successTasks) > 0 {
			if err := runService.BulkUpdateNextRun(ctx, successTasks); err != nil {
				log.Printf("[Worker-%s] ❌ Bulk Update Error: %v", workerID, err)
			} else {
				log.Printf("[Worker-%s] 💾 Successfully updated %d tasks in database", workerID, len(successTasks))
			}
		}
	}
}

func calculateNextRun(task *model.RunAutomation) (time.Time, error) {
	now := time.Now()

	switch task.Frequency {
	case "once":
		return time.Time{}, nil
	case "daily":
		return utils.CalculateDailyNextRun(now, task.StartDate, time.Local)
	case "weekly":
		return utils.CalculateWeeklyNextRun(now, task.StartDate, task.DayOfWeek, time.Local)
	case "monthly":
		return utils.CalculateMonthlyNextRun(now, task.StartDate, int(task.DayOfMonth), time.Local)
	case "yearly":
		return utils.CalculateYearlyNextRun(now, task.StartDate, int(task.DayOfMonth), int(task.MonthOfYear), time.Local)
	default:
		return time.Time{}, fmt.Errorf("unsupported frequency: %s", task.Frequency)
	}
}
