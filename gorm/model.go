package gorm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"reflect"
	"strings"
	"time"
)

type Post struct{ gorm.Model }
type Category struct{ gorm.Model }
type PostCategory struct{ gorm.Model }
type Box struct{ gorm.Model }

func (b *Box) TableName() string {
	return "my_box"
}

type TypeMap struct {
	gorm.Model

	FInt        int
	FUInt       uint
	FFloat32    float32
	FFloat64    float64
	FString     string
	FTime       time.Time
	FBytesSlice []byte

	FIntP     *int
	FUIntP    *uint
	FFloat32P *float32
	FFloat64P *float64
	FStringP  *string
	FTimeP    time.Time
}

func Migrate() {
	if err := DB.Debug().AutoMigrate(&Blog{}, &IAndC{}, &FieldTag{},
		&TypeMap{}, &Post{}, &Category{}, &PostCategory{}, &Box{}); err != nil {
		log.Fatal()
	}
}

func PointerDiff() {
	// 模型的零值
	typeMap := &TypeMap{}
	fmt.Printf("%+v\n", typeMap)

	fmt.Println("=========================")

	// 查询数据 NULL 对应的值
	DB.First(typeMap, 1)
	fmt.Printf("%+v\n", typeMap)
}

type CustomTypeModel struct {
	gorm.Model

	FTime    time.Time
	FNumTime sql.NullTime

	FString     string
	FNullString sql.NullString

	FUUID     uuid.UUID
	FNullUUID uuid.NullUUID
}

func CustomType() {
	// 初始化模型
	ctm := &CustomTypeModel{}

	// 迁移数据表
	DB.AutoMigrate(ctm)

	// 创建
	ctm.FTime = time.Now()        // 当前时间
	ctm.FNumTime = sql.NullTime{} // 零值，Valid默认为false
	ctm.FString = ""
	ctm.FNullString = sql.NullString{}

	ctm.FUUID = uuid.New()
	ctm.FNullUUID = uuid.NullUUID{}

	fmt.Println(DB.Create(ctm).Error)

	// 查询
	DB.First(ctm, ctm.ID)

	// 判断字段是否为NULL
	if ctm.FString == "" {
		fmt.Println("FString is NULL")
	} else {
		fmt.Println("FNString is NOT NULL")
	}

	if !ctm.FNullString.Valid {
		fmt.Println("FNullString is NULL")
	} else {
		fmt.Println("FNullString is NOT NULL")
	}
}

type FieldTag struct {
	gorm.Model

	// string 类型的处理
	FStringDefault string
	FTypeChar      string `gorm:"type:char(32)"`
	FTypeVarChar   string `gorm:"type:varchar(255)"`
	FTypeText      string `gorm:"type:text"`
	FTypeBlob      []byte `gorm:"type:blob"`
	FTypeEnum      string `gorm:"type:enum('Go', 'GORM', 'MySQL')"`
	FTypeSet       string `gorm:"type:set('Go', 'GORM', 'MySQL')"`

	FColNum     string `gorm:"column:custom_column_name"`
	FColNotNull string `gorm:"type:varchar(255);not null;default"`
	FColDefault string `gorm:"type:varchar(255);not null;default:gorm middle ware"`
	FColComment string `gorm:"type:varchar(255);comment:带有注释的字段"`
}

type IAndC struct {
	// 基础索引类型
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"type:varchar(255);unique"`
	Age   uint8  `gorm:"index;check:age>=18 AND email is not null"`

	// 符合索引
	FirstName string `gorm:"index:name"`
	LastName  string `gorm:"index:name"`
	// 顺序关键顺序
	// 默认的 priority:10
	FirstName1 string `gorm:"index:name,priority:2"`
	LastName1  string `gorm:"index:name,priority:1"`

	// 索引选项，前缀长度，排序方式，comment
	Height      float32 `gorm:"index:,sort:desc"`
	AddressHash string  `gorm:"index:,length:12,comment:前12个字符作为索引关键字"`
}

func IAndCCreate() {
	iac := &IAndC{}
	iac.Age = 18
	if err := DB.Create(iac).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", iac)
}

type Service struct {
	gorm.Model

	Url         string `gorm:"-"`
	Schema      string
	Host        string `gorm:"<-:false"`
	path        string `gorm:"<-:update"`
	QueryString string `gorm:">-:false"`
}

func ServiceCURD() {
	s := &Service{}
	s.Schema = "http"
	s.Url = "http:www.mashibing.com/"
	DB.Create(s)
}

type Paper struct {
	gorm.Model

	Subject string
	// 使用 json 序列化器进行处理
	Tags []string `gorm:"serializer:json"`
	// 使用自定义的编码器进行处理
	Categories []string `gorm:"serializer:csv"`
}

func PaperCurd() {
	if err := DB.AutoMigrate(&Paper{}); err != nil {
		log.Fatal()
	}

	paper := &Paper{}
	paper.Subject = "使用Serializer操作Tags字段"
	paper.Tags = []string{"Go", "Serializer", "Gorm", "MySQL"}
	// create 会执行序列化工作
	if err := DB.Create(paper).Error; err != nil {
		log.Fatal(err)
	}

	// 查询
	newPaper := &Paper{}
	// first 会执行反序列化工作
	DB.First(newPaper, 1)
	fmt.Printf("%+v\n", newPaper)
}

// 定义实现了序列化器接口的类型
type CSVSerializer struct{}

// Scan
// context.Context 对象
// field 模型的字段对应的类型
// dst 目标值（最终结果赋值到dst）
// dbValue 从数据库读取的值
// 错误
func (CSVSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value,
	dbValue any) error {
	// 初始化一个用来存储字段值的变量
	var fieldValue []string

	// 解析读取到的数据表的数据
	if dbValue != nil { // 不是 NULL
		// 支持解析的只有string和[]byte
		// 使用类型检测进行判定
		var str string
		switch v := dbValue.(type) {
		case string:
			str = v
		case []byte:
			str = string(v)
		default:
			return fmt.Errorf("failed to unmarshal CSV value: %#v", dbValue)
		}
		// 二：核心：将数据表中的字段使用逗号分割，形成 []string
		fieldValue = strings.Split(str, ",")
	}

	// 将处理好的数据，设置到dst上
	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(fieldValue))
	return nil
}

// 实现Value
// fieldValue 模型的字段值
func (CSVSerializer) Value(ctx context.Context, field *schema.Field,
	dst reflect.Value, fieldValue any) (any, error) {
	return strings.Join(fieldValue.([]string), ","), nil
}

func CustomSerializer() {
	// 注册序列化器
	schema.RegisterSerializer("csv", CSVSerializer{})

	// 测试
	if err := DB.AutoMigrate(&Paper{}); err != nil {
		log.Fatal()
	}

	paper := &Paper{}
	paper.Subject = "使用自定义的Serializer操作Categories字段"
	paper.Tags = []string{"Go", "Serializer", "Gorm", "MySQL"}
	paper.Categories = []string{"Go", "Serializer", "Gorm", "MySQL"}
	// create 会执行序列化工作
	if err := DB.Create(paper).Error; err != nil {
		log.Fatal(err)
	}

	// 查询
	newPaper := &Paper{}
	// first 会执行反序列化工作
	DB.First(newPaper, paper.ID)
	fmt.Printf("%+v\n", newPaper)
}

type Blog struct {
	gorm.Model

	BlogBasic
	Author `gorm:"embeddedPrefix:author_"`
}

type BlogBasic struct {
	Subject string
	Summary string
	Content string
}
