#!/bin/bash
professional="【专业版】"
enterprise="【企业版】"
# 定义可安装的插件
plugins=(
    "sqle-tidb-plugin"
    "sqle-td-mysql-plugin"
    "sqle-pg-plugin"
    "sqle-ob-mysql-plugin"
    "sqle-mycat-plugin"
    "sqle-ms-plugin"
    "sqle-db2-plugin.jar"
    "sqle-oracle-plugin.jar"
    "sqle-tbase-plugin"
    "sqle-goldendb-plugin"
    "sqle-dm-plugin.jar"
    "sqle-ob-for-oracle-plugin.jar"
)
# 若需要添加一个插件的配置，需要在plugins、get_plugin_edition和install_plugin三处进行配置
# 若有对环境依赖，需要在check_env中增加环境依赖检测项


# 全局变量
info="[INFO]"
error="[ERROR]"
warn="[WARN]"
config_file="./etc/config.yaml"
plugins_dir="./plugins"
successful_plugins=()
failed_plugins=()

choose_plugins() {
    echo "$split_line"
    echo ">>>> 步骤1: 请选择要安装的插件 <<<<"
    echo "$split_line"
    for ((i=0; i<${#plugins[@]}; i++)); do
        edition=$(get_plugin_edition "${plugins[i]}")
        echo "$edition $((i+1)). ${plugins[i]} "
    done
    
    echo "$split_line"
    echo "$info 请输入插件序号, 以空格分隔, 如: 1 2 3"
    echo "$info 若要安装所有插件，请输入'all'或'ALL', 如: all"
    read -p "$info 请输入>> " input
    echo "$split_line"
    if [[ -z "$input" ]]; then
        echo "<<<<请输入正确的插件序号>>>>"
        return 1
    fi
    
    # 检查是否输入了 "all"
    if [[ "$input" == "all"||"$input" == "ALL" ]]; then
        selected_plugins=("${plugins[@]}")
        return 0
    fi
    
    # 将用户输入的序号转换为插件列表
    selected_plugins=()
    for index in $input; do
        if [[ ! $index =~ ^[0-9]+$ ]] || [[ $index -lt 1 ]] || [[ $index -gt ${#plugins[@]} ]]; then
            echo "输入的序号 $index 无效, 请重新输入有效的序号"
            return 1
        fi
        # 检查序号是否有效，并将对应插件添加到selected_plugins数组中
        index=$((index-1)) # 转换为数组索引（从0开始）
        if [[ $index -ge 0 ]] && [[ $index -lt ${#plugins[@]} ]]; then
            selected_plugins+=("${plugins[$index]}")
        else
            echo "序号 $index 超出插件列表范围，请重新运行脚本并输入有效的序号"
            return 1
        fi
    done
}

get_plugin_edition() {
    local plugin_name="$1"
    case "$plugin_name" in
        sqle-dm-plugin.jar | sqle-goldendb-plugin | sqle-ob-for-oracle-plugin.jar | sqle-tbase-plugin)
            echo "【企业版】"
        ;;
        *)
            echo "【企业版】【专业版】"
        ;;
    esac
}

check_plugin_env() {
    echo ">>>> 步骤2: 预检查插件的运行条件 <<<<"
    all_checks_passed=true
    for plugin in "${selected_plugins[@]}"; do
        if ! check_env "$plugin"; then
            all_checks_passed=false
        fi
    done
    echo "$split_line"
    if $all_checks_passed; then
        echo "$info 插件运行条件预检查通过"
        return 0
    fi
    echo "$error 插件运行条件预检查不通过"
    echo "$info 插件配置请参考SQLE插件配置文档:"
    echo "$info https://actiontech.github.io/sqle-docs/docs/dev-manual/plugins/howtouse"
    echo "$info 请确认失败原因，解决失败问题，然后对上述失败的插件重新配置"
    echo "若需要社区支持，欢迎加入SQLE社区交流群："
    echo "https://actiontech.github.io/sqle-docs/docs/support/community-support"
    return 1
}

install_plugins() {
    echo "$split_line"
    echo ">>>> 步骤3: 安装配置插件 <<<<"
    echo "$split_line"
    
    all_install_success=true
    add_plugin_config_to_sqle_service
    add_plugin_config_to_dms
    ## 默认添加MySQL配置
    if ! install_plugin_to_config "MySQL" "/logo/mysql.svg" "" ;then
        return 1
    fi
    for plugin in "${selected_plugins[@]}"; do
        if ! install_plugin "$plugin"; then
            echo "$error 插件 $plugin 配置失败"
            echo "$split_line"
            all_install_success=false
            # 如果插件不在 failed_plugins 中才添加
            if [[ ! " ${failed_plugins[@]} " =~ " $plugin " ]]; then
                failed_plugins+=("$plugin")
            fi
        else
            echo "$info 插件 $plugin 配置完成"
            echo "$split_line"
            # 如果插件不在 successful_plugins 中才添加
            if [[ ! " ${successful_plugins[@]} " =~ " $plugin " ]]; then
                successful_plugins+=("$plugin")
            fi
        fi
    done
    
    # 返回0表示所有插件成功安装，返回1表示有插件安装失败
    if $all_install_success; then
        return 0
    else
        return 1
    fi
}

failed_message(){
    echo ""
    echo "X_X"
}

check_env() {
    plugin=$1
    echo "$split_line"
    echo "$info 正在检查 $plugin 的运行条件"
    # 检查插件文件是否存在
    plugin_file="$work_directory/plugins/$plugin"
    unused_plugin_file="$work_directory/unused_plugins/$plugin"
    if ! check_file_existence $plugin_file; then
        if ! check_file_existence $unused_plugin_file; then
            echo "$error 查找插件可执行文件: $plugin_file, 文件不存在"
            return 1
        else
            if ! move_to_plugins $plugin;then
                echo "$error 移动插件 $plugin_file 到plugins文件夹失败"
            fi
        fi
    fi
    # 检查插件依赖环境是否已配置
    case "$plugin" in
        sqle-ms-plugin)
            # 检查 .NET 6 是否安装
            if check_dotnet6; then
                return 0
            else
                return 1
            fi
        ;;
        sqle-db2-plugin.jar | sqle-oracle-plugin.jar|sqle-dm-plugin.jar|sqle-ob-for-oracle-plugin.jar) # todo 确认达梦和oboracle是否需要java1.8还是其他版本也可以
            if check_java8; then
                return 0
            else
                return 1
            fi
        ;;
    esac
    return 0
}

# 定义函数，检查文件是否存在
check_file_existence() {
    local file_path="$1"
    
    if [ -f "$file_path" ]; then
        return 0  # 文件存在
    else
        return 1  # 文件不存在
    fi
}

check_dotnet6() {
    command -v dotnet >/dev/null 2>&1
    if [ $? -eq 0 ]; then
        version=$(dotnet --version)
        if [[ $version == 6.* ]]; then
            echo "$info .NET6 已安装 (版本 $version)"
            return 0
        else
            echo "$warn 已安装的 .NET 版本为 $version, 但不是 .NET6"
            return 1
        fi
    else
        echo "$error .NET6未安装"
        return 1
    fi
}

check_java8() {
    command -v java >/dev/null 2>&1
    if [ $? -eq 0 ]; then
        version=$(java -version 2>&1 | awk -F '"' '/version/ {print $2}')
        if [[ $version == 1.8.* ]]; then
            echo "$info Java JDK 1.8 已安装 (版本 $version)"
            return 0
        else
            echo "$warn 已安装的 Java 版本为 $version, 但不是 JDK 1.8"
            return 2
        fi
    else
        echo "$error Java JDK 1.8 未安装"
        return 1
    fi
}

backup_config_file() {
    local date_str=$(date +"%Y-%m-%d_%H-%M-%S")
    local backup_dir="./etc/backup"
    local backup_file="$backup_dir/config.yaml.back.$date_str"
    # 检查 backup 文件夹是否存在，如果不存在则创建
    if [ ! -d "$backup_dir" ]; then
        mkdir -p "$backup_dir"
        if [ $? -ne 0 ]; then
            echo "$error 无法创建备份文件夹 $backup_dir"
            return 1
        fi
    fi
    
    # 检查配置文件是否存在
    if [ -f "$config_file" ]; then
        # 复制配置文件到备份文件
        if cp "$config_file" "$backup_file"; then
            echo "$info 配置文件已备份为 $backup_file"
            return 0
        else
            echo "$error 备份配置文件失败"
            return 1
        fi
    else
        echo "$error 配置文件 $config_file 不存在"
        return 1
    fi
}

# 函数：安装插件
install_plugin() {
    plugin=$1
    if ! add_execute_permission $plugin; then
        echo "$error 为插件增加可执行权限失败"
        return 1
    fi
    case "$plugin" in
        sqle-ms-plugin)
            if ! install_plugin_to_config "SQL Server" "/logo/sql_server.svg" "" ;then
                return 1
            fi
        ;;
        sqle-td-mysql-plugin)
            if ! install_plugin_to_config "TDSQL For InnoDB" "/logo/tdsql_for_innodb.jpeg" "" ;then
                return 1
            fi
        ;;
        sqle-tidb-plugin)
            if ! install_plugin_to_config "TiDB" "/logo/ti_db.png" "" ;then
                return 1
            fi
        ;;
        sqle-pg-plugin)
            if ! install_plugin_to_config "PostgreSQL" "/logo/postgres_sql.svg" "" ;then
                return 1
            fi
        ;;
        sqle-oracle-plugin.jar)
            if ! install_plugin_to_config "Oracle" "/logo/oracle.svg" "[\n        {\n          \"key\":\"service_name\",\n          \"value\":\"XE\",\n          \"desc\":\"service name\",\n          \"type\":\"string\"\n        }\n      ]" ;then
                return 1
            fi
            if ! install_plugin_to_plugin_config "sqle-oracle-plugin.jar" "java -jar $plugins_dir/sqle-oracle-plugin.jar" ;then
                return 1
            fi
        ;;
        sqle-db2-plugin.jar)
            if ! install_plugin_to_config "DB2" "/logo/db2.jpeg" "[\n        {\n          \"key\":\"database_name\",\n          \"value\":\"\",\n          \"desc\":\"数据库名\",\n          \"type\":\"string\"\n        }\n      ]" ;then
                return 1
            fi
            if ! install_plugin_to_plugin_config "sqle-db2-plugin.jar" "java -jar $plugins_dir/sqle-db2-plugin.jar" ;then
                return 1
            fi
        ;;
        sqle-mycat-plugin)
            if ! install_plugin_to_config "Mycat" "/logo/mycat.png" "" ;then
                return 1
            fi
        ;;
        sqle-ob-mysql-plugin)
            if ! install_plugin_to_config "OceanBase For MySQL" "/logo/ob_for_mysql.png" "[\n        {\n          \"key\":\"tenant_name\",\n          \"value\":\"\",\n          \"desc\":\"tenant\",\n          \"type\":\"string\"\n        }\n      ]" ;then
                return 1
            fi
        ;;
        sqle-tbase-plugin)
            if ! install_plugin_to_config "TBase" "/logo/tbase.png" "" ;then
                return 1
            fi
        ;;
        sqle-goldendb-plugin)
            if ! install_plugin_to_config "GoldenDB" "/logo/golden_db.png" "" ;then
                return 1
            fi
        ;;
        sqle-dm-plugin.jar)
            if ! install_plugin_to_config "DM" "/logo/dameng_db.png" "" ;then
                return 1
            fi
            if ! install_plugin_to_plugin_config "sqle-dm-plugin.jar" "java -jar $plugins_dir/sqle-dm-plugin.jar" ;then
                return 1
            fi
        ;;
        sqle-ob-for-oracle-plugin.jar)
            if ! install_plugin_to_config "OceanBase For Oracle" "/logo/ob_for_oracle.png" "[\n        {\n          \"key\":\"service_name\",\n          \"value\":\"SYS\",\n          \"desc\":\"service name\",\n          \"type\":\"string\"\n        }\n      ]" ;then
                return 1
            fi
            if ! install_plugin_to_plugin_config "sqle-ob-for-oracle-plugin.jar" "java -jar $plugins_dir/sqle-ob-for-oracle-plugin.jar" ;then
                return 1
            fi
        ;;
        # 其他插件的安装可以在这里添加
        *)
            # 默认的安装流程
            echo "$error 暂不支持插件 $plugin 的自动配置"
            return 1
        ;;
    esac
    return 0
}

