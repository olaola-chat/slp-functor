#!/bin/bash  

cd `dirname $0`


force=""

if [ "$1" = "-force" ]; then
    force="-force"
fi

bash deploy_do.sh ${force} slp-functor.rpc.voice_lover_admin
bash deploy_do.sh ${force} slp-functor.rpc.voice_lover_main
