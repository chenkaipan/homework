package services

import (
	"homework04/models"
	"homework04/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}
func (s *UserService) Register(req models.CreateUserRequest) (uint, error) {

	var count int64
	s.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)

	if count > 0 {
		return 0, utils.NewBizError(400, "username already exists")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return 0, utils.NewBizError(500, "failed to create user")
	}

	return user.ID, nil
}

func (s *UserService) Login(req models.LoginRequest) (string, error) {
	var user models.User
	user.Username = req.Username
	user.Password = req.Password
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return "", utils.NewBizError(400, "username or password incorrect")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", utils.NewBizError(400, "username or password incorrect")
	}
	utils.InfoLogger.Println("登录")
	utils.InfoLogger.Println(user.ID)
	utils.InfoLogger.Println(user.Username)
	token, _ := utils.GenerateToken(user.ID, user.Username)

	return token, nil
}
