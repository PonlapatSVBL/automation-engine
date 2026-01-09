package main

import (
	"log"
	"os"
	"time"

	"automation-engine/internal/api"
	"automation-engine/internal/middleware"
	"automation-engine/internal/repository"
	"automation-engine/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	myConfig "github.com/go-sql-driver/mysql"

	_ "automation-engine/docs"
)

// @title           Automation Engine API
// @version         1.0
// @description     API for Automation System
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	txManager := repository.NewTransactionManager(db)
	conditionRepo := repository.NewConditionRepository(db)
	operatorRepo := repository.NewOperatorRepository(db)
	unitRepo := repository.NewUnitRepository(db)
	actionRepo := repository.NewActionRepository(db)
	conditionOperatorRepo := repository.NewConditionOperatorRepository(db)
	conditionUnitRepo := repository.NewConditionUnitRepository(db)
	conditionActionRepo := repository.NewConditionActionRepository(db)
	automationRepo := repository.NewAutomationRepository(db)
	automationActionRepo := repository.NewAutomationActionRepository(db)
	automationConditionRepo := repository.NewAutomationConditionRepository(db)
	automationExecutionRepo := repository.NewAutomationExecutionRepository(db)

	// 2. ประกอบร่างจิ๊กซอว์ (Dependency Injection)
	// DefinitionService จะสร้าง ActionRepository ภายในตัวมันเองตามที่คุณเขียนไว้
	definitionService := service.NewDefinitionService(
		txManager,
		actionRepo,
		conditionRepo,
		operatorRepo,
		unitRepo,
	)
	policyService := service.NewPolicyService(
		txManager,
		conditionRepo,
		operatorRepo,
		unitRepo,
		actionRepo,
		conditionOperatorRepo,
		conditionUnitRepo,
		conditionActionRepo,
	)
	runService := service.NewRunService(
		txManager,
		automationRepo,
		automationActionRepo,
		automationConditionRepo,
	)
	logService := service.NewLogService(
		txManager,
		automationExecutionRepo,
	)

	// สร้าง Handler โดยส่ง Service เข้าไป
	authHandler := api.NewAuthHandler()
	definitionHandler := api.NewDefinitionHandler(definitionService)
	policyHandler := api.NewPolicyHandler(policyService)
	runHandler := api.NewRunHandler(runService)
	logHandler := api.NewLogHandler(logService)

	// 3. เริ่มต้นระบบ HTTP Server ด้วย Gin
	r := gin.Default()

	// Route สำหรับ Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ui/policy-rule-config", func(c *gin.Context) {
		c.File("./web/policy-rule-config.html")
	})

	// 4. กำหนด Routes (เส้นทาง API)
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/login", authHandler.Login) // เส้นนี้ไม่ต้องใช้ Token
	}

	// Protected Routes (ต้องมี JWT)
	protected := apiV1.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// ย้ายกลุ่ม definition มาไว้ที่นี่
		definitionGroup := protected.Group("/definition")
		{
			definitionGroup.GET("/actions", definitionHandler.GetActionByID)
			definitionGroup.POST("/actions", definitionHandler.CreateAction)
		}

		policyGroup := protected.Group("/policy")
		{
			policyGroup.GET("/rule-config", policyHandler.GetPolicyRuleConfig)
			policyGroup.POST("/condition-operators", policyHandler.CreateConditionOperators)
			policyGroup.POST("/condition-units", policyHandler.CreateConditionUnits)
			policyGroup.POST("/condition-actions", policyHandler.CreateConditionActions)
		}

		runGroup := protected.Group("/run")
		{
			runGroup.POST("/automation", runHandler.CreateAutomation)
		}

		logGroup := protected.Group("/logs")
		{
			logGroup.GET("/automation-execution", logHandler.CreateAutomationExecution)
		}
	}

	// 5. รัน Server
	log.Println("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
