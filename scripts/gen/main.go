package main

import (
	"fmt"
	"log"
	"os"
	"time"

	myConfig "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DB")

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
		panic(fmt.Sprintf("ไม่สามารถเชื่อมต่อ Database ได้: %v", err))
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/domain/query",
		ModelPkgPath: "model",
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)
	allModels := g.GenerateAllTable()
	g.ApplyBasic(allModels...)
	g.Execute()

	fmt.Println("Generate domain models สำเร็จแล้ว!")
}
