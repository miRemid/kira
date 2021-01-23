package main

import (
	"fmt"
	"log"
	"time"

	hystrixGo "github.com/afex/hystrix-go/hystrix"
	"github.com/casbin/casbin/v2"
	"github.com/miRemid/kira/common"

	"github.com/miRemid/kira/services/user/handler"
	"github.com/miRemid/kira/services/user/pb"
	"github.com/miRemid/kira/services/user/repository"
	"github.com/miRemid/kira/services/user/route"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getConnect() string {
	username := common.Getenv("MYSQL_USERNAME", "shi")
	password := common.Getenv("MYSQL_PASSWORD", "123456")
	database := common.Getenv("MYSQL_DATABASE", "kira")
	address := common.Getenv("MYSQL_ADDRESS", "127.0.0.1:3306")
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&Local", username, password, address, database)
}

func connect() (*gorm.DB, error) {
	connect := getConnect()
	_conn, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connect,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db, _ := _conn.DB()
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return _conn, nil
}

func startAPIService() {
	e, err := casbin.NewEnforcer("./casbin/model.conf", "./casbin/permission.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := route.Route(e)
	service := web.NewService(
		web.Name("go.micro.api.user"),
		web.Address(common.Getenv("API_ADDRESS", ":5002")),
		web.Handler(r),
		web.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
	)
	route.Init(client.DefaultClient)
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func startMicroService() {
	service := micro.NewService(
		micro.Name("kira.micro.service.user"),
		micro.Version("latest"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(common.Getenv("REGISTRY_ADDRESS", "127.0.0.1:2379")),
		)),
		micro.WrapClient(hystrix.NewClientWrapper()),
	)
	service.Init()
	hystrixGo.DefaultMaxConcurrent = 5
	hystrixGo.DefaultTimeout = 300

	db, err := connect()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "connect to database"))
	}

	repo, err := repository.NewUserRepository(service.Client(), db)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "new repo"))
	}

	userHandler := handler.UserHandler{
		Repo: repo,
	}

	if err := pb.RegisterUserServiceHandler(service.Server(), userHandler); err != nil {
		log.Fatal(errors.WithMessage(err, "register handler"))
	}

	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run service"))
	}
}

func main() {
	log.SetFlags(log.Llongfile)
	go startAPIService()
	startMicroService()
}
