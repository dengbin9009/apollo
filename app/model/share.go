package model

import "github.com/jinzhu/gorm"

import "fmt"

// 分享的状态
const (
	ShareStatusCommon = 0
)

// Share 分享
type Share struct {
	Model
	UserID     uint   `gorm:"index"`             // 用户 ID
	User       *User  `gorm:"foreignkey:UserID"` // 用户
	URL        string `gorm:"varchar(1024)"`     // URL
	Title      string `gorm:"varchar(100)"`      // 分享的文章标题
	Review     string `gorm:"type:text"`         // 评论
	Status     int    `gorm:"default:0"`         // 状态
	Tag        string `gorm:"index;varchar(30)"` // 标签
	ClickCount int    `gorm:"default:0"`         // 点击数量(浏览数量)
	StarCount  int    `gorm:"default:0"`         // 赞的数量
}

// QueryBatch 检索一组
func (s *Share) QueryBatch(offset, pageSize int, userPreload bool, args ...interface{}) (shares []Share, total int, err error) {
	db := dbArgs(mydb, args...)
	err = db.Model(Share{}).Count(&total).Error
	if err != nil {
		return
	}
	if userPreload {
		db = db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,nickname")
		})
	}
	err = db.Offset(offset).
		Limit(pageSize).Order("id desc").
		Find(&shares).Error
	return
}

// QueryByID 根据 id 检索
func (s *Share) QueryByID(id interface{}) (share Share, err error) {
	db := mydb.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,nickname")
	})
	err = db.Where("id = ?", id).First(&share).Error
	return
}

// Create 创建分享
func (s *Share) Create() error {
	return mydb.Save(s).Error
}

// Update 更新
func (s *Share) Update() error {
	if s.ID == 0 {
		return fmt.Errorf("Share 更新必须设置 ID")
	}
	return mydb.Model(s).Updates(s).Error
}

// Click 点击量
func (s *Share) Click(shareID uint) error {
	err := mydb.Exec("UPDATE shares SET click_count = click_count + 1 WHERE id = ?", shareID).Error
	return err
}

// Star 点赞的数量
func (s *Share) Star(shareID uint) error {
	err := mydb.Exec("UPDATE shares SET star_count = star_count + 1 WHERE id = ?", shareID).Error
	return err
}

// AggregateTagCount 某个 Tag 有多少个 share
func (s *Share) AggregateTagCount(tagName string) (count int, err error) {
	err = mydb.Model(s).Where("tag=?", tagName).Count(&count).Error
	return
}
