package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/sharding"
	"time"
)

var DB *gorm.DB

func init() {
	username := "root"   //账号
	password := "root"   //密码
	host := "127.0.0.1"  //数据库地址，可以是Ip或者域名
	port := 3306         //数据库端口
	Dbname := "gorm_001" //数据库名
	timeout := "10s"     //连接超时，10秒

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	err = db.Use(sharding.Register(sharding.Config{
		ShardingKey:         "id",
		NumberOfShards:      2,
		PrimaryKeyGenerator: sharding.PKSnowflake, //主键生成规则，这里面雪花算法
	}, "user_models"))
	if err != nil {
		fmt.Println("分片加载失败", err.Error())
		return
	}
	fmt.Println("分片生成成功")
	err = db.AutoMigrate(&UserModel{})
	if err != nil {
		fmt.Println("表结构生成失败", err.Error())
		return
	}
	fmt.Println("表结构生成成功")
	DB = db
}

type Model struct {
	ID        int64     `gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserModel struct {
	Model
	UserName string `json:"userName"`
	Pwd      string `json:"-"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

func main() {
	//DB.Create(&UserModel{
	//	Model: Model{
	//		ID: 3,
	//	},
	//	UserName: "fengfeng1",
	//})
	//DB.Create(&UserModel{
	//	Model: Model{
	//		ID: 4,
	//	},
	//	UserName: "zhangsan1",
	//})

	var user UserModel
	DB.Take(&user, "id = ?", 2)
	fmt.Println(user)
	user = UserModel{}
	DB.Table("user_models_1").Take(&user, "user_name = ?", "fengfeng")
	fmt.Println(user)
}
