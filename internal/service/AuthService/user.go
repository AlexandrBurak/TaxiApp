package AuthService

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/AlexandrBurak/TaxiApp/internal/cache"
	"github.com/AlexandrBurak/TaxiApp/internal/logger"
	"github.com/AlexandrBurak/TaxiApp/internal/model"
	"github.com/AlexandrBurak/TaxiApp/internal/repository"
)

type Repository interface {
	AddNewUser(ctx context.Context, user model.User) error
	GetUserPhoneAndPasswordByPhone(ctx context.Context, user model.User) (model.User, error)
	Exists(ctx context.Context, user model.User) (bool, error)
}

type Cache interface {
	CacheUser(ctx context.Context, user model.User) error
	ExistsInCache(ctx context.Context, phone string) (model.User, error)
}

type Logger interface {
	Log(err error)
	Error(err error)
}

type Service struct {
	Repos Repository
	Log   Logger
	Cache Cache
}

func NewService(repository repository.Repository, log logger.Logger, cache cache.Cache) Service {
	return Service{Repos: &repository, Log: &log, Cache: &cache}
}

func (s *Service) Save(ctx context.Context, user model.User) error {
	exists, err := s.Repos.Exists(ctx, user)
	if exists {
		s.Log.Error(ErrUserAlreadyExists)
		return ErrUserAlreadyExists
	}
	if err != nil {
		s.Log.Error(err)
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	err = s.Repos.AddNewUser(ctx, user)
	if err != nil {
		s.Log.Error(err)
		return err
	}
	err = s.Cache.CacheUser(ctx, user)
	if err != nil {
		s.Log.Error(err)
		return err
	}
	return nil
}

func (s *Service) Authorize(ctx context.Context, user model.Login) error {
	searchUser := model.User{
		Phone:    user.Phone,
		Password: user.Password,
	}
	comparedUser, err := s.Cache.ExistsInCache(ctx, user.Phone)
	if err != nil {
		s.Log.Error(err)
		comparedUser, err = s.Repos.GetUserPhoneAndPasswordByPhone(ctx, searchUser)
		if err != nil {
			s.Log.Error(err)
			s.Log.Error(errors.New(comparedUser.Password + comparedUser.Phone))
			return err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(comparedUser.Password), []byte(searchUser.Password))
	if err != nil {
		s.Log.Error(err)
		return ErrWrongPassword
	}
	return nil
}
