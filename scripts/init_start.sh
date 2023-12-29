#!/bin/bash

# 本脚本用于在rpm安装完成后，启动dms服务和sqled服务

detectSystemServiceStatus() {
    detectedRunningTimes=0
    for i in {1..30}; do
        sleep 1
    	  if systemctl status "$1" &>/dev/null; then
            ((detectedRunningTimes++))
            if [[ $detectedRunningTimes == 3 ]]; then
                return 0
            fi
        else
            detectedRunningTimes=0
        fi
    done

    return 1
}

systemctl daemon-reload
systemctl start dms.service

# 10秒dms还没有启动成功退出
if detectSystemServiceStatus "dms.service"; then
    echo "init and start dms success!"
else
    echo "duration 30 seconds; start dms failed"
    exit 3
fi

systemctl start sqled.service

# 10秒sqle还没有启动成功退出
if detectSystemServiceStatus "sqled.service"; then
    echo "init and start sqled success!"
else
    echo "duration 30 seconds; start sqled failed"
fi