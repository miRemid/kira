module github.com/miRemid/kira

go 1.14

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/casbin/casbin/v2 v2.19.8
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/golang/protobuf v1.4.3
	github.com/gomodule/redigo v1.8.3
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/micro/v2 v2.9.3
	github.com/minio/minio-go/v7 v7.0.7
	github.com/pkg/errors v0.9.1
	github.com/rs/xid v1.2.1
	github.com/segmentio/ksuid v1.0.3
	github.com/teris-io/shortid v0.0.0-20201117134242-e59966efd125
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	google.golang.org/protobuf v1.25.0
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.9
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
