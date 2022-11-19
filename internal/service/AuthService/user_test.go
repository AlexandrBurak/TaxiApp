package AuthService

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/zxcghoulhunter/InnoTaxi/internal/mock"
	"github.com/zxcghoulhunter/InnoTaxi/internal/model"
)

func TestService_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()
	logger := mock.NewMockLogger(ctrl)
	repos := mock.NewMockRepository(ctrl)
	cache := mock.NewMockCache(ctrl)
	user := model.User{Name: "adsadasd", Phone: "123123123", Email: "123@qawe", Password: "zxczxccz"}
	repos.EXPECT().Exists(gomock.Any(), user).Return(false, nil)
	repos.EXPECT().AddNewUser(gomock.Any(), gomock.Any()).Return(nil)
	cache.EXPECT().CacheUser(gomock.Any(), gomock.Any()).Return(nil)
	service := Service{repos, logger, cache}

	err := service.Save(ctx, user)
	if err != nil {
		t.Error("")
	}
}

func TestService_Save_UserExistsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()
	repos := mock.NewMockRepository(ctrl)
	logger := mock.NewMockLogger(ctrl)
	cache := mock.NewMockCache(ctrl)
	user := model.User{Name: "adsadasd", Phone: "123123123", Email: "123@qawe", Password: "zxczxccz"}
	repos.EXPECT().Exists(gomock.Any(), user).Return(true, nil)
	logger.EXPECT().Error(gomock.Any()).Return()

	service := Service{repos, logger, cache}

	err := service.Save(ctx, user)
	if err != ErrUserAlreadyExists {
		t.Error("")
	}
}

func TestService_Save_UnhandledError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()
	repos := mock.NewMockRepository(ctrl)
	logger := mock.NewMockLogger(ctrl)
	cache := mock.NewMockCache(ctrl)
	user := model.User{Name: "adsadasd", Phone: "123123123", Email: "123@qawe", Password: "zxczxccz"}
	repos.EXPECT().Exists(gomock.Any(), user).Return(false, errors.New("zxzc"))
	logger.EXPECT().Error(gomock.Any()).Return()

	service := Service{repos, logger, cache}

	err := service.Save(ctx, user)
	if err == nil {
		t.Error("")
	}
}

func TestService_Authorize_UnhandledError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()
	repos := mock.NewMockRepository(ctrl)
	logger := mock.NewMockLogger(ctrl)
	cache := mock.NewMockCache(ctrl)
	user := model.User{Phone: "123123123", Password: "2345w891"}
	cache.EXPECT().ExistsInCache(gomock.Any(), gomock.Any()).Return(user, nil)
	//repos.EXPECT().GetUserPhoneAndPasswordByPhone(gomock.Any(), user).Return(model.User{Phone: "123123123", Password: ""}, nil)
	logger.EXPECT().Error(gomock.Any()).Return()
	service := Service{repos, logger, cache}

	err := service.Authorize(ctx, model.Login{Phone: user.Phone, Password: "2345w891"})
	if err == nil {
		t.Error("")
	}
}

func TestService_Authorize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()
	repos := mock.NewMockRepository(ctrl)
	logger := mock.NewMockLogger(ctrl)
	cache := mock.NewMockCache(ctrl)
	user := model.User{Phone: "123123123", Password: "2345w891"}
	//repos.EXPECT().GetUserPhoneAndPasswordByPhone(gomock.Any(), user).Return(model.User{Phone: "123123123", Password: "$2a$04$ddAHC4UK4JF0iId7AZR3wOccM4d.sasMvX0Ldsr/.ZrZXYsj7CZeW"}, nil)
	service := Service{repos, logger, cache}
	cache.EXPECT().ExistsInCache(gomock.Any(), gomock.Any()).Return(model.User{Phone: "123123123", Password: "$2a$04$ddAHC4UK4JF0iId7AZR3wOccM4d.sasMvX0Ldsr/.ZrZXYsj7CZeW"}, nil)
	err := service.Authorize(ctx, model.Login{Phone: user.Phone, Password: "2345w891"})
	if err != nil {
		t.Error("")
	}
}

func TestService_Authorize_Wrong_Password(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()
	repos := mock.NewMockRepository(ctrl)
	cache := mock.NewMockCache(ctrl)
	logger := mock.NewMockLogger(ctrl)
	user := model.User{Phone: "123123123", Password: "2345w892"}
	repos.EXPECT().GetUserPhoneAndPasswordByPhone(gomock.Any(), user).Return(model.User{Phone: "123123123", Password: "$2a$04$ddAHC4UK4JF0iId7AZR3wOccM4d.sasMvX0Ldsr/.ZrZXYsj7CZeW"}, nil)
	logger.EXPECT().Error(gomock.Any()).Return().AnyTimes()
	cache.EXPECT().ExistsInCache(gomock.Any(), gomock.Any()).Return(user, errors.New("asd"))
	service := Service{repos, logger, cache}

	err := service.Authorize(ctx, model.Login{Phone: user.Phone, Password: "2345w892"})
	if err != ErrWrongPassword {
		t.Error("")
	}

}
