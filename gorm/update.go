package gorm

import (
	"fmt"
	"log"
)

func UpdatePK() {
	var c Content
	// 无主键
	if err := DB.Save(&c).Error; err != nil {
		log.Fatalln(err)
	}

	// [21.694ms] [rows:1] INSERT INTO `wsq_content` (`created_at`,
	// `updated_at`,`deleted_at`,`subject`,`likes`,`views`,`publish_time`,
	// `author_id`) VALUES ('2024-04-08 09:53:56.502','2024-04-08 09:53:56.502',
	// NULL,'',99,99,'2024-04-08 09:53:56.501',0) ON DUPLICATE KEY UPDATE
	// `updated_at`='2024-04-08 09:53:56.502',`deleted_at`=VALUES(`deleted_at`),
	// `subject`=VALUES(`subject`),`likes`=VALUES(`likes`),`views`=VALUES(`views`),
	// `publish_time`=VALUES(`publish_time`),`author_id`=VALUES(`author_id`)
	fmt.Printf("%+v\n", c)

	// 具有主键 ID
	if err := DB.Save(&c).Error; err != nil {
		log.Fatalln(err)
	}
	// [7.087ms] [rows:1] UPDATE `wsq_content` SET `created_at`='2024-04-08 09:57:45.023',
	// `updated_at`='2024-04-08 09:57:45.045',`deleted_at`=NULL,`subject`='',`likes`=99,
	// `views`=99,`publish_time`='2024-04-08 09:57:45.022',`author_id`=0 WHERE
	// `wsq_content`.`deleted_at` IS NULL AND `id` = 24
	fmt.Printf("%+v\n", c)
}

func UpdateWhere() {
	// 更新的字段值数据
	// map 推荐
	values := map[string]any{
		"subject": "Where Update Row",
		"likes":   10001,
	}

	// 执行带有条件的更新
	result := DB.Model(&Content{}).
		Where("likes > ?", 100).
		Updates(values)
	if result.Error != nil {
		log.Fatalln(result.Error)
	}

	// 获取更新结果，更新的记录数量（受影响的记录数）
	// 指的修改的
}
