package gorm

import (
	"log"
)

func DeleteWhere() {
	result1 := DB.Delete(&Content{}, "likes < ?", 100)
	if result1.Error != nil {
		log.Fatalln(result1.Error)
	}

	result2 := DB.Where("likes < ?", 100).Delete(&Content{})
	if result2.Error != nil {
		log.Fatalln(result2.Error)
	}
}

func FindDeleted() {
	var c Content
	DB.Delete(&c, 13)

	if err := DB.First(&c, 13).Error; err != nil {
		log.Println(err)
	}

	if err := DB.Unscoped().First(&c, 13).Error; err != nil {
		log.Fatalln(err)
	}
	// SELECT * FROM `wsq_content` WHERE `wsq_content`.`id` = 13 ORDER BY `wsq_content`.`id` LIMIT 1
	log.Printf("%+v\n", c)
}

func DeleteHard() {
	var c Content
	if err := DB.Unscoped().Delete(&c, 14).Error; err != nil {
		log.Fatalln(err)
	}
	// DELETE FROM `wsq_content` WHERE `wsq_content`.`id` = 14
}
