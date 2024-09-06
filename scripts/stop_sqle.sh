#!/bin/bash

systemctl stop dms.service
systemctl stop sqled.service
echo "<<<< SQLE已关闭 >>>>"