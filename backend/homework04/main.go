package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"homework04/config"
	"homework04/handlers"
	"homework04/middleware"
	"homework04/models"
	"homework04/services"
	"homework04/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 数据库初始化
func InitDB() (*gorm.DB, error) {

	var db *gorm.DB
	var err error

	switch config.AppConfig.DBType {

	case "sqlite":
		os.MkdirAll("./db", os.ModePerm)

		db, err = gorm.Open(sqlite.Open(config.AppConfig.SqlitePath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})

	case "mysql":

		db, err = gorm.Open(mysql.Open(config.AppConfig.MySQLDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})

	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.AppConfig.DBType)

	}

	return db, err
}

var db *gorm.DB

func main() {
	utils.InitLogger()
	// 定义命令行参数
	portFlag := flag.String("port", "", "Port to run the server on")
	dbFlag := flag.String("db", "", "database type")
	flag.Parse()
	utils.InfoLogger.Printf("启动参数: port=%s, db=%s\n", *portFlag, *dbFlag)

	//读取配置文件
	err := config.LoadConfig("config/config.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("端口:", config.AppConfig.Port)

	//启动命令行参数优先
	port := config.AppConfig.Port
	if *portFlag != "" {
		port = *portFlag
	}
	dbType := config.AppConfig.DBType

	if *dbFlag != "" {
		dbType = *dbFlag
	}
	config.AppConfig.DBType = dbType

	db, err = InitDB()
	if err != nil {
		utils.ErrorLogger.Fatalf("init db failed: %v", err)
		panic(err)
	}
	// 自动迁移模型
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	// 初始化服务
	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService, []byte(config.AppConfig.JWTSecret))
	postService := services.NewPostService(db)
	postHandlers := handlers.NewPostHandler(postService)
	commentService := services.NewCommentService(db)
	commentHandler := handlers.NewCommentHandler(commentService)

	//启动 Gin 服务
	r := gin.Default()
	//统一通过这个中间件处理异常recovery.
	r.Use(middleware.RecoveryMiddleware())

	//不需要 token 的接口组
	publicGroup := r.Group("/open/api")
	{
		// 测试接口
		publicGroup.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, utils.Response{
				Code:    200,
				Message: "Hello World!---open",
			})
		})

		// 用户注册
		publicGroup.POST("/user/register", userHandler.Register)

		// 用户登录
		publicGroup.POST("/user/login", userHandler.Login)

		//获取文章列表或详情
		publicGroup.GET("/posts", postHandlers.GetPosts)
		publicGroup.GET("/posts/:id", postHandlers.GetPost)
		//获取某篇文章评论列表
		publicGroup.GET("/posts/comments/:post_id", commentHandler.ListComments)
	}

	//需要 token 的接口组
	authGroup := r.Group("auth/api")
	authGroup.Use(middleware.JWTAuthMiddleware()) // 中间件验证 JWT
	{
		authGroup.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, utils.Response{
				Code:    200,
				Message: "Hello World!---auth",
			})
		})
		//创建/更新/删除文章
		authGroup.POST("/post", postHandlers.CreatePost)
		authGroup.PUT("/post/:id", postHandlers.UpdatePost)
		authGroup.DELETE("/post/:id", postHandlers.DeletePost)
		//创建评论
		authGroup.POST("/posts/comments/:post_id", commentHandler.CreateComment)
	}

	utils.InfoLogger.Printf("启动服务器: 端口=%s, 数据库=%s\n", port, config.AppConfig.DBType)
	r.Run(":" + port)
}