# 函数：为 ../plugins 文件夹下的指定文件增加执行权限
add_execute_permission() {
    local file_name="$1"
    
    # 检查 ../plugins 目录是否存在
    if [ -d "$plugins_dir" ]; then
        # 检查文件是否存在
        if [ -f "$plugins_dir/$file_name" ]; then
            # 尝试给文件增加执行权限
            if chmod +x "$plugins_dir/$file_name"; then
                echo "$info 插件 $file_name 可执行权限配置完成"
                return 0
            else
                echo "$error 设置 $file_name 执行权限失败"
                return 1
            fi
        else
            echo "$error 找不到 $plugins_dir/$file_name 文件"
            return 1
        fi
    else
        echo "$error 找不到 $plugins_dir 目录"
        return 1
    fi
}

move_to_unused_plugins() {
    local file_name="$1"
    local unused_plugin_dir="./unused_plugins"
    
    # 检查 unused_plugin 目录是否存在，不存在则创建
    if [ ! -d "$unused_plugin_dir" ]; then
        mkdir "$unused_plugin_dir"
        echo "$info 目录 $unused_plugin_dir 已创建"
    fi
    
    # 检查 ../plugins 目录是否存在
    if [ -d "$plugins_dir" ]; then
        # 检查文件是否存在
        if [ -f "$plugins_dir/$file_name" ]; then
            # 尝试移动文件到 unused_plugin 目录
            if mv "$plugins_dir/$file_name" "$unused_plugin_dir/"; then
                echo "$info 插件文件 $file_name 已移动到 $unused_plugin_dir"
                return 0
            else
                echo "$error 移动插件 $file_name 失败"
                return 1
            fi
        else
            # 没有文件就不移动了
            return 0
        fi
    else
        echo "$error 找不到 $plugins_dir 目录"
        return 1
    fi
}

