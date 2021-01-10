module github.com/miRemid/kira

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.4.3
	github.com/micro/go-micro/v2 v2.9.1
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
