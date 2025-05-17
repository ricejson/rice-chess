package service

import (
	"context"
	"errors"
	"github.com/ricejson/rice_chess/internal/domain"
	"github.com/ricejson/rice_chess/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate = repository.ErrUserDuplicate
	ErrUserNotExists = errors.New("user not found")
)

// UserService 用户服务
type UserService interface {
	// Login 登录业务
	Login(ctx context.Context, username string, pwd string) (domain.User, error)
	// Register 注册业务
	Register(ctx context.Context, username string, pwd string) error
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserServiceImpl(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (svc *UserServiceImpl) Login(ctx context.Context, username string, pwd string) (domain.User, error) {
	// 1. 根据用户名查找用户
	dbUser, err := svc.repo.FindByName(ctx, username)
	if err == repository.ErrUserNameNotFound {
		return domain.User{}, ErrUserNotExists
	}
	// 2. 比对密码是否一致
	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(pwd)); err != nil {
		return domain.User{}, ErrUserNotExists
	}
	return dbUser, nil
}

func (svc *UserServiceImpl) Register(ctx context.Context, username string, pwd string) error {
	// 1. 密码加密
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// 2. 保存用户
	err = svc.repo.Insert(ctx, domain.User{
		Username: username,
		Password: string(encryptedPwd),
	})
	if err == repository.ErrUserDuplicate {
		return ErrUserDuplicate
	}
	return nil
}
