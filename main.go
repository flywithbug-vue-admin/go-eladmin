package main

import (
	"flag"
	"fmt"
	"go-eladmin/config"
	"go-eladmin/core/jwt"
	"go-eladmin/core/mongo"
	"go-eladmin/email"
	"go-eladmin/log_writer"
	mongoIndex "go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"go-eladmin/server"
	"go-eladmin/server/handler/file_handler"
	"os"

	log "github.com/flywithbug/log4go"
)

//log 启动配置
func setLog() {
	//log日志写入文件
	//w := log.NewFileWriter()
	//w.SetPathPattern(config.Conf().LogPath)
	//log.Register(w)

	//log日志控制台输出
	//c := log.NewConsoleWriter()
	//c.SetColor(true)
	//log.Register(c)

	fmt.Println("init log")
	//日志保存到db
	w := log_writer.NewDBWriter()
	log.Register(w)

	//log日志控制台输出
	l := log_writer.NewConsoleWriter()
	l.SetColor(true)
	log.Register(l)

	log.SetLevel(1)
	log.SetLayout("2006-01-02 15:04:05")
	fmt.Println("\tcomplete")

}

func setFileConfig() {
	fmt.Println("设置图片文件存储路径")
	file_handler.SetLocalImageFilePath(config.Conf().LocalFilePath)
}

func setJWTKey() {
	//signingKey read
	fmt.Println("配置JWT 秘钥文件")
	jwt.ReadSigningKey(config.Conf().PrivateKeyPath, config.Conf().PublicKeyPath)
}

func init() {
	fmt.Println("----------------init----------------")
	//配置文件
	configPath := flag.String("config", "config.json", "Configuration file to use")
	flag.Parse()
	fmt.Println("ParseConfig")
	err := config.ReadConfig(*configPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("\tComplete")
}

func setMongoDB() {
	//mongodb启动连接
	//设置数据库名字
	fmt.Println("init mongodb")
	shareDB.SetDocMangerDBName(config.Conf().DBConfig.DBName)
	shareDB.SetMonitorDBName(config.Conf().MonitorDBConfig.DBName)

	fmt.Println("\tdoc_manager")
	err := mongo.RegisterMongo(config.Conf().DBConfig.Url, config.Conf().DBConfig.DBName)
	if err != nil {
		panic(err)
	}
	fmt.Println("\tmonitor")
	err = mongo.RegisterMongo(config.Conf().MonitorDBConfig.Url, config.Conf().MonitorDBConfig.DBName)
	if err != nil {
		panic(err)
	}
	fmt.Println("\tindex")

	//模型唯一索引
	mongoIndex.CreateMgoIndex()
	fmt.Println("\tcomplete")

}

func setMail() {
	var err error
	email.Mail, err = config.Conf().MailConfig.Dialer()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(0)
	}
}

func main() {
	setMongoDB()

	//log配置
	setLog()
	defer log.Close()
	//文件存储位置
	setFileConfig()

	setMail()

	//jwt验证
	setJWTKey()
	fmt.Println("Server")
	//启动ApiServer服务
	server.StartServer(config.Conf().Port, config.Conf().StaticPath, config.Conf().RouterPrefix, config.Conf().AuthPrefix)
}
