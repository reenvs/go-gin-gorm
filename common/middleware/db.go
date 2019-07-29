package middleware

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"
)

func GetDbPrepareHandler(dbName, dbSource string, enableLog bool, contextDbName string) gin.HandlerFunc {
	db, err := gorm.Open(dbName, dbSource)
	if err != nil {
		return nil
	}

	if enableLog {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	db.DB().SetMaxIdleConns(50)
	db.DB().SetMaxOpenConns(150)
	//if host, err := os.Hostname(); host == "starsx-text" && err == nil {
	//	db.DB().SetConnMaxLifetime(time.Second * 3500)
	//}

	return func(c *gin.Context) {
		c.Set(contextDbName, db)
		//c.Header("Server-Version", constant.Version)
		c.Next()
	}
}

func GetEnginePrepareHandler(dbName, dbSource string, enableLog bool, contextDbName string) gin.HandlerFunc {
	db, err := xorm.NewEngine(dbName, dbSource)
	if err != nil {
		return nil
	}

	if enableLog {
		// todo 设置log
	}

	db.DB().SetMaxIdleConns(50)
	db.DB().SetMaxOpenConns(150)
	//if host, err := os.Hostname(); host == "starsx-text" && err == nil {
	//	db.DB().SetConnMaxLifetime(time.Second * 3500)
	//}

	return func(c *gin.Context) {
		c.Set(contextDbName, db)
		//c.Header("Server-Version", constant.Version)
		c.Next()
	}
}

func GetDbHandler(db *gorm.DB, contextDbName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(contextDbName, db)
		//c.Header("Server-Version", constant.Version)
		c.Next()
	}
}
