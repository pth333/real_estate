package initialize

func Run() {
	LoadConfig()
	InitMysql()
	app := InitRouter()
	app.Listen(":8000")
}

func RunCrawler() {
	LoadConfig()
	InitMysql()
}
