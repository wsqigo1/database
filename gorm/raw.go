package gorm

import "log"

// 原生查询测试
func RawSelect() {
	type Result struct {
		ID           uint
		Subject      string
		Likes, Views int
	}
	var rs []Result

	// SQL
	sql := "SELECT `id`, `subject`, `likes`, `views` FROM `wsq_content` WHERE `likes` > ? ORDER BY `likes` DESC LIMIT ?"

	// 执行SQL，并扫描结果
	if err := DB.Raw(sql, 99, 12).Scan(&rs).Error; err != nil {
		log.Fatalln(err)
	}

	// SELECT `id`, `subject`, `likes`, `views` FROM `wsq_content`
	// WHERE `likes` > 99 ORDER BY `likes` DESC LIMIT 12

	log.Println(rs)
}

// 执行类的 SQL 原生
func RawExec() {
	// SQL
	sql := "UPDATE `wsq_content` SET `subject` = CONCAT(`subject`, '-new postfix') WHERE `id` BETWEEN ? AND ?"

	// 执行，获取结果
	result := DB.Exec(sql, 20, 30)
	if result.Error != nil {
		log.Fatalln(result.Error)
	}
	// UPDATE `wsq_content` SET `subject` =
	// CONCAT(`subject`, '-new postfix') WHERE `id` BETWEEN 30 AND 40

	log.Println(result.RowsAffected)
}

// sql.ROw 或 sql.Rows 类型的结果处理
func RowsAndRow() {
	// sql
	sql := "SELECT `id`, `subject`, `likes`, `views` FROM `wsq_content` WHERE `likes` > ? ORDER BY `likes` DESC LIMIT ?"

	// 执行，获取 rows
	rows, err := DB.Raw(sql, 99, 12).Rows()
	if err != nil {
		log.Fatalln(err)
	}

	// 遍历 rows
	for rows.Next() {
		// 扫描的列独立的变量
		var id uint
		var subject string
		var likes, views int
		rows.Scan(&id, &subject, &likes, &views)
	}
}
