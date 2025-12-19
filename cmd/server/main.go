package main

import (
	"log"
	"time"

	"automation-engine/internal/api"
	"automation-engine/internal/middleware"
	"automation-engine/internal/service"

	"github.com/gin-gonic/gin"
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
	// 1. เชื่อมต่อ Database (GORM)
	cfg := myConfig.Config{
		User:   "",
		Passwd: "",
		Net:    "tcp",
		Addr:   ":3306",
		DBName: "hms_automation_engine",
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

	// 2. ประกอบร่างจิ๊กซอว์ (Dependency Injection)
	// DefinitionService จะสร้าง ActionRepository ภายในตัวมันเองตามที่คุณเขียนไว้
	defService := service.NewDefinitionService(db)

	// สร้าง Handler โดยส่ง Service เข้าไป
	defHandler := api.NewDefinitionHandler(defService)
	authHandler := api.NewAuthHandler()

	// 3. เริ่มต้นระบบ HTTP Server ด้วย Gin
	r := gin.Default()

	// Route สำหรับ Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		defGroup := protected.Group("/definition")
		{
			defGroup.GET("/actions", defHandler.GetActionByID)
			defGroup.POST("/actions", defHandler.CreateAction)
		}
	}

	// 5. รัน Server
	log.Println("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
