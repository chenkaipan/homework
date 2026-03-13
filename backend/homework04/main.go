package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Config struct {
	Port       string `json:"port"`
	DBType     string `json:"db_type"`
	SqlitePath string `json:"sqlite_path"`
	MySQLDSN   string `json:"mysql_dsn"`
	JWTSecret  string `json:"jwt_secret"`
}
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}
type BizError struct {
	Code int
	Msg  string
}

func (e *BizError) Error() string {
	return e.Msg
}

// 工厂函数
func NewBizError(code int, msg string) *BizError {
	return &BizError{
		Code: code,
		Msg:  msg,
	}
}

// 生成 JWT
func GenerateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 24 小时过期
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 验证 JWT
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
}

// JWT 中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "Authorization header missing",
			})
			c.Abort()
			return
		}

		token, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "invalid token claims",
			})
			c.Abort()
			return
		}

		// 保存用户信息到上下文
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))

		c.Next()
	}
}

// 全局异常处理
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			if bizErr, ok := err.Err.(*BizError); ok {
				// 记录业务错误
				errorLogger.Printf("BizError: code=%d, msg=%s\n", bizErr.Code, bizErr.Msg)
				// 返回 HTTP 状态码
				httpStatus := 400
				if bizErr.Code >= 500 {
					httpStatus = 500
				} else if bizErr.Code == 401 {
					httpStatus = 401
				} else if bizErr.Code == 403 {
					httpStatus = 403
				} else if bizErr.Code == 404 {
					httpStatus = 404
				}

				c.JSON(httpStatus, Response{
					Code:    bizErr.Code,
					Message: bizErr.Msg,
				})
			} else {
				// 未知错误
				errorLogger.Printf("UnknownError: %v\n", err.Err)
				c.JSON(http.StatusInternalServerError, Response{
					Code:    500,
					Message: "internal server error",
				})
			}
			c.Abort()
		}
	}
}

// 数据库初始化
func InitDB(config Config) (*gorm.DB, error) {

	var db *gorm.DB
	var err error

	switch config.DBType {

	case "sqlite":
		os.MkdirAll("./db", os.ModePerm)

		db, err = gorm.Open(sqlite.Open(config.SqlitePath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})

	case "mysql":

		db, err = gorm.Open(mysql.Open(config.MySQLDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})

	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.DBType)

	}

	return db, err
}

// 注册函数
func RegisterHandler(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewBizError(400, err.Error()))
		return
	}

	var count int64
	db.Model(&User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.Error(NewBizError(400, "username already exists"))
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := db.Create(&user).Error; err != nil {
		c.Error(NewBizError(500, "failed to create user"))
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    gin.H{"user_id": user.ID},
	})
}

// 登录函数
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewBizError(400, err.Error()))
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.Error(NewBizError(400, "username or password incorrect"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.Error(NewBizError(400, "username or password incorrect"))
		return
	}

	token, _ := GenerateToken(user.ID, user.Username)
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "login success",
		Data:    gin.H{"token": token},
	})
}

// 创建文章
func CreatePostHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewBizError(400, err.Error()))
		return
	}

	post := Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := db.Create(&post).Error; err != nil {
		c.Error(NewBizError(500, "failed to create post"))
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "post created",
		Data:    gin.H{"post_id": post.ID},
	})
}

// 获取文章列表
func ListPostsHandler(c *gin.Context) {
	var posts []Post
	if err := db.Preload("User").Find(&posts).Error; err != nil {
		c.Error(NewBizError(500, "failed to list posts"))
		return
	}

	data := make([]gin.H, len(posts))
	for i, p := range posts {
		data[i] = gin.H{
			"id":      p.ID,
			"title":   p.Title,
			"content": p.Content,
			"author":  p.User.Username,
			"created": p.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    data,
	})
}

// 获取单篇文章
func GetPostHandler(c *gin.Context) {
	id := c.Param("id")
	var post Post
	if err := db.Preload("User").First(&post, id).Error; err != nil {
		c.Error(NewBizError(500, "post not found"))
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data: gin.H{
			"id":      post.ID,
			"title":   post.Title,
			"content": post.Content,
			"author":  post.User.Username,
			"created": post.CreatedAt,
		},
	})
}

// 更新文章
func UpdatePostHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var post Post
	if err := db.First(&post, id).Error; err != nil {
		c.Error(NewBizError(404, "post not found"))
		return
	}

	if post.UserID != userID {
		c.Error(NewBizError(403, "forbidden: not the author"))
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewBizError(400, err.Error()))
		return
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := db.Save(&post).Error; err != nil {
		c.Error(NewBizError(500, "failed to update post"))
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "post updated",
	})
}

