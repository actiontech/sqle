#!/bin/bash

SQLE_BASE="/opt/sqle"
CONF="${SQLE_BASE}/etc/config.yaml"

if grep -q "# sqle-init" ${CONF}; then
    echo "sqle config has been init"
else
    # DMS 配置
    sed -i "s|host: 172.41.134.5|host:|" ${CONF}
    sed -i "s|port: 8978|port:|" ${CONF}
    sed -i "s|admin_user: administrator|admin_user:|" ${CONF}
    sed -i "s|admin_password: 123456|admin_password:|" ${CONF}
    sed -i "s|username: root|username: ${MYSQL_USER}|g" ${CONF}
    sed -i "s|password: 123|password: ${MYSQL_PASSWORD}|g" ${CONF}
    sed -i "s|host: 127.0.0.1|host: ${MYSQL_HOST}|g" ${CONF}
    sed -i "s|port: 3306|port: ${MYSQL_PORT}|g" ${CONF}
    sed -i "s|database: dms|database: ${MYSQL_DMS_SCHEMA}|g" ${CONF}
    sed -i "s|path: logs|path: ${SQLE_BASE}\/logs|g" ${CONF}
    # SQLE 配置
    sed -i "s|debug_log: false|debug_log: ${DEBUG}|g" ${CONF}
    sed -i "s|log_path: './logs'|log_path: '${SQLE_BASE}/logs'|g" ${CONF}
    sed -i "s|\.\/plugins|${SQLE_BASE}\/plugins|g" ${CONF}
    sed -i "s|true|${AUTO_MIGRATE_TABLE}|g" ${CONF}
    sed -i "s|mysql_host: '127.0.0.1'|mysql_host: '${MYSQL_HOST}'|g" ${CONF}
    sed -i "s|mysql_port: '3306'|mysql_port: ${MYSQL_PORT}|g" ${CONF}
    sed -i "s|mysql_user: 'root'|mysql_user: '${MYSQL_USER}'|g" ${CONF}
    sed -i "s|mysql_password: 'pass'|mysql_password: '${MYSQL_PASSWORD}'|g" ${CONF}
    sed -i "s|mysql_schema: 'sqle'|mysql_schema: '${MYSQL_SQLE_SCHEMA}'|g" ${CONF}
    
    echo -e "\n # sqle-init; DO NOT EDIT" >> ${CONF}
fi

echo "start sqle server..."
cd ${SQLE_BASE}
# 启动dms
./bin/dms -conf ${CONF}  >> std_dms.log 2>&1  &
sleep 3
# 检查dms进程是否存在
while ! pgrep -x "dms" > /dev/null; do
    sleep 1
done
# 启动sqle
./bin/sqled --config ${CONF} >> std.log 2>&1
