#!/bin/bash

cd `dirname $0`

force=""

if [ "$1" = "-force" ]; then
    force="-force"
fi


names=(cache.profile cache.aggregation cache.room cache.follow)
names+=(exhibition treat.task wakeup.push autochat cloud.view)
names+=(gift.reward gift.wall login)

# ./consul-game-config.sh "172.16.0.179:8500"

for name in "${names[@]}"; do
  ./deploy_do.sh ${force} "rbp.cmd.${name}"
done


bash deploy_do.sh ${force} rbp.cmd.room.operation
bash deploy_do.sh ${force} rbp.cmd.room.vip
bash deploy_do.sh ${force} rbp.cmd.room.tag
bash deploy_do.sh ${force} rbp.cmd.special.attention
bash deploy_do.sh ${force} rbp.cmd.facedet
bash deploy_do.sh ${force} rbp.cmd.hour.rank
bash deploy_do.sh ${force} rbp.cmd.defend.rank
bash deploy_do.sh ${force} rbp.cmd.rank.match
bash deploy_do.sh ${force} rbp.cmd.argus
bash deploy_do.sh ${force} rbp.cmd.familiar
bash deploy_do.sh ${force} rbp.cmd.match
bash deploy_do.sh ${force} rbp.cmd.viability
bash deploy_do.sh ${force} rbp.cmd.room.purrecommend
bash deploy_do.sh ${force} rbp.cmd.sms
bash deploy_do.sh ${force} rbp.cmd.room.data
bash deploy_do.sh ${force} rbp.cmd.roomdata.sum
bash deploy_do.sh ${force} rbp.cmd.achieve
bash deploy_do.sh ${force} rbp.cmd.anchor
bash deploy_do.sh ${force} rbp.cmd.cphouse
bash deploy_do.sh ${force} rbp.cmd.smartgreet
bash deploy_do.sh ${force} rbp.cmd.confess
bash deploy_do.sh ${force} rbp.cmd.roomlist
bash deploy_do.sh ${force} rbp.cmd.cron
bash deploy_do.sh ${force} rbp.cmd.monitor

./install_logbus.sh prod

echo "ok"
exit 0
