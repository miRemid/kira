ROOT_DIR=$PWD

echo "\033[32mCompile: \033[0m $1"
CGO_ENABLED=0 go build -a -ldflags '-s' -o ${ROOT_DIR}/deploy/bin/$1_service ${ROOT_DIR}/services/$1/main.go
echo "\033[32mFinish: \033[0m ${ROOT_DIR}/deploy/bin/$1_service"

cd ${ROOT_DIR}/deploy
echo "\033[32mPacket: \033[0m $1_service"
docker build -t kira/$1 -f ./build/$1/Dockerfile .
echo "\033[32mFinish: \033[0m kira/$1"
