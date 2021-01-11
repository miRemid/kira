ROOT_DIR=/home/mio/Workspace/kira

services="
auth
gateway
user
file
"

build_service() {
    CGO_ENABLED=0 go build -a -ldflags '-s' -o ${ROOT_DIR}/deploy/bin/$1_service ${ROOT_DIR}/services/$1/main.go
    echo -e "\033[32m编译完成: \033[0m ${ROOT_DIR}/deploy/bin/$1"
}

build_image() {
    docker build -t kira/$1 -f ./services/$1/Dockerfile .
    echo -e "\033[32m镜像打包完成: \033[0m kira/$1\n"
}

cd ${ROOT_DIR}

# mkdir -p ${ROOT_DIR}/deploy/bin && rm -f ${ROOT_DIR}/deploy/bin/*

# for service in $services
# do
#     build_service $service
# done

echo -e "\033[32mFinish building services\033[0m"

cd ${ROOT_DIR}/deploy/
for service in $services
do
    build_image $service
done

echo -e "\033[32mFinish building docker images\033[0m"