package main

import (
	"github.com/kisara71/WeBook/webook/internal/repository"
	"github.com/kisara71/WeBook/webook/internal/repository/cache"
	"github.com/kisara71/WeBook/webook/internal/repository/dao"
	"github.com/kisara71/WeBook/webook/internal/service"
	memory_sms "github.com/kisara71/WeBook/webook/internal/service/sms/memory-sms"
	"github.com/kisara71/WeBook/webook/internal/web/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kisara71/WeBook/webook/internal/web"
)

func main() {

	db := initDB()
	server := initWebserver()
	uh := initUser(db)
	uh.RegisterRoutes(server)

	if err := server.Run(":8080"); err != nil {
		panic(err)
		return
	}
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	if err = dao.InitTable(db); err != nil {
		panic(err)
	}
	return db
}

func initWebserver() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			} else {
				return strings.Contains(origin, "kisara71.xyz")
			}
		},
		MaxAge: 1 * time.Second,
	}))
	//store, err := redis.NewStore(16, "tcp", "localhost:13322", "", "", []byte("d25MZ9waMGelpa9GrTQcawfIeL1YrORY"))
	//store := cookie.NewStore([]byte("secret"))
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("ssid", store))
	//
	//server.Use(middleware.NewLoginMiddlerBuilder().Build())
	server.Use(middleware.NewLoginJwtVerMiddleWare([]string{
		"/users/login",
		"/users/signup",
		"/users/login_sms/code/send",
		"/users/login_sms",
	}).Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:13322",
	})
	client := cache.NewUserCache(redisClient, time.Minute*15)
	userRepo := repository.NewUserRepository(ud, client)
	userSvc := service.NewUserService(userRepo)

	codeCache := cache.NewCodeCache(redisClient)
	smsRepo := repository.NewCodeRepository(codeCache)
	//smsSvc, _ := aliyun_sms.NewService("dysmsapi.aliyuncs.com")
	smsSvc := memory_sms.NewService()
	codeSvc := service.NewCodeService(smsRepo, smsSvc)
	uhr := web.InitUserHandler(userSvc, codeSvc)
	return uhr
}
