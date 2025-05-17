package repository

import (
	"context"
	"github.com/ricejson/rice_chess/internal/domain"
	"github.com/ricejson/rice_chess/internal/repository/dao"
)

var (
	ErrUserDuplicate    = dao.ErrUserDuplicate
	ErrUserNameNotFound = dao.ErrUserNameNotFound
)

type UserRepository interface {
	Insert(ctx context.Context, user domain.User) error
	FindByName(ctx context.Context, name string) (domain.User, error)
}

type CachedUserRepository struct {
	dao dao.UserDAO
}

func NewCachedUserRepository(dao dao.UserDAO) *CachedUserRepository {
	return &CachedUserRepository{
		dao: dao,
	}
}

func (repo *CachedUserRepository) Insert(ctx context.Context, user domain.User) error {
	err := repo.dao.Insert(ctx, dao.User{
		Username: user.Username,
		Password: user.Password,
	})
	if err == dao.ErrUserDuplicate {
		return ErrUserDuplicate
	}
	return err
}

func (repo *CachedUserRepository) FindByName(ctx context.Context, name string) (domain.User, error) {
	user, err := repo.dao.FindByName(ctx, name)
	if err == dao.ErrUserNameNotFound {
		return domain.User{}, ErrUserNameNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		UserId:     user.UserId,
		Username:   user.Username,
		Score:      user.Score,
		TotalCount: user.TotalCount,
		WinCount:   user.WinCount,
		Ctime:      user.Ctime,
		Utime:      user.Utime,
	}, nil
}
