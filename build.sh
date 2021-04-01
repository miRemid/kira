ROOT_DIR=$PWD
services=`ls $ROOT_DIR/services`
for service in $services
do
    echo "\033[32mCompile: \033[0m $service"
    cd ${ROOT_DIR}/services/$service
    make static docker
    echo "\033[32mFinish: \033[0m $service"
done