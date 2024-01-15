#!/bin/bash

# 本脚本用于在rpm安装完成后，启动dms服务和sqled服务

detectSystemServiceStatus() {
    detectedRunningTimes=0
    for i in {1..30}; do
        sleep 1
    	  if systemctl status "$1" &>/dev/null; then
            ((detectedRunningTimes++))
            if [[ $detectedRunningTimes == 3 ]]; then
                echo "init and start $1 success!"
                return 0
            fi
        else
            detectedRunningTimes=0
        fi
    done

    echo "duration 30 seconds; start $1 failed"
    return 1
}

systemctl daemon-reload

systemctl start dms.service
if ! detectSystemServiceStatus "dms.service"; then
    exit 3
fi

systemctl start sqled.service
if ! detectSystemServiceStatus "sqled.service"; then
    exit 3
fi

systemctl start neo4j.service
if ! detectSystemServiceStatus "neo4j.service"; then
    exit 3
fi

systemctl start provision.service
if ! detectSystemServiceStatus "provision.service"; then
    exit 3
fi