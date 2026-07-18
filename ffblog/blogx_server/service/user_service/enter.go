// service/user_service/enter.go
package user_service

import "blogx_server/models"

type UserService struct {
	userModel models.UserModel
}

func NewUserService(user models.UserModel) *UserService {
	return &UserService{
		userModel: user,
	}
}
