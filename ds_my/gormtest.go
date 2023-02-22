package main

import (
	"sync"

	"github.com/ds_my/dal/query"

	"github.com/ds_my/common"
	"github.com/ds_my/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ConnQuery *query.Query
var once sync.Once

// MySQLInit 初始化，将ConnQuery与数据库绑定
func MySQLInit() {
	once.Do(func() {
		ConnQuery = getQueryConnection()
	})
}

func getQueryConnection() *query.Query {
	var err error
	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(common.MySqlDSN))
	if err != nil {
		utils.Log.Fatal("数据库连接错误" + err.Error())
	} else {
		utils.Log.Info("MySQL连接成功")
	}
	ConnQuery = query.Use(db)
	return ConnQuery
}
