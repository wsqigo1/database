package gorm

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

func GetByPk() {
	// migrate
	DB.AutoMigrate(&Content{}, &ContentStrPK{})

	// 查询单条
	c := Content{}
	if err := DB.First(&c, 10).Error; err != nil {
		log.Println(err)
	}

	// 字符串类型的主键
	cStr := ContentStrPK{}
	if err := DB.First(&cStr, "id = ?", "some pk").Error; err != nil {
		log.Println(err)
	}

	// 查询多条
	var cs []Content
	if err := DB.Find(&cs, []uint{10, 11, 12}).Error; err != nil {
		log.Println(err)
	}

	// 字符串类型的主键
	var cStrs []ContentStrPK
	if err := DB.Find(&cStrs, "id in ?", []string{"some", "pk", "item"}).Error; err != nil {
		log.Println(err)
	}
}

func GetToMap() {
	// 单条
	c := map[string]any{}
	if err := DB.Model(&Content{}).First(&c, 13).Error; err != nil {
		log.Println(err)
	}
	// 需要接口类型断言，才能继续处理
	if c["id"].(uint) == 13 {
		fmt.Println("id bingo")
	}

	// time 类型的处理
	t, err := time.Parse("2006-01-02 15:04:05.000 -0700 CST", "2024-04-06 17:56:14.499 +0800 CST")
	if err != nil {
		log.Println(err)
	}
	if c["created_at"].(time.Time) == t {
		fmt.Println("created_at bingo")
	}

	// 多条
	var cs []map[string]any
	if err := DB.Model(&Content{}).Find(&cs, []uint{13, 14, 15}).Error; err != nil {
		log.Println(err)
	}
	for _, c := range cs {
		fmt.Println(c["id"].(uint), c["subject"].(string), c["created_at"])
	}
}

func GetPluck() {
	// 使用切片存储
	var subjects []sql.NullString
	if err := DB.Model(&Content{}).Pluck("concat(subject, '-', likes)", &subjects).Error; err != nil {
		log.Println()
	}

	for _, subject := range subjects {
		// NullString 的使用
		//if subject.Valid {
		//	fmt.Println(subject.String)
		//} else {
		//	fmt.Println("[NULL]")
		//}

		fmt.Println(subject)
	}
}

func GetPluckExp() {
	// 使用切片存储，如果表达式可以保证NULL不会出现了，就可以不适用NullType了
	var subjects []string
	// 字段为表达式的结果
	if err := DB.Model(&Content{}).Pluck("concat(coalesce(subject, '[no subject]'), '-', likes)", &subjects).Error; err != nil {
		log.Println()
	}

	for _, subject := range subjects {
		fmt.Println(subject)
	}
}

func GetSelect() {
	var c Content

	// 基本的字段名
	if err := DB.Select("subject", "likes").First(&c, 13).Error; err != nil {
		log.Fatal(err)
	}

	// 字段表达式
	if err := DB.Select("subject", "likes", "concat(subject, '-', views) AS sv").
		First(&c, 13).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", c)
}

func GetDistinct() {
	var cs []Content

	if err := DB.Distinct("*").Find(&cs).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cs)
}

func WhereMethod() {
	var cs []Content

	// inline条件
	if err := DB.Find(&cs, "likes > ? AND subject like ?", 100, "gorm%").Error; err != nil {
		log.Fatal(err)
	}
	// SELECT * FROM `wsq_content` WHERE (likes > 100 AND subject like
	// 'gorm%') AND `wsq_content`.`deleted_at` IS NULL

	// Where，通常在动态拼凑条件时使用

	//query := DB.Where("likes > ?", 100)
	//subject := "from user data"
	//// 当前用户输出subject，不为空字符串时，才拼凑subject条件
	//if subject != "" {
	//	query.Where("subject like ?", subject+"%")
	//}
	//if err := query.Find(&cs).Error; err != nil {
	//	log.Fatal(err)
	//}

	// OR 逻辑运算
	//query := DB.Where("likes > ?", 100)
	//subject := "gorm"
	//// 当前用户输出subject，不为空字符串时，才拼凑subject条件
	//if subject != "" {
	//	query.Or("subject like ?", subject+"%")
	//}
	//if err := query.Find(&cs).Error; err != nil {
	//	log.Fatal(err)
	//}
	// SELECT * FROM `wsq_content` WHERE (likes > 100 OR subject like
	// 'gorm%') AND `wsq_content`.`deleted_at` IS NULL

	// Not 逻辑运算
	query := DB.Where("likes > ?", 100)
	subject := "gorm"
	// 当前用户输出subject，不为空字符串时，才拼凑subject条件
	if subject != "" {
		query.Or(DB.Not("subject like ?", subject+"%"))
	}
	if err := query.Find(&cs).Error; err != nil {
		log.Fatal(err)
	}
}

