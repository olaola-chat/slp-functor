#!/bin/bash

cd "$(dirname "$0")" || exit
cd ..

git pull origin dev
make build

path="/root/go/src/github.com/olaola-chat/rbp-functor"
targetPath="/home/webroot/rbp-functor"
logPath="/home/log"

if [ ! -d "$targetPath" ]; then
	mkdir -p "$targetPath"
fi

if [ ! -d "$logPath" ]; then
	mkdir -p "$logPath"
fi

#复制目录过去
dirs=("bin" "config")
for i in ${!dirs[@]}; do
	cp -rf "${path}/${dirs[i]}" "${targetPath}"
	if [ $? -ne 0 ]; then
		echo "error to copy agent to target";
		exit 1
	fi
done

#copy supervisor 配置文件
#改成自动获取目录里的文件来
#重启supervisor守护进程
files=("http" "rpc.voice_lover_admin" "rpc.voice_lover_main")
for i in ${!files[@]}; do
	superFile="/etc/supervisor/conf.d/rbp-functor.${files[i]}.conf"
	localFile="${path}/deploy/dev/rbp-functor.${files[i]}.conf"
	cp -f "${localFile}" "${superFile}"
	supervisorctl restart "rbp-functor.${files[i]}"
done

echo "ok"
exit 0


