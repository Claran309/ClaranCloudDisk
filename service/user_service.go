package services

import (
	"ClaranCloudDisk/dao/mysql"
	"ClaranCloudDisk/model"
	"ClaranCloudDisk/util"
	"ClaranCloudDisk/util/jwt_util"
	"ClaranCloudDisk/util/minIO"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type UserService struct {
	UserRepo    mysql.UserRepository
	TokenRepo   mysql.TokenRepository
	jwtUtil     jwt_util.Util
	AvatarDIR   string
	minioClient *minIO.MinIOClient
}

func NewUserService(userRepo mysql.UserRepository, tokenRepo mysql.TokenRepository, jwtUtil jwt_util.Util, avatarDIR string, minioClient *minIO.MinIOClient) *UserService {
	return &UserService{
		UserRepo:    userRepo,
		TokenRepo:   tokenRepo,
		jwtUtil:     jwtUtil,
		AvatarDIR:   avatarDIR,
		minioClient: minioClient,
	}
}

func (s *UserService) Register(req *model.RegisterRequest) (*model.User, *model.InvitationCode, error) {
	//邀请码是否正确
	invitationCode, err := s.UserRepo.ValidateInvitationCode(req.InviteCode)
	if err != nil {
		return nil, nil, errors.New("邀请码不正确" + err.Error())
	}

	//密码时候否符合格式
	var flagPassword bool
	for i := 0; i < len(req.Password); i++ {
		if !((req.Password[i] >= 'a' && req.Password[i] <= 'z') || (req.Password[i] >= '0' && req.Password[i] <= '9') || (req.Password[i] >= 'A' && req.Password[i] <= 'Z')) {
			flagPassword = true
		}
	}
	if flagPassword {
		return nil, nil, errors.New("password format Error")
	}

	//邮箱是否符合格式
	if !strings.Contains(req.Email, "@") {
		return nil, nil, errors.New("email format Error")
	}

	//加密密码
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, nil, errors.New("password hash failed")
	}

	//创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword, // 加密存储
		Email:    req.Email,
		Role:     "user",
		Storage:  0,
	}

	// 如果是服务器用户
	if invitationCode.Code == "FirstAdminCode" && invitationCode.CreatorUserID == 1 && invitationCode.IsUsed == false && invitationCode.UserID == 1 {
		user.Role = "admin"
	}

	//传入数据库
	if err := s.UserRepo.AddUser(user); err != nil {
		return nil, nil, err
	}

	user, _ = s.UserRepo.SelectByUsername(user.Username)

	//使用邀请码
	err = s.UserRepo.UseInvitationCode(invitationCode.Code, user.UserID)
	if err != nil {
		return nil, nil, errors.New("使用邀请码时出现错误")
	}

	return user, &invitationCode, nil
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

	user, _ = s.UserRepo.SelectByUsername(user.Username)
	if user.IsBanned == true {
		return "", nil, errors.New("用户已被封禁"), ""
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

	var userID int
	if userIDFloat, ok := claims["user_id"].(float64); ok {
		// 安全转换：float64 转 int
		userID = int(userIDFloat)
	}
	newToken, err := s.jwtUtil.GenerateToken(userID, claims["username"].(string), claims["role"].(string), 1)
	if err != nil {
		return "", errors.New("token generate failed")
	}

	return newToken, nil
}

func (s *UserService) CheckStorage(UserID int) (int64, error) {
	UsedStorage, err := s.UserRepo.GetStorage(UserID)
	if err != nil {
		return -1, errors.New("get storage failed")
	}
	return UsedStorage, nil
}

func (s *UserService) Logout(token string) error {
	//加入token黑名单
	err := s.TokenRepo.AddBlackList(token)
	if err != nil {
		return errors.New("add black list failed: " + err.Error())
	}
	return nil
}

