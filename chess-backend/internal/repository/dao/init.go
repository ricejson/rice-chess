package dao

import "gorm.io/gorm"

// InitTables 数据库模型迁移
func InitTables(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
