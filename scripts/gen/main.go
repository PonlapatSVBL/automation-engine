package main

import (
	"fmt"
	"time"

	myConfig "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
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
