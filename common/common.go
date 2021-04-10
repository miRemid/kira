package common

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	AnonymousKey = "ANONYMOUS_FILES_ID"
	AnonyEvent   = "kira.micro.service.upload.anony"
	AnonyBucket  = "anony"
)

func Getenv(key, replace string) string {
	res := os.Getenv(key)
	if res == "" {
		return replace
	}
	return res
}

func getConnect() string {
	username := Getenv("MYSQL_USERNAME", "shi")
	password := Getenv("MYSQL_PASSWORD", "123456")
	database := Getenv("MYSQL_DATABASE", "kira")
	address := Getenv("MYSQL_ADDRESS", "127.0.0.1:3306")
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, address, database)
}

func DBConnect() (*gorm.DB, error) {
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
	openConns, err := strconv.Atoi(Getenv("MAX_OPEN_CONNS", "0"))
	if err != nil {
		return nil, errors.WithMessage(err, "MAX_OPEN_CONNS must be a integer")
	}
	idelConns, err := strconv.Atoi(Getenv("MAX_IDLE_CONNS", "0"))
	if err != nil {
		return nil, errors.WithMessage(err, "MAX_IDLE_CONNS must be a integer")
	}
	db.SetMaxOpenConns(openConns)
	db.SetMaxIdleConns(idelConns)
	db.SetConnMaxLifetime(time.Minute * 5)
	return _conn, nil
}

func MinioConnect() (*minio.Client, error) {
	var ssl = false
	endpoint := Getenv("MINIO_ENDPOINT", "127.0.0.1:9900")
	accessKey := Getenv("MINIO_ACCESSKEY", "kira")
	secretKey := Getenv("MINIO_SECRETKEY", "1234567890")
	secure := Getenv("MINIO_SECURE", "false")
	if secure == "true" {
		ssl = true
	}
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: ssl,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}
