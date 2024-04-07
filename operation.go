package redis

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type User struct {
	gorm.Model

	Username string
	Name     string
	Email    string
	Birthday *time.Time
}

func OperatorType() {
	DB.AutoMigrate(&User{})

	var users []User

	// 一步操作
	err := DB.Where("birthday is NOT NULL").
		Where("email like ?", "@163.com%").
		Order("name DESC").
		Find(&users).Error // select all
	if err != nil {
		log.Fatal(err)
	}

	// 分布操作
	query := DB.Where("birthday IS NOT NULL")
	query = query.Where("email like ?", "@163.com%")
	query = query.Order("name DESC")
	query.Find(&users)
}

func CreateBasic() {
	DB.AutoMigrate(&Content{})

	// 模型映射记录，操作模型字段，就是操作记录的列
	c1 := Content{}
	c1.Subject = "GORM的使用"

	result1 := DB.Create(&c1)
	if result1.Error != nil {
		log.Fatal(result1.Error)
	}
	fmt.Println(c1.ID, result1.RowsAffected)

	// map 指定数据
	// 设置 map 的 values
	values := map[string]any{
		"Subject":     "Map指定值",
		"PublishTime": time.Now(),
	}
	// create
	result2 := DB.Model(&Content{}).Create(values)
	if result2.Error != nil {
		log.Fatal(result2.Error)
	}
	// 测试输出
	fmt.Println(result2.RowsAffected)
}

func CreateMulti() {
	DB.AutoMigrate(&Content{})

	// 定义模型的切片
	models := []Content{
		{Subject: "标题1"},
		{Subject: "标题2"},
		{Subject: "标题3"},
	}

	// create 插入
	result := DB.Create(&models)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	fmt.Println("RowsAffected:", result.RowsAffected)
	for _, m := range models {
		fmt.Println("ID:", m.ID)
	}

	// 切片结构同样支持
	vs := []map[string]any{
		{"Subject": "标题4"},
		{"Subject": "标题5"},
		{"Subject": "标题6"},
	}
	result2 := DB.Model(&Content{}).Create(vs)
	if result2.Error != nil {
		log.Fatal(result2.Error)
	}
	fmt.Println("RowsAffected:", result2.RowsAffected)
}

func CreateInBatches() {
	DB.AutoMigrate(&Content{})

	// 定义模型的切片
	models := []Content{
		{Subject: "标题1"},
		{Subject: "标题2"},
		{Subject: "标题3"},
	}

	// create 插入
	result := DB.CreateInBatches(&models, 2)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	fmt.Println("RowsAffected:", result.RowsAffected)
	for _, m := range models {
		fmt.Println("ID:", m.ID)
	}

	// 切片结构同样支持
	vs := []map[string]any{
		{"Subject": "标题4"},
		{"Subject": "标题5"},
		{"Subject": "标题6"},
	}
	result2 := DB.Model(&Content{}).CreateInBatches(vs, 2)
	if result2.Error != nil {
		log.Fatal(result2.Error)
	}
	fmt.Println("RowsAffected:", result2.RowsAffected)
}

func Upsert() {
	DB.AutoMigrate(&Content{})

	c1 := Content{}
	c1.Subject = "原始标题"
	c1.Likes = 10
	DB.Create(&c1)
	fmt.Println(c1)

	// 主键冲突的错误
	//c2 := Content{}
	//c2.ID = c1.ID
	//c1.Subject = "新标题"
	//c1.Likes = 20
	//if err := DB.Create(&c2).Error; err != nil {
	//	log.Fatal(err)
	//	// Error 1062 (23000): Duplicate entry '13' for key 'wsq_content.PRIMARY'
	//}

	// 冲突后，更新全部字段
	// INSERT INTO `wsq_content` (`created_at`, `updated_at`, `deleted_at`, `subject`, `likes`, `publish_time`, `id`)
	// VALUES ('2024-04-06 17:59:33.261', '2024-04-06 17:59:33.261', NULL, '新标题', 20, NULL, 15)
	// ON DUPLICATE KEY UPDATE `updated_at`='2024-04-06 17:59:33.261',
	//                         `deleted_at`=VALUES(`deleted_at`),
	//                         `subject`=VALUES(`subject`),
	//                         `likes`=VALUES(`likes`),
	//                         `publish_time`=VALUES(`publish_time`)
	//c3 := Content{}
	//c3.ID = c1.ID
	//c3.Subject = "新标题"
	//c3.Likes = 20
	//if err := DB.
	//	Clauses(clause.OnConflict{UpdateAll: true}).
	//	Create(&c3).Error; err != nil {
	//	log.Fatal(err)
	//}

	// 冲突后，更新部分字段
	c4 := Content{}
	c4.ID = c1.ID
	c4.Subject = "新标题"
	c4.Likes = 20
	if err := DB.
		Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"likes"}),
		}).
		Create(&c4).Error; err != nil {
		log.Fatal(err)
	}
}

func DefaultValue() {
	DB.AutoMigrate(&Content{})

	c1 := Content{}
	c1.Subject = "原始标题"
	c1.Likes = 0
	//views := uint(0)
	//c1.Views = &views
	DB.Create(&c1)
	fmt.Println(c1.Likes, c1.Views)
}

func SelectOmit() {
	DB.AutoMigrate(&Content{})

	c1 := Content{}
	c1.Subject = "原始标题"
	c1.Likes = 10
	c1.Views = 99
	now := time.Now()
	c1.PublishTime = &now

	//DB.Select("Subject", "Likes", "UpdatedAt").Create(&c1)
	//INSERT INTO `wsq_content` (`created_at`,`updated_at`,`subject`,`likes`)
	//VALUES ('2024-04-06 19:27:27.766','2024-04-06 19:27:27.766','原始标题',10)

	DB.Omit("Subject", "Likes", "UpdatedAt").Create(&c1)
}

func CreateHook() {
	DB.AutoMigrate(&Content{})

	c1 := Content{}
	err := DB.Create(&c1).Error
	if err != nil {
		log.Fatal(err)
	}
}
