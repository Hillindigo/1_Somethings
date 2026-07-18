package service

import (
	"errors"

	"github.com/goagent/mojian/internal/config"
	"github.com/goagent/mojian/internal/model"
	"github.com/goagent/mojian/internal/repository"
	"github.com/goagent/mojian/pkg/errcode"
	"github.com/goagent/mojian/pkg/utils"
	"gorm.io/gorm"
)

// UserService 用户业务逻辑层
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 创建用户 Service 实例
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register 用户注册，校验用户名和邮箱唯一性后创建用户
func (s *UserService) Register(req *model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrUserAlreadyExists))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 检查邮箱是否已存在
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrEmailAlreadyUsed))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 密码哈希加密
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrUserCreateFailed))
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: hash,
		Email:        req.Email,
		Role:         0, // 默认普通用户
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrUserCreateFailed))
	}

	return user, nil
}

// Login 用户登录，校验密码后生成 JWT Token
func (s *UserService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrUserNotFound))
		}
		return nil, err
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New(errcode.GetMessage(errcode.ErrPasswordIncorrect))
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Role, config.GlobalConfig.JWT.Secret, config.GlobalConfig.JWT.Expire)
	if err != nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrTokenGenerate))
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetUser 根据用户 ID 获取用户信息
func (s *UserService) GetUser(id uint) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrUserNotFound))
		}
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New(errcode.GetMessage(errcode.ErrUserNotFound))
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// ListUsers 获取用户列表（分页），管理员使用
func (s *UserService) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize)
}

// UpdateUserRole 更新用户角色，管理员使用
func (s *UserService) UpdateUserRole(id uint, role int8) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(errcode.GetMessage(errcode.ErrUserNotFound))
		}
		return nil, err
	}

	user.Role = role
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser 删除用户，管理员使用
func (s *UserService) DeleteUser(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(errcode.GetMessage(errcode.ErrUserNotFound))
		}
		return err
	}
	return s.repo.Delete(id)
}
