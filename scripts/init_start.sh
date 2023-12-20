#!/bin/bash

# 本脚本用于在rpm安装完成后，启动dms服务和sqled服务

systemctl daemon-reload
systemctl start dms.service

for i in {1..10}; do
	systemctl status dms.service &>/dev/null
    if [  $? -eq 0 ]; then
        echo "Init and start dms success!"
        break
    fi
    sleep 1
done

systemctl start sqled.service
for i in {1..10}; do
	systemctl status sqled.service &>/dev/null
    if [  $? -eq 0 ]; then
        echo "Init and start sqled success!"
        exit 0
    fi
    sleep 1
done

echo "init and start sqle fail! Please check sqle log"