move_to_plugins() {
    local file_name="$1"
    local unused_plugin_dir="./unused_plugins"
    
    # 检查 ../plugins 目录是否存在，不存在则创建
    if [ ! -d "$plugins_dir" ]; then
        mkdir "$plugins_dir"
        echo "$info 目录 $plugins_dir 已创建"
    fi
    
    # 检查 unused_plugin 目录是否存在
    if [ -d "$unused_plugin_dir" ]; then
        # 检查文件是否存在
        if [ -f "$unused_plugin_dir/$file_name" ]; then
            # 尝试移动文件到 plugins 目录
            if mv "$unused_plugin_dir/$file_name" "$plugins_dir/"; then
                echo "$info 插件文件 $file_name 已移动到 $plugins_dir"
                return 0
            else
                echo "$error 移动插件 $file_name 失败"
                return 1
            fi
        else
            # 文件不存在就不移动了
            return 0
        fi
    else
        echo "$error 找不到 $unused_plugin_dir 目录"
        return 1
    fi
}

clear_config_sections() {
    # 清除 database_driver_options 部分
    local database_driver_line=$(grep -n "^  database_driver_options:" "$config_file" | cut -d: -f1)
    if [[ -n "$database_driver_line" ]]; then
        # 找到 database_driver_options 块的结束位置
        local end_line=$(sed -n "$((database_driver_line+1)),\$p" "$config_file" | grep -n -m 1 "^\(  \)\{0,1\}[^ ]" | cut -d: -f1)
        if [[ -z "$end_line" ]]; then
            end_line=$(wc -l < "$config_file")
        else
            end_line=$((database_driver_line + end_line))
        fi
        # 确保 end_line 不超出文件末尾
        end_line=$(($end_line < $(wc -l < "$config_file") ? $end_line : $(wc -l < "$config_file")))
        # 删除 database_driver_options 部分
        sed -i "${database_driver_line},$((end_line-1))d" "$config_file"
    fi
    
    # 清除 sqle 下的 plugin_config 部分
    sqle_line=$(grep -n "^sqle:" "$config_file" | cut -d: -f1)
    if [[ -n "$sqle_line" ]]; then
        plugin_config_line=$(sed -n "$sqle_line,\$p" "$config_file" | grep -n "^    plugin_config:" | cut -d: -f1)
        if [[ -n "$plugin_config_line" ]]; then
            # 计算绝对行号
            plugin_config_line=$((sqle_line + plugin_config_line - 1))
            # 找到 plugin_config 块的结束位置 行开头匹配0个2个或者4个空格
            end_line=$(sed -n "$((plugin_config_line+1)),\$p" "$config_file" | grep -n -m 1 "^\(  \)\{0,2\}[^ ]" | cut -d: -f1)
            if [[ -z "$end_line" ]]; then
                end_line=$(wc -l < "$config_file")
            else
                end_line=$((plugin_config_line + end_line))
            fi
            # 删除 plugin_config 部分
            sed -i "${plugin_config_line},$((end_line-1))d" "$config_file"
        fi
    fi
}

