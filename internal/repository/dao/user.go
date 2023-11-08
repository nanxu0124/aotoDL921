package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateUserId = errors.New("用户名已存在")
	ErrUserNotFound        = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	UserId   string `gorm:"unique"`
	Password string

	Ctime int64
	Utime int64
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateUserId
		}
	}
	return err
}

func (dao *UserDao) FindByUserId(ctx context.Context, userId string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("user_id = ?", userId).First(&u).Error
	return u, err
}
