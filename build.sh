ROOT_DIR=$PWD

services="
auth
user
file
upload
site
"

build_service() {
    echo -e "\033[32mCompile: \033[0m $1_service"
    CGO_ENABLED=0 go build -a -ldflags '-s' -o ${ROOT_DIR}/deploy/bin/$1_service ${ROOT_DIR}/services/$1/main.go
    echo -e "\033[32mFinish: \033[0m ${ROOT_DIR}/deploy/bin/$1_service"
}

build_image() {
    echo -e "\033[32mPacket: \033[0m $1_service"
    docker build -t kira/$1 -f ./build/$1/Dockerfile .
    echo -e "\033[32mFinish: \033[0m kira/$1\n"
}

cd ${ROOT_DIR}

mkdir -p ${ROOT_DIR}/deploy/bin && rm -f ${ROOT_DIR}/deploy/bin/*_service

for service in $services
do
    build_service $service
done

echo -e "\033[32mFinish building services\033[0m"

cd ${ROOT_DIR}/deploy/
for service in $services
do
    build_image $service
done

echo -e "\033[32mFinish building docker images\033[0m"