# 添加配置到database_driver_options层级下
install_plugin_to_config() {
    
    # 接受参数：db_type、logo_path 和 params
    db_type=$1
    logo_path=$2
    params=$3
    if check_config_for_content "db_type: $db_type" ; then
        echo "$info 插件配置 plugin_config 配置完成"
        return 0
    fi
    # 构造要添加的内容
    new_content="    - db_type: $db_type\n      logo_path: \"$logo_path\"\n      params: $params"
    
    # 插入到 database_driver_options: 层级下
    if grep -q "^  database_driver_options:" "$config_file"; then
        # 查找 "  database_driver_options:" 并在其后面插入内容
        sed -i "/^  database_driver_options:/a\\$new_content" "$config_file"
        echo "$info 插件配置 database_driver_options 配置完成"
    else
        echo "$error 错误: 找不到 'database_driver_options:' 关键字，未能配置插件 $db_type"
        return 1
    fi
    return 0
}

# 添加配置到plugin_config层级下
install_plugin_to_plugin_config() {
    # 接受参数：plugin_name 和 plugin_params
    plugin_name=$1
    cmd=$2
    if check_config_for_content "plugin_name: $plugin_name" ; then
        echo "$info 插件配置 plugin_config 配置完成"
        return 0
    fi
    # 构造要添加的内容
    new_content="      - plugin_name: $plugin_name\n        cmd: $cmd"
    
    # 插入到 plugin_config: 层级下
    if grep -q "^    plugin_config:" "$config_file"; then
        # 查找 "    plugin_config:" 并在其后面插入内容
        sed -i "/^    plugin_config:/a\\$new_content" "$config_file"
        echo "$info 插件配置 plugin_config 配置完成"
    else
        echo "$info 错误: 找不到 'plugin_config:' 关键字，未能配置插件 $plugin_name"
        return 1
    fi
    return 0
}

