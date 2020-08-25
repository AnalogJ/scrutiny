package database

import (
	"fmt"
	"github.com/analogj/scrutiny/webapp/backend/pkg/models/db"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func DatabaseHandler(dbPath string) gin.HandlerFunc {
	//var database *gorm.DB
	fmt.Printf("Trying to connect to database stored: %s", dbPath)
	database, err := gorm.Open("sqlite3", dbPath)

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&db.Device{})
	database.AutoMigrate(&db.SelfTest{})
	database.AutoMigrate(&db.Smart{})
	database.AutoMigrate(&db.SmartAtaAttribute{})
	database.AutoMigrate(&db.SmartNvmeAttribute{})
	database.AutoMigrate(&db.SmartScsiAttribute{})

	//TODO: detrmine where we can call defer database.Close()
	return func(c *gin.Context) {
		c.Set("DB", database)
		c.Next()
	}
}
