build(){
    path=`pwd`
    cd "services"
    for dir in `ls`; do
        cd $dir

        make build

        mv ./${dir}_service ${path}/deploy/bin


        cd ..
    done
    cd ..
}

build