// 删除文章
func DeletePostHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var post Post
	if err := db.First(&post, id).Error; err != nil {
		c.Error(NewBizError(404, "post not found"))
		return
	}

	if post.UserID != userID {
		c.Error(NewBizError(403, "forbidden: not the author"))
		return
	}

	if err := db.Delete(&post).Error; err != nil {
		c.Error(NewBizError(500, "failed to delete post"))
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "post deleted",
	})
}

// 创建评论
func CreateCommentHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	postID := c.Param("post_id")

	var post Post
	if err := db.First(&post, postID).Error; err != nil {
		c.Error(NewBizError(404, "post not found"))
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewBizError(400, err.Error()))
		return
	}

	comment := Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  post.ID,
	}

	if err := db.Create(&comment).Error; err != nil {
		c.Error(NewBizError(500, "failed to create comment"))
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "comment created",
		Data:    gin.H{"comment_id": comment.ID},
	})
}

// 获取某篇文章的评论列表
func ListCommentsHandler(c *gin.Context) {
	postID := c.Param("post_id")

	var comments []Comment
	if err := db.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		c.Error(NewBizError(500, "failed to list comments"))
		return
	}

	data := make([]gin.H, len(comments))
	for i, cm := range comments {
		data[i] = gin.H{
			"id":      cm.ID,
			"content": cm.Content,
			"author":  cm.User.Username,
			"created": cm.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    data,
	})
}

var jwtSecret []byte

var db *gorm.DB

var (
	infoLogger  = log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "ERROR\t", log.LstdFlags)
)

func main() {
	// 定义命令行参数
	portFlag := flag.String("port", "", "Port to run the server on")
	dbFlag := flag.String("db", "", "database type")
	flag.Parse()
	infoLogger.Printf("启动参数: port=%s, db=%s\n", *portFlag, *dbFlag)

	//读取配置文件
	config := Config{Port: "8080"} // 默认端口
	file, err := os.Open("config.json")

	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			infoLogger.Println("读取配置文件失败，使用默认端口 8080")
		}
	} else {
		infoLogger.Println("未找到配置文件，使用默认端口 8080")
	}

	if config.JWTSecret != "" {
		jwtSecret = []byte(config.JWTSecret)
	} else {
		infoLogger.Println("未在配置文件中找到 JWT 密钥，将使用默认值")
		jwtSecret = []byte("default_secret_key")
	}

	//启动命令行参数优先
	port := config.Port
	if *portFlag != "" {
		port = *portFlag
	}
	dbType := config.DBType

	if *dbFlag != "" {
		dbType = *dbFlag
	}
	config.DBType = dbType

	db, err = InitDB(config)
	if err != nil {
		errorLogger.Println(err)
		panic(err)
	}
	// 自动迁移模型
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	//启动 Gin 服务
	r := gin.Default()
	//统一通过这个中间件处理异常
	r.Use(RecoveryMiddleware())

	//不需要 token 的接口组
	publicGroup := r.Group("/open/api")
	{
		// 测试接口
		publicGroup.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, Response{
				Code:    200,
				Message: "Hello World!---open",
			})
		})

		// 用户注册
		publicGroup.POST("/user/register", RegisterHandler)

		// 用户登录
		publicGroup.POST("/user/login", LoginHandler)

		//获取文章列表或详情
		publicGroup.GET("/posts", ListPostsHandler)
		publicGroup.GET("/posts/:id", GetPostHandler)
		//获取某篇文章评论列表
		publicGroup.GET("/posts/comments/:post_id", ListCommentsHandler)
	}

	//需要 token 的接口组
	authGroup := r.Group("auth/api")
	authGroup.Use(JWTAuthMiddleware()) // 中间件验证 JWT
	{
		authGroup.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, Response{
				Code:    200,
				Message: "Hello World!---auth",
			})
		})
		//创建/更新/删除文章
		authGroup.POST("/post", CreatePostHandler)
		authGroup.PUT("/post/:id", UpdatePostHandler)
		authGroup.DELETE("/post/:id", DeletePostHandler)
		//创建评论
		authGroup.POST("/posts/comments/:post_id", CreateCommentHandler)
	}

	infoLogger.Printf("启动服务器: 端口=%s, 数据库=%s\n", port, config.DBType)
	r.Run(":" + port)
}
