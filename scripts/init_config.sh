#!/bin/bash

# 主程序入口
# 获取用户输入的新路径作为work_directory（如果有）
work_directory="/opt/sqle"
split_line="-----------------------------------------------------------------"

if [ $# -gt 0 ]; then
    work_directory="$1"
fi

# 切换到 work_directory
cd "$work_directory" || { echo "Failed to change directory to $work_directory"; exit 1; }
echo $split_line
echo ">>>>> SQLE 为您提供 SQL 全生命周期质量管控能力 <<<<<"
echo $split_line
# 输出当前目录以确认切换成功
echo "[INFO] 当前工作目录: $(pwd)"

# 引用 functions.sh 文件
if ! source ./scripts/functions.sh; then
    echo $split_line
    echo "<<<< 请重新输入正确的工作目录 >>>>"
    echo $split_line
    exit 1
fi

echo "$info 为您自动备份配置文件"
backup_config_file
clear_config_sections
# 主循环
while true; do
    if ! choose_plugins;then
        continue
    fi
    
    if check_plugin_env; then
        if install_plugins; then
            echo ">>>> 本轮所有选择的插件均成功配置 <<<<"
            echo "$split_line"
            echo "累计成功配置的插件有："
            echo ""
            for ((i=0; i<${#successful_plugins[@]}; i++)); do
                edition=$(get_plugin_edition "${plugins[i]}")
                echo "$edition ${successful_plugins[i]}"
            done
            echo "$split_line"
            read -p "输入 'y'或'Y' 重启SQLE, 输入任意字符继续选择需要配置的插件>> " choice
            if [[ "$choice" == "y" || "$choice" == "Y" ]]; then
                echo "$split_line"
                # 将不存在于successful_plugins中的插件移动到unused_plugin
                echo "$info 本次启动未使用的插件文件将移动到unused_plugins中"
                for plugin in "${plugins[@]}"; do
                    if [[ ! " ${successful_plugins[@]} " =~ " ${plugin} " ]]; then
                        if ! move_to_unused_plugins "$plugin"; then
                            echo "$error 移动插件 $plugin 到 unused_plugins 目录失败"
                        fi
                    fi
                done
                restart_sqle
                error_code=$?
                if [ $error_code -ne 0 ];then
                    if [ $error_code -eq 1 ]; then
                      echo "$error 重启SQLE失败，原因是：设置权限失败，脚本已退出"
                    elif [ $error_code -eq 2 ]; then
                      echo "$error 重启SQLE失败，原因是：停止脚本不存在，脚本已退出"
                    else
                      echo "$error 重启SQLE失败, 脚本已退出"
                    fi
                fi
                break
            else
                continue
            fi
        else
            echo "累计成功配置的插件有："
            for ((i=0; i<${#successful_plugins[@]}; i++)); do
                edition=$(get_plugin_edition "${plugins[i]}")
                echo "$edition ${successful_plugins[i]}"
            done
            echo "以下插件配置失败: "
            for ((i=0; i<${#failed_plugins[@]}; i++)); do
                edition=$(get_plugin_edition "${plugins[i]}")
                echo "$edition $((i+1)). ${failed_plugins[i]}"
            done
            echo "$split_line"
            echo "$info 插件配置请参考SQLE插件配置文档:"
            echo "$info https://actiontech.github.io/sqle-docs/docs/dev-manual/plugins/howtouse"
            echo "$info 请确认失败原因，解决失败问题，然后对上述失败的插件重新配置"
            break
        fi
    else
        failed_message
        break
    fi
done