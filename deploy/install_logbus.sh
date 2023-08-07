#!/bin/bash

cd `dirname $0`

env=$1
if [[ $env = "" ]];then
        env="local"
fi

ps aux | grep logbus

if [ ! -f /home/ecs-user/webroot/logbus/logbus ]; then
  echo "Download Logbus"
  wget -O /home/ecs-user/admin/logbus2.tar.gz https://download.thinkingdata.cn/tools/release/ta-logBus-v2-linux-amd64.tar.gz
  tar -xzf /home/ecs-user/admin/logbus2.tar.gz -C /home/ecs-user/webroot/
fi

if diff "../config/logbus.${env}.json" "/home/ecs-user/webroot/logbus/conf/daemon.json" >/dev/null; then
  # 文件无差异，不需要更新
  echo "config file are same"
  cd /home/ecs-user/webroot/logbus
  /home/ecs-user/webroot/logbus/logbus progress
  exit 0
else
  echo "config file are different!!!"
  cp -rf "../config/logbus.${env}.json" "/home/ecs-user/webroot/logbus/conf/daemon.json"
  cd /home/ecs-user/webroot/logbus
  /home/ecs-user/webroot/logbus/logbus restart
  /home/ecs-user/webroot/logbus/logbus progress
fi

ps aux | grep logbus
