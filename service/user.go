package services

import (
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/model"
	"ClaranCloudDisk/util"
	"ClaranCloudDisk/util/jwt_util"
	"errors"
	"fmt"
	"strings"
)

type UserService struct {
	UserRepo mysql.UserRepository
	jwtUtil  jwt_util.Util
}

func NewUserService(userRepo mysql.UserRepository, jwtUtil jwt_util.Util) *UserService {
	return &UserService{
		UserRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (s *UserService) Register(req *model.RegisterRequest) (*model.User, error) {
	//密码时候否符合格式
	var flagPassword bool
	for i := 0; i < len(req.Password); i++ {
		if !((req.Password[i] >= 'a' && req.Password[i] <= 'z') || (req.Password[i] >= '0' && req.Password[i] <= '9') || (req.Password[i] >= 'A' && req.Password[i] <= 'Z')) {
			flagPassword = true
		}
	}
	if flagPassword {
		return nil, errors.New("password format Error")
	}

	//邮箱是否符合格式
	if !strings.Contains(req.Email, "@") {
		return nil, errors.New("email format Error")
	}

	//加密密码
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("password hash failed")
	}

	//创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword, // 加密存储
		Email:    req.Email,
		Role:     req.Role,
	}

	//传入数据库
	if err := s.UserRepo.AddUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(loginKey, password string) (string, *model.User, error, string) {
	//判断是邮箱登录还是用户名登录
	var user *model.User
	var at, point bool
	for i := 0; i < len(loginKey); i++ {
		if loginKey[i] == '@' {
			at = true
		}
		if loginKey[i] == '.' {
			point = true
		}
	}
	if at && point { // 邮箱登录
		userByEmail, err := s.UserRepo.SelectByEmail(loginKey)
		if err != nil {
			return "", nil, err, ""
		}
		user = userByEmail
	} else { // 用户名登录
		userByUsername, err := s.UserRepo.SelectByUsername(loginKey)
		if err != nil {
			return "", nil, err, ""
		}
		user = userByUsername
	}

	//检查用户是否存在
	if flag := s.UserRepo.Exists(user.Username, user.Email); !flag {
		return "", nil, errors.New("username Not Exist"), ""
	}

	//检验密码正确性
	if !util.CheckPassword(user.Password, password) {
		return "", nil, errors.New("password error"), ""
	}

	//access token
	token, err := s.jwtUtil.GenerateToken(user.UserID, user.Username, user.Role, 1)
	//refresh token
	refreshToken, err := s.jwtUtil.GenerateToken(user.UserID, user.Username, user.Role, 255)
	if err != nil {
		return "", nil, errors.New("token Error"), ""
	}

	return token, user, nil, refreshToken
}

func (s *UserService) Refresh(refreshToken model.RefreshTokenRequest) (string, error) {
	//验证token
	token, err := s.jwtUtil.ValidateToken(refreshToken.RefreshToken)
	if err != nil {
		return "", errors.New("refresh Token Error")
	}

	//提取声明
	claims, err := s.jwtUtil.ExtractClaims(token)
	if err != nil {
		return "", errors.New("extract claims failed")
	}

	newToken, err := s.jwtUtil.GenerateToken(claims["user_id"].(int), claims["username"].(string), claims["role"].(string), 1)
	if err != nil {
		return "", errors.New("token generate failed")
	}

	return newToken, nil
}

func (s *UserService) CheckStorage(UserID int) (string, error) {
	UsedStorage, err := s.UserRepo.GetStorage(UserID)
	if err != nil {
		return "", errors.New("get storage failed")
	}
	return fmt.Sprintf("%dG", UsedStorage), nil
}
