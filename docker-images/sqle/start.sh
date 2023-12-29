#!/bin/bash

SQLE_BASE="/opt/sqle"
CONF="${SQLE_BASE}/etc/config.yml"


if [ ! -f "${CONF}" ];then
  cat > ${CONF} <<EOF
dms:
  id: 1
  api:
    addr: 0.0.0.0
    port: 10000
  service:
    database:
      username: ${MYSQL_USER}
      password: ${MYSQL_PASSWORD}
      host: ${MYSQL_HOST}
      port: ${MYSQL_PORT}
      database: ${MYSQL_DMS_SCHEMA}
      debug: ${DEBUG}
      auto_migrate: true
    log:
      level: DEBUG
      path: logs
      max_size_mb: 100
      max_backup_number: 10      
  cloudbeaver:
    enable_https: 
    host: 
    port: 
    admin_user: 
    admin_password: 
  database_driver_options:
    - db_type: MySQL
      logo_path: "/logo/mysql.png"
      params: 
  secret_key:  
sqle:
  id: 1
  dms_server_address: http://127.0.0.1:10000
  api:
    addr: 127.0.0.1
    port: 10001
    enable_https: false
    cert_file_path: './etc/cert.pem'
    key_file_path: './etc/key.pem'
  secret_key:
  service:
    auto_migrate_table: ${AUTO_MIGRATE_TABLE}
    debug_log: ${DEBUG}
    log_path: '${SQLE_BASE}/logs'
    log_max_size_mb: 1024
    log_max_backup_number: 2
    plugin_path: '${SQLE_BASE}/plugins'
    database:
      mysql_host: '${MYSQL_HOST}'
      mysql_port:  ${MYSQL_PORT}
      mysql_user: '${MYSQL_USER}'
      mysql_password: '${MYSQL_PASSWORD}'
      mysql_schema: '${MYSQL_SQLE_SCHEMA}'
EOF
fi

echo "start sqle server..."
cd ${SQLE_BASE}
# 启动dms
./bin/dms -conf ${CONF} &
sleep 3
# 检查dms进程是否存在
while ! pgrep -x "dms" > /dev/null; do
    sleep 1
done
# 启动sqle
./bin/sqled --config ${CONF}
