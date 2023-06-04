package comment

import (
	"fmt"
	"ginEssential/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	_ = viper.GetString("datasource.driverName")
	username := viper.GetString("datasource.username") //账号
	password := viper.GetString("datasource.password") //密码
	host := viper.GetString("datasource.host")         //数据库地址，可以是Ip或者域名
	port := viper.GetString("datasource.port")         //数据库端口
	database := viper.GetString("datasource.database") //数据库名
	timeout := "10s"                                   //连接超时，10秒

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, database, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	// 连接成功
	fmt.Println(db)
	db.AutoMigrate(&model.User{})
	DB = db
}
func GetDB() *gorm.DB {
	InitDB()
	return DB
}