func WhereType() {
	var cs []Content

	// (1 or 2) and (3 and (4 or 5))
	// 1 or 2
	//condA := DB.Where("likes > ?", 10).Or("likes <= ?", 100) // 1
	//// 3 and (4 or 5)
	//condB := DB.Where("views > ?", 20).Where(DB.Where("views <= ?", 200).
	//	Or(DB.Where("subject like ?", "gorm%")))
	//query := DB.Where(condA).Where(condB)
	// SELECT * FROM `wsq_content` WHERE (likes > 10 OR likes <= 100) AND
	// (views > 20 AND (views <= 200 OR subject like 'gorm%')) AND `wsq_content`.`deleted_at` IS NULL

	// map 构建条件, and, = in
	//query := DB.Where(map[string]any{
	//	"views": 100,
	//	"id":    []uint{1, 2, 3, 4, 5},
	//})

	// struct 条件构建
	query := DB.Where(Content{
		Views:   100,
		Subject: "GORM",
	})
	if err := query.Find(&cs).Error; err != nil {
		log.Fatal(err)
	}
}

func PlaceHolder() {
	var cs []Content

	// 具名，绑定名字sql.Named()结构
	//query := DB.Where("likes = @like AND subject like @subject",
	//	sql.Named("subject", "gorm%"), sql.Named("like", 100))

	// gorm 还支持使用map的形式具名绑定
	query := DB.Where("likes = @like AND subject like @subject",
		map[string]any{
			"subject": "gorm%",
			"like":    100,
		})

	if err := query.Find(&cs).Error; err != nil {
		log.Fatal(err)
	}
}

func OrderBy() {
	var cs []Content

	ids := []uint{2, 3, 1}
	//query := DB.Order("FIELD(id, 2, 3, 1)")
	query := DB.Clauses(clause.OrderBy{
		Expression: clause.Expr{
			SQL:                "FIELD(id, ?)",
			Vars:               []any{ids},
			WithoutParentheses: true,
		},
	})

	if err := query.Find(&cs, ids).Error; err != nil {
		log.Fatal(err)
	}

	for _, c := range cs {
		fmt.Println(c.ID)
	}
}

// 默认的值
const (
	DefaultPage     = 1
	DefaultPageSize = 12
)

// 定义分页必要数据结构
type Pager struct {
	Page, PageSize int
}

// 翻页程序
func Pagination(pager Pager) {
	// 确定 offset 和 pageSize
	page := DefaultPage
	if pager.PageSize != 0 {
		page = pager.PageSize
	}

	pageSize := DefaultPageSize
	if pager.PageSize != 0 {
		pageSize = pager.PageSize
	}

	// 计算 offset
	// page, pageSize, offset
	// 1, 10, 0
	// 2, 10, 10
	// 3, 10, 20
	offset := pageSize * (page - 1)

	var cs []Content
	if err := DB.Offset(offset).Limit(pageSize).Find(&cs).Error; err != nil {
		log.Fatal(err)
	}
}

// 用于得到 func(db *gorm.DB) *gorm.DB 类型函数
// 为什么不直接定义函数，因为需要 func(db *gorm.DB) *gorm.DB 与分页信息产生联系
func Paginate(pager Pager) func(db *gorm.DB) *gorm.DB {
	// 确定 offset 和 pageSize
	page := DefaultPage
	if pager.PageSize != 0 {
		page = pager.PageSize
	}

	pageSize := DefaultPageSize
	if pager.PageSize != 0 {
		pageSize = pager.PageSize
	}

	offset := pageSize * (page - 1)

	return func(db *gorm.DB) *gorm.DB {
		// 使用闭包的变量，实现翻页的业务逻辑
		return db.Limit(pageSize).Offset(offset)
	}
}

