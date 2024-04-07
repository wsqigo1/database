package redis

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

// 创建模型
type Article struct {
	// 潜入基础模型
	gorm.Model
	// 定义字段
	Subject     string
	Likes       uint
	Published   bool
	PublishTime time.Time
	AuthorID    uint
}

func BasicUsage() {
	// 定义 dsn
	dsn := "root:secret@tcp(127.0.0.1:13306)/gormExample?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接服务器（池）
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 基于模型完成表结构（设计）的迁移（
	if err := db.AutoMigrate(&Article{}); err != nil {
		log.Fatal(err)
	}
	log.Println("migrate success")
}

func Create() {
	// 构建Article类型数据
	article := &Article{
		Subject:     "GORM 的 CURD 基础操作",
		Likes:       0,
		Published:   true,
		PublishTime: time.Now(),
		AuthorID:    42,
	}

	// DB.Create 完成数据库的insert
	if err := DB.Create(article).Error; err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Println(article)
}

func Retrieve(id uint) {
	// 初始化Article模型，零值
	article := &Article{}

	if err := DB.First(article, id).Error; err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Println(article)
}

func Update() {
	// 获取需要更新的对象
	article := &Article{}
	if err := DB.First(article, 1).Error; err != nil {
		log.Fatal(err)
	}

	// 更新对象对象
	article.AuthorID = 23
	article.Likes = 101
	article.Subject = "新的文章标题"

	// 存储，DB.Save()
	if err := DB.Save(article).Error; err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Println(article)
}

func Delete() {
	// 获取模型对象
	article := &Article{}
	if err := DB.First(article, 1).Error; err != nil {
		log.Fatal(err)
	}

	if err := DB.Delete(article).Error; err != nil {
		log.Fatal(err)
	}

	// print
	fmt.Println("article was deleted")
}

func Debug() {
	// insert
	article := &Article{
		Subject:     "GORM 的 CURD 基础操作",
		PublishTime: time.Now(),
	}

	if err := DB.Debug().Create(article).Error; err != nil {
		log.Fatal(err)
	}

	// select
	if err := DB.Debug().First(article, article.ID).Error; err != nil {
		log.Fatal(err)
	}
}

func Log() {
	Update()
}
