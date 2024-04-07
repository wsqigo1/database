package redis

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Content struct {
	gorm.Model

	Subject     string
	Likes       uint `gorm:"default:99"`
	Views       uint `gorm:"default:99"`
	PublishTime *time.Time

	// 不需要迁移
	// 禁用写操作
	Sv string `gorm:"-:migration;<-:false"`

	// 作者ID
	AuthorID uint
}

func NewContent() Content {
	return Content{
		Likes: 99,
		Views: 99,
	}
}

func (c *Content) BeforeCreate(db *gorm.DB) error {
	// 业务
	if c.PublishTime == nil {
		now := time.Now()
		c.PublishTime = &now
	}

	// 配置
	db.Statement.AddClause(clause.OnConflict{UpdateAll: true})

	return nil
}

func (c *Content) AfterCreate(db *gorm.DB) error {
	return errors.New("custom error")
}

type ContentStrPK struct {
	ID          string `gorm:"primaryKey"`
	Subject     string
	Likes       uint `gorm:"default:99"`
	Views       uint `gorm:"default:99"`
	PublishTime *time.Time
}