// 测试重用的翻页逻辑
func PaginationScope(pager Pager) {
	var cs []Content
	if err := DB.Scopes(Paginate(pager)).Find(&cs).Error; err != nil {
		log.Fatal(err)
	}
}

func GroupHaving() {
	DB.AutoMigrate(&Content{})

	// 定义查询结构类型
	type Result struct {
		// 分组字段
		AuthorID uint

		// 合计字段
		TotalViews int
		TotalLikes int
		AvgViews   float64
	}

	// 执行分组合计过滤查询
	var rs []Result
	if err := DB.Model(&Content{}).
		Select("author_id", "SUM(views) as total_views",
			"SUM(likes) as total_likes", "AVG(views) as avg_views").
		Group("author_id").
		Having("total_views > ?", 99).
		Find(&rs).Error; err != nil {
		log.Fatalln(err)
	}
	// SQL
	// SELECT `author_id`,SUM(views) as total_views,SUM(likes) as
	// total_likes,AVG(views) as avg_views FROM `wsq_content` WHERE
	// `wsq_content`.`deleted_at` IS NULL GROUP BY `author_id` HAVING total_views > 99
}

func CountPage(pager Pager) {
	// 集中的条件，用于统计数量和获取某页记录
	query := DB.Model(&Content{}).
		Where("likes > ?", 99)

	// total rows count
	var count int64
	if err := query.Count(&count).Error; err != nil {
		log.Println(err)
	}
	// 计算总页数 ceil(count / pageSize)

	// rows per page
	var cs []Content
	if err := query.Scopes(Paginate(pager)).Find(&cs); err != nil {
		log.Fatal(err)
	}
}

func Iterator() {
	// 利用 DB.Rows() 获取 Rows 对象
	rows, err := DB.Model(&Content{}).Rows()
	if err != nil {
		log.Fatal(err)
	}
	// 注意：保证使用过后关闭 rows 结果集
	defer rows.Close()
	fmt.Println(rows)

	// 迭代的从Rows中扫描记录到模型
	for rows.Next() {
		var c Content
		if err := DB.ScanRows(rows, &c); err != nil {
			log.Fatal(err)
		}
		fmt.Println(c.Subject)
	}
}

func Locking() {
	var cs []Content

	if err := DB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&cs).Error; err != nil {
		log.Fatal(err)
	}
	// [22.163ms] [rows:20] SELECT * FROM `wsq_content` WHERE `wsq_content`.`deleted_at`
	// IS NULL FOR UPDATE

	if err := DB.
		Clauses(clause.Locking{Strength: "SHARE"}).
		Find(&cs).Error; err != nil {
		log.Fatal(err)
	}
	// [22.163ms] [rows:20] SELECT * FROM `wsq_content` WHERE `wsq_content`.`deleted_at`
	// IS NULL FOR SHARE
}

func SubQuery() {
	// migrate
	DB.AutoMigrate(&Author{}, &Content{})

	// 条件子查询
	// select * from content where author_id in (select id from author where status = 0);
	// 子查询，不需要使用终结方法 Find 完成查询，只需要构建语句即可
	whereSubQuery := DB.Model(&Author{}).
		Select("id").
		Where("status = ?", 0)
	var cs []Content
	if err := DB.Where("author_id IN (?)", whereSubQuery).Find(&cs).Error; err != nil {
		log.Fatal(err)
	}

	// from 子查询
	// select * from (select subject, likes from content
	// where publish_time is null) ass temp where likes > 10;
	fromSubQuery := DB.Model(&Content{}).
		Select("subject", "likes").
		Where("publish_time IS NULL")
	type Result struct {
		Subject string
		Likes   int
	}
	var rs []Result
	if err := DB.Table("(?) AS temp", fromSubQuery).
		Where("likes > ?", 10).
		Find(&rs).Error; err != nil {
		log.Fatal(err)
	}
}

func FindHook() {
	var c Content
	if err := DB.First(&c, 13).Error; err != nil {
		log.Fatal(err)
	}

	fmt.Println(c.AuthorID)
}
