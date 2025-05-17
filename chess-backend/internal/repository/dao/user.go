package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate    = errors.New("user duplicate")
	ErrUserNameNotFound = errors.New("username not found")
)

// UserDAO 用户数据访问层接口
type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByName(ctx context.Context, name string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewGORMUserDAO(db *gorm.DB) *GORMUserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

// Insert 插入用户记录
func (dao *GORMUserDAO) Insert(ctx context.Context, user User) error {
	currentTime := time.Now().UnixMilli()
	// 更新创建时间
	user.Ctime = currentTime
	// 更新更新时间
	user.Utime = currentTime
	err := dao.db.WithContext(ctx).Create(&user).Error
	// 特殊处理重复键错误
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexConflictError = 1062
		if mysqlErr.Number == uniqueIndexConflictError {
			return ErrUserDuplicate
		}
	}
	return err
}

// FindByName 根据用户名查找用户
func (dao *GORMUserDAO) FindByName(ctx context.Context, name string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("username=?", name).Find(&user).Error
	if err != nil {
		return User{}, ErrUserNameNotFound
	}
	return user, nil
}

// User 用户实体
type User struct {
	UserId   int64  `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"type:varchar(16);unique"`
	Password string `gorm:"type:varchar(256)"`
	// 天梯积分
	Score int `gorm:"default:1000"`
	// 总场数
	TotalCount int
	// 胜场数
	WinCount int
	// 创建时间（毫秒数）
	Ctime int64
	// 更新时间（毫秒数）
	Utime int64
}
