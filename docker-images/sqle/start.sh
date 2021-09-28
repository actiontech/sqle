#!/bin/bash

SQLE_BASE="/opt/sqle"
SQLE_CONF="${SQLE_BASE}/etc/sqled.yml"

if [ ! -f "${SQLE_CONF}" ];then
  cat > ${SQLE_CONF} <<EOF
  server:
    sqle_config:
      server_port: 10000
      auto_migrate_table: ${AUTO_MIGRATE_TABLE}
      debug_log: ${DEBUG}
      log_path: '${SQLE_BASE}/logs'
      plugin_path: '${SQLE_BASE}/plugins'
    db_config:
      mysql_cnf:
        mysql_host: '${MYSQL_HOST}'
        mysql_port: ${MYSQL_PORT}
        mysql_user: '${MYSQL_USER}'
        mysql_password: '${MYSQL_PASSWORD}'
        mysql_schema: '${MYSQL_SCHEMA}'
EOF
fi

echo "start sqle server..."
cd ${SQLE_BASE}
./bin/sqled --config ${SQLE_CONF}