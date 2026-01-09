package main

import (
	"automation-engine/internal/domain/model"
	"automation-engine/internal/service"
	"automation-engine/internal/utils"
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
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

	c := cron.New()

	// ตั้ง Cron ทำงานทุก 1 นาที
	c.AddFunc("* * * * *", func() {
		go runWorker("Worker-at-"+time.Now().Format("15:04:05"), runService)
	})

	c.Start()
	log.Println("Scheduler Started... Press Ctrl+C to exit")

	select {}
}

func runWorker(workerName string, runService service.RunService) {
	ctx := context.Background()
	log.Printf("[%s] Worker cycle started", workerName)

	for {
		// 1. ดึงงานและจองสถานะเป็น PROCESSING (ผ่าน Skip Locked)
		tasks, err := runService.FetchAndLockTasks(ctx, 20)
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
		}

		// 4. อัปเดตครั้งเดียว (1 Database Round-trip)
		if err := runService.BulkUpdateNextRun(ctx, updatedTasks); err != nil {
			log.Printf("[%s] Bulk Update Error: %v", workerName, err)
		}
	}
}
