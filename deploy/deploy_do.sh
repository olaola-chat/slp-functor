#!/bin/bash  

cd `dirname $0`
cd ..

if [ $# -lt 1 ];
then
        echo "USAGE: $0 [-force] serviceName"
        exit 1
fi

serviceName=$1
force=false

if [ "$1" = "-force" ]; then
    serviceName=$2
    force=true
fi

path=$(pwd)

targetPath="/home/ecs-user/webroot/${serviceName}"
logPath="/home/ecs-user/log/${serviceName}"
supervisorPath="/home/ecs-user/.local/etc/supervisor/conf.d"

superFile="${supervisorPath}/${serviceName}.conf"
localFile="${path}/deploy/prod/${serviceName}.conf"

if [ ! -f "${localFile}" ];then
    echo "serviceName is not valid";
    exit 1;
fi

if [ ! -d "$targetPath" ]; then
	mkdir -p "$targetPath"
fi

if [ ! -d "$logPath" ]; then
	mkdir -p "$logPath"
fi

if [ ! -d "$supervisorPath" ]; then
	mkdir -p "$supervisorPath"
fi


copy_dirs() {
    #复制目录过去
    local dirs=("template" "public" "bin" "config" "i18n")
    for i in ${!dirs[@]}; do
        cp -rf "${path}/${dirs[i]}" "${targetPath}"
        if [ $? -ne 0 ]; then
            echo "error to copy agent to target";
            exit 1
        fi
    done
}

if [ ! -f "${superFile}" ];then
    cp -f "${localFile}" "${superFile}"

    copy_dirs

    #更新 supervisor 配置
    #系统会自动启动进程
    supervisorctl update "${serviceName}"

else
    #检查文件是否发生变化了
    diff "${localFile}" "${superFile}" > /dev/null
    if [ $? -eq 0 ]; then
        echo "${serviceName} file are same"

        if [ "$force" = "false" ]; then
            #文件没变, 什么都不做
            exit 0;
        fi

        copy_dirs

        supervisorctl restart "$serviceName"
    else
        echo "${serviceName} file are different"

        cp -f "${localFile}" "${superFile}"

        copy_dirs

        supervisorctl update "$serviceName"
    fi
	
	if [ $? -ne 0 ]; then
		echo "error with supervisorctl $serviceName";
		exit 1
	fi
fi
	
# 判断进程状态
for k in {1..5}
do
    sleep 1
    v=`supervisorctl status "$serviceName" | grep "RUNNING" | wc -l`
    if [ $v -eq "0" ]; then
    echo "error status with $serviceName";
        exit 1
    else
        echo "check status ${k} ok"
    fi
done

echo "ok"
exit 0

