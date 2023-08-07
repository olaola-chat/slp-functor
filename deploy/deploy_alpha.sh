#!/bin/bash  

httpExitFun(){
	curl "http://127.0.0.1:9090/unregister"
}

path="/home/admin/rbp-go"
targetPath="/home/webroot/rbp-go"
logPath="/home/log"

if [ ! -d "$targetPath" ]; then
	mkdir -p "$targetPath"
fi

if [ ! -d "$targetPath/config" ]; then
	mkdir -p "$targetPath/config"
fi

if [ ! -d "$logPath" ]; then
	mkdir -p "$logPath"
fi

#复制目录过去
dirs=("template" "public" "bin" "config" "deploy" "i18n")
for i in ${!dirs[@]}; do
	cp -rf "${path}/${dirs[i]}" "${targetPath}"
	if [ $? -ne 0 ]; then
		echo "error to copy agent to target";
		exit 1
	fi
done

#copy supervisor 配置文件
#改成自动获取目录里的文件来
#只有当配置文件发生改变时，才会重启进程
#就是说发布时，想要重启哪些进程，就更改下对应文件的APP_VERSION
#先这样部署，稍后更改，不能一台机器部署所有rpc服务
#不擅长写sh, 直接写个php or golang 执行？
files=("http" "rpc.user" "rpc.room.info" "rpc.pay" "rpc.cache" "rpc.rce" "rpc.im" "rpc.screenmessage" "rpc.recommend" "rpc.user.gift" "rpc.user.relation" "rpc.risk" "rpc.store" "rpc.music.song_source")
for i in ${!files[@]}; do
	superFile="/etc/supervisor/conf.d/rbp.${files[i]}.conf"
	localFile="${path}/deploy/${files[i]}.conf"
	if [ ! -f "${superFile}" ];then
		cp -f "${localFile}" "${superFile}"
		
		if [ "${files[i]}" == "http" ]; then
			#暂停nginx服务
			httpExitFun
		fi
		
		#更新 supervisor 配置
		#系统会自动启动进程
		supervisorctl update "rbp.${files[i]}"
	else
		#检查文件是否发生变化了
		diff "${localFile}" "${superFile}" > /dev/null
		if [ $? -eq 0 ]; then
			echo "${files[i]} file are same"
			#重启进程
			#supervisorctl restart "${files[i]}"
			#文件没变，不重启
			continue
		else
		    echo "${files[i]} file are different"
			cp -f "${localFile}" "${superFile}"
			
			if [ "${files[i]}" == "http" ]; then
				#暂停nginx服务
				httpExitFun
			fi
		
			supervisorctl update "rbp.${files[i]}"
		fi
	fi
	
	if [ $? -ne 0 ]; then
		echo "error with supervisorctl ${files[i]}";
		exit 1
	fi
	
	# 判断进程状态
	for k in {1..5}
	do
	    sleep 1
	    v=`supervisorctl status "rbp.${files[i]}" | grep "RUNNING" | wc -l`
	    if [ $v -eq "0" ]; then
	        echo "error status with rbp.${files[i]}";
	        exit 1
	    else
	        echo "check status ${k} ok"
	    fi
	done
done

echo "ok"
exit 0


