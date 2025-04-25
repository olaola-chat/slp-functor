FROM slp-acr-registry.cn-hangzhou.cr.aliyuncs.com/slp/slp-ubuntu:22.04

WORKDIR /home/ecs-user/webroot/slp-functor

COPY bin/rpc /home/ecs-user/webroot/slp-functor/bin/slp-functor-rpc
COPY config /home/ecs-user/webroot/slp-functor/config
COPY i18n /home/ecs-user/webroot/slp-functor/i18n
COPY public /home/ecs-user/webroot/slp-functor/public
COPY template /home/ecs-user/webroot/slp-functor/template

