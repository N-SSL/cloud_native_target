package routes

import (
	"github.com/N-SSL/container-target/MySQL"
	"github.com/N-SSL/container-target/controllers"
	"github.com/N-SSL/container-target/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"


	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	var store sessions.Store
	var err error
	store, err = redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret")) //user redis for session
	if err != nil{
		store = cookie.NewStore([]byte("secret"))
	}

	//允许跨域
	r.Use(middlewares.Cors())
	//开启缓存
	r.Use(sessions.Sessions("session_id", store))

	//所有的数据库请求
	r1 := r.Group("/api/k8s")

	//登录注册和测试
	r1.POST("/login",MySQL.JustLogin)
	r1.POST("/register",MySQL.JustRegister)
	r1.GET("/ping/:text", controllers.Ping)

	//开启用户认证，获取环境
	r1.Use(middlewares.UserMiddleware())
	r1.GET("/listAll", MySQL.ListAllImages)
	r1.GET("/getImageByName/:imageName", MySQL.GetImageByName)
	r1.GET("/home",MySQL.GetPersonalInfo)

	//用户对环境进行操作
	r1.POST("/startEnv", MySQL.StartEnv)
	r1.POST("/endEnv", MySQL.EndEnv)
	r1.POST("/restartEnv", MySQL.RestartEnv)

	//管理员管理环境，强行关闭镜像
	r1.Use(middlewares.JudgeIsAdmin())
	r1.POST("/addNewEnv", MySQL.AddNewEnv)
	r1.POST("/deleteEnv", MySQL.DeleteEnv)
	r1.POST("/updateEnv", MySQL.UpdateEnv)
	r1.POST("/endEnvForce", MySQL.EndEnvForce)
	r1.GET("/listAllRunningEnv", MySQL.ListAllRunningEnv)

	//开启镜像
	r2 := r.Group("/app")
	r2.GET("/:id/*path", MySQL.AppHandler)

}