# 函数：检查并添加 plugin_config 到 sqle 的 service 层级下
add_plugin_config_to_sqle_service() {
    # 找到 sqle: 的行号
    sqle_line=$(grep -n "^sqle:" "$config_file" | cut -d: -f1)
    
    if [[ -z "$sqle_line" ]]; then
        echo "$info 错误: 找不到 'sqle:' 关键字"
        return 1
    fi
    
    # 确保在 sqle: 部分找到了 service:
    service_line=$(sed -n "$sqle_line,\$p" "$config_file" | grep -n "^  service:" | head -n 1 | cut -d: -f1)
    
    if [[ -z "$service_line" ]]; then
        echo "$info 错误: 找不到 'service:' 关键字"
        return 1
    fi
    
    # 计算绝对行号
    service_line=$((sqle_line + service_line - 1))
    
    # 检查 plugin_config 是否已经存在于 service: 下
    if ! sed -n "$service_line,\$p" "$config_file" | grep -q "^    plugin_config:"; then
        # 插入 plugin_config 到 service: 下
        sed -i "${service_line}a\\    plugin_config:" "$config_file"
        echo "$info 自动在配置文件sqle.service下添加plugin_config配置项"
    fi
    
    return 0
}

# 函数：检查并添加 database_driver_options 到 dms 层级下
add_plugin_config_to_dms() {
    # 检查 dms: 下是否已经存在 database_driver_options:
    if grep -q "^dms:" "$config_file"; then
        if ! grep -q "^  database_driver_options:" "$config_file"; then
            # 如果 database_driver_options 不存在，则添加
            sed -i "/^dms:/a\\  database_driver_options:" "$config_file"
        fi
    else
        echo "$info 错误: 找不到 'dms:' 关键字"
        return 1
    fi
    return 0
}

