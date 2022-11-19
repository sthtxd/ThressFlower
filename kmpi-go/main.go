package main

import (
	"fmt"
	"kmpi-go/controller"
	"kmpi-go/log"
	"kmpi-go/socket"
	_ "kmpi-go/task"
	"kmpi-go/util"
	"runtime"

	"github.com/fatih/color"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const PORT = "9394"

func main() {
	//调节并发数为cpu核数的2倍
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	startSocket()
	err := initController()
	if err != nil {
		log.Error("initController error", err.Error())
	}

}
func startSocket() {
	go func() {
		socket.StartSocketClient()
	}()

}
func initController() error {
	gin.DefaultWriter = log.BaseGinLog()
	r := gin.New()
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: int(20 * 60), Path: "/"})
	r.Use(sessions.Sessions("mysession", store))
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		return fmt.Sprintf("[GIN] %v |%3d| %13v | %15s | %-7s  %#v %s |\"%s\" \n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
			param.Request.UserAgent(),
		)
	}))
	r.LoadHTMLGlob("./resource/*.html")             // 添加入口index.html
	r.LoadHTMLFiles("./resource/static/*/*")        // 添加资源路径
	r.Static("/static", "./resource/static")        // 添加资源路径
	r.StaticFile("/admin", "./resource/index.html") //前端接口

	r.POST("/api/login", controller.Login)

	v1 := r.Group("/api")
	{
		//v1.Use(midware.SessionMidware())
		v1.GET("/getAllDeviceData", controller.GetAllDeviceData)
		v1.GET("/getTemperature", controller.GetTemperature)
		v1.GET("/getAgvPosition", controller.GetAgvPosition)
		v1.GET("/downloadExcel", controller.DownloadExcel)

		v1.POST("/saveTemperature", controller.SaveFurnaceTemperature)
		v1.POST("/searchLog", controller.SearchFurnaceLog)
		v1.POST("/searchOperationLog", controller.SearchFurnaceOperationLog)
		v1.POST("/saveOperationExcel", controller.SaveOperationExcel)
		v1.POST("/saveExcel", controller.SaveExcel)
		v1.POST("/checkDownloadId", controller.CheckDownloadId)
		v1.POST("/searchLogByType", controller.SearchFurnaceLogByType)
		v1.POST("/searchAgvAndWeight", controller.SearchAgvAndWeight)
	}
	r.POST("/api/adminLogin", controller.AdminLogin)
	v2 := r.Group("/api")
	{
		v2.POST("/updateUser", controller.UpdateUser)
		v2.POST("/addUser", controller.AddUser)
		v2.POST("/deleteUser", controller.DeleteUser)
		v2.POST("/getAllUser", controller.GetAllUser)
		v2.POST("/modifyPassword", controller.ModifyPassword)

	}

	printBanner()
	err := r.Run(":" + PORT) // listen and serve

	return err
	//	r.RunTLS(":9393", "resource/client.pem", "resource/client.key")

}
func printBanner() error {
	localIp, err := util.GetIp()
	if err != nil {
		return err
	}
	color.Yellow(*localIp + ":" + PORT + "/admin")

	return nil
}
