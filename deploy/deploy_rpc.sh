#!/bin/bash  

cd `dirname $0`


force=""

if [ "$1" = "-force" ]; then
    force="-force"
fi

bash deploy_do.sh ${force} rbp.rpc.user
bash deploy_do.sh ${force} rbp.rpc.room.info
bash deploy_do.sh ${force} rbp.rpc.pay
bash deploy_do.sh ${force} rbp.rpc.cache
bash deploy_do.sh ${force} rbp.rpc.rce
bash deploy_do.sh ${force} rbp.rpc.im
bash deploy_do.sh ${force} rbp.rpc.ame
bash deploy_do.sh ${force} rbp.rpc.location
bash deploy_do.sh ${force} rbp.rpc.room.tag
bash deploy_do.sh ${force} rbp.rpc.wlc
bash deploy_do.sh ${force} rbp.rpc.match
bash deploy_do.sh ${force} rbp.rpc.screenmessage
bash deploy_do.sh ${force} rbp.rpc.recommend
bash deploy_do.sh ${force} rbp.rpc.user.gift
bash deploy_do.sh ${force} rbp.rpc.user.relation
bash deploy_do.sh ${force} rbp.rpc.risk
bash deploy_do.sh ${force} rbp.rpc.consume
bash deploy_do.sh ${force} rbp.rpc.user.viability
bash deploy_do.sh ${force} rbp.rpc.store
bash deploy_do.sh ${force} rbp.rpc.music.song_source
bash deploy_do.sh ${force} rbp.rpc.user.anonymous
bash deploy_do.sh ${force} rbp.rpc.admin.verify