func (s *UserService) UpdateInfo(UserID int, req model.UpdateRequest) (model.User, error) {
	if req.Username != "" {
		err := s.UserRepo.UpdateUsername(UserID, req.Username)
		if err != nil {
			return model.User{}, err
		}
	}

	if req.Email != "" {
		//邮箱是否符合格式
		if !strings.Contains(req.Email, "@") {
			return model.User{}, errors.New("email format Error")
		}
		err := s.UserRepo.UpdateEmail(UserID, req.Email)
		if err != nil {
			return model.User{}, err
		}
	}

	if req.Role != "" {
		err := s.UserRepo.UpdateRole(UserID, req.Role)
		if err != nil {
			return model.User{}, err
		}
	}

	if req.Password != "" {
		//密码时候否符合格式
		var flagPassword bool
		for i := 0; i < len(req.Password); i++ {
			if !((req.Password[i] >= 'a' && req.Password[i] <= 'z') || (req.Password[i] >= '0' && req.Password[i] <= '9') || (req.Password[i] >= 'A' && req.Password[i] <= 'Z')) {
				flagPassword = true
			}
		}
		if flagPassword {
			return model.User{}, errors.New("password format Error")
		}
		//加密密码
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			return model.User{}, errors.New("password hash failed")
		}
		err = s.UserRepo.UpdatePassword(UserID, hashedPassword)
		if err != nil {
			return model.User{}, err
		}
	}

	user, err := s.UserRepo.SelectByUserID(UserID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) GenerateInvitationCode(userID int) (model.InvitationCode, error) {
	user, err := s.UserRepo.SelectByUserID(userID)
	if err != nil {
		return model.InvitationCode{}, errors.New("user not found")
	}
	//判定是否超出非admin生成数量
	if user.Role != "admin" {
		//VIP可生成三个
		if user.IsVIP == true {
			if user.GeneratedInvitationCodeNum >= 3 {
				return model.InvitationCode{}, errors.New("非管理员的VIP用户邀请码生成数量不能超过3个")
			}
		}
		//非VIP只能生成一个
		if user.IsVIP == false {
			if user.GeneratedInvitationCodeNum >= 1 {
				return model.InvitationCode{}, errors.New("非管理员的普通用户邀请码生成数量不能超过3个")
			}
		}
	}

	//生成邀请码
	charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	inviteCode := make([]byte, 10)
	for i := 0; i < 10; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(62))
		if err != nil {
			return model.InvitationCode{}, errors.New("rand error")
		}
		inviteCode[i] = charset[n.Int64()]
	}

	var invitationCode = model.InvitationCode{
		Code:          string(inviteCode),
		CreatorUserID: userID,
	}

	//数据层
	//保存邀请码
	err = s.UserRepo.CreateInvitationCode(invitationCode)
	if err != nil {
		return model.InvitationCode{}, errors.New("create invitation code failed")
	}
	//更新用户邀请码数量
	err = s.UserRepo.AddInvitationCodeNum(userID)
	if err != nil {
		return model.InvitationCode{}, errors.New("add invitation code failed")
	}

	return invitationCode, nil
}

func (s *UserService) InvitationCodeList(userID int) ([]model.InvitationCode, int64, error) {
	invitationCodes, total, err := s.UserRepo.GetInvitationCodeList(userID)
	if err != nil {
		return nil, -1, errors.New("get invitation code list failed")
	}
	return invitationCodes, total, nil
}

func (s *UserService) UploadAvatar(file *multipart.FileHeader, userID int, userName string) (string, string, string, error) {
	// 验证文件类型
	contentType := file.Header.Get("Content-Type")
	allowedTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/jpg",
	}

	validType := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			validType = true
			break
		}
	}

	if !validType {
		return "", "", "", errors.New("不支持的文件类型，请上传图片文件 (jpg, png, gif, webp)")
	}

	// 限制5MB
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		return "", "", "", errors.New("文件大小不能超过5MB")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", "", "", errors.New("无法打开文件")
	}
	defer src.Close()

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// 如果没有扩展名，根据Content-Type推断
		exts, _ := mime.ExtensionsByType(contentType)
		if len(exts) > 0 {
			ext = exts[0]
		} else {
			ext = ".jpg" // 默认
		}
	}

	// 使用username+UID生成唯一文件名
	fileName := fmt.Sprintf("avatar_%d_%s%s", userID, userName, ext)

	//目录
	userDir := filepath.Join(s.AvatarDIR, fmt.Sprintf("user_%d", userID))

	//获取字节数组
	bytes, err := io.ReadAll(src)
	if err != nil {
		return "", "", "", errors.New("读取文件失败")
	}

	// 创建目标文件
	dstPath := filepath.Join(userDir, fileName)

	//保存文件到minIO
	err = s.minioClient.Save(context.Background(), dstPath, bytes, ext)
	if err != nil {
		return "", "", "", err
	}

	////============================================非minIO=============================================================
	//// 确保目录存在
	//if err := os.MkdirAll(userDir, 0755); err != nil {
	//	return "", "", "", errors.New("创建存储目录失败")
	//}
	//
	//
	//dst, err := os.Create(dstPath)
	//if err != nil {
	//	return "", "", "", errors.New("创建目标文件失败")
	//}
	//defer dst.Close()
	//// 创建目标文件
	//	dstPath := filepath.Join(userDir, fileName)
	//// 复制文件内容到服务器
	//_, err = io.Copy(dst, src)
	//if err != nil {
	//	// 删除不完整的文件
	//	os.Remove(dstPath)
	//	return "", "", "", errors.New("保存文件失败")
	//}
	////================================================================================================================

	// 构建访问URL
	// 因为这个URL是新上传的头像，所以不管之前有没有上传过头像都没关系
	//avatarURL := filepath.Join(s.AvatarDIR, fmt.Sprintf("/user_%d/%s", userID, fileName))

	// 将访问URL保存到数据层
	err = s.UserRepo.UploadAvatar(userID /*avatarURL*/, dstPath)
	if err != nil {
		return "", "", "", errors.New("upload database failed: " + err.Error())
	}

	return dstPath, fileName, contentType, nil
}

func (s *UserService) GetAvatar(userID int) (string, error) {
	//访问数据层
	avatarPath, err := s.UserRepo.GetAvatar(userID)
	if err != nil {
		return "", errors.New("get avatar failed")
	}
	return avatarPath, nil
}
