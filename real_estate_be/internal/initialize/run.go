package initialize

func Run() {
	LoadConfig()
	InitMysql()
	InitRouter()
}
