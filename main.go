package main

import (
	"PracticeProject/router"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
func main() {
	// TODO 后续把一些配置抽离出来，根据环境来加载不同的配置(dev/paas)，可通过命令行参数来表示不同的环境

	r := router.Router()
	r.Run() // 监听并在 0.0.0.0:8080上启动服务
}