# 函数：检查是否存在指定的内容
check_config_for_content() {
    search_content=$1
    if grep -q "$search_content" "$config_file"; then
        echo "$warn 插件配置 $search_content 存在于配置文件中"
        echo "$warn 该部分配置未覆盖，也未添加新的配置配置"
        echo "$warn 若有手动修改，请人工确认配置的正确性"
        return 0
    else
        return 1
    fi
}

launch_sqle() {
    echo "$split_line"
    echo ">>>> 正在启动SQLE <<<<"
    echo "$split_line"
    if chmod +x "$work_directory/scripts/init_start.sh"; then
        echo "$info init_start.sh 可执行权限配置完成"
    else
        echo "$error 设置 init_start.sh 执行权限失败"
        return 1
    fi
    # 检查 init_start.sh 是否存在
    if [ ! -f "$work_directory/scripts/init_start.sh" ]; then
        echo "$error init_start.sh 脚本不存在"
        echo "$split_line"
        exit 1
    fi
    
    # 启动 init_start.sh 脚本
    echo "请稍等......"
    bash "$work_directory/scripts/init_start.sh"
    # 检查启动是否成功
    if [ $? -eq 0 ]; then
        echo "$info SQLE启动成功"
        echo "$split_line"
        echo ">>>>> SQLE 为您提供 SQL 全生命周期质量管控能力 <<<<<"
        echo ""
        echo "爱可生官网-SQLE详情页见: "
        echo "https://www.actionsky.com/sqle"
        echo ""
        echo "快速开始请参考："
        echo "https://actiontech.github.io/sqle-docs/docs/quick-usage"
        echo ""
        echo "用户使用手册请参考: "
        echo "https://actiontech.github.io/sqle-docs/docs/user-manual/project/intro"
        echo ""
        echo "若需要社区支持，欢迎加入SQLE社区交流群："
        echo "https://actiontech.github.io/sqle-docs/docs/support/community-support"
        echo ""
        echo "$split_line"
    else
        echo "$error SQLE启动失败"
        echo "$split_line"
        echo ">>>>> SQLE 为您提供 SQL 全生命周期质量管控能力 <<<<<"
        echo ""
        echo "若需要社区支持，欢迎加入SQLE社区交流群："
        echo "https://actiontech.github.io/sqle-docs/docs/support/community-support"
        echo ""
        echo "$split_line"
    fi
}

restart_sqle() {
  echo "$split_line"
  echo ">>>> 正在重新启动SQLE <<<<"
  echo "$split_line"
  if pgrep -x sqled >/dev/null
  then
    if chmod +x "$work_directory/scripts/stop_sqle.sh"; then
          echo "$info stop_sqle.sh 可执行权限配置完成"
      else
          echo "$error 设置 stop_sqle.sh 执行权限失败"
          return 1
    fi
    # 检查 stop_sqle.sh 是否存在
    if [ ! -f "$work_directory/scripts/stop_sqle.sh" ]; then
            echo "$error stop_sqle.sh 脚本不存在"
            echo "$split_line"
            exit 1
    fi
    # 启动 stop_sqle.sh脚本
    bash "$work_directory/scripts/stop_sqle.sh"
    launch_sqle
  else
    launch_sqle
  fi
}

