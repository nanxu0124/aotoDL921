package repository

import (
	"autoDL921/internal/domain"
	"autoDL921/internal/repository/dao"
	"context"
)

var (
	ErrUserDuplicateUserId = dao.ErrUserDuplicateUserId
	ErrUserNotFound        = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		UserId:   u.UserId,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByUserId(ctx context.Context, userId string) (domain.User, error) {
	u, err := repo.dao.FindByUserId(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		UserId:   u.UserId,
		Password: u.Password,
	}, nil
}
