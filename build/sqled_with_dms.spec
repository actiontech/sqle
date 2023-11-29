Summary: Actiontech Sqle
Name: sqle
Version: %{commit}
%if %{?_with_qa:1}%{!?_with_qa:0}
Release: qa.%{os_version}
%else
Release: rel.%{os_version}
%endif
Source0: %{name}.tar.gz
License: Commercial
Group: Actiontech
Prefix: /usr/local/%{name}

%description
Acitontech Sqle

%define debug_package %{nil}
%define _source_filedigest_algorithm md5
%define _binary_filedigest_algorithm md5
%define _source_payload w0.gzdio
%define _binary_payload w0.gzdio

##########

%prep
%setup -q

##########

##########

%build
## build is done in outside, please see Makefile.

##########

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/usr/local/%{name}/bin
mkdir -p $RPM_BUILD_ROOT/usr/local/%{name}/etc
mkdir -p %{_builddir}/%{buildsubdir}/%{name}/plugins
cp %{_builddir}/%{buildsubdir}/dms/config.yaml $RPM_BUILD_ROOT/usr/local/%{name}/etc/config.yaml
cp -R %{_builddir}/%{buildsubdir}/dms/build/logo $RPM_BUILD_ROOT/usr/local/%{name}/static/
cp %{_builddir}/%{buildsubdir}/%{name}/bin/sqled $RPM_BUILD_ROOT/usr/local/%{name}/bin/sqled
cp %{_builddir}/%{buildsubdir}/%{name}/bin/scannerd $RPM_BUILD_ROOT/usr/local/%{name}/bin/scannerd
cp %{_builddir}/%{buildsubdir}/%{name}/bin/dms $RPM_BUILD_ROOT/usr/local/%{name}/bin/dms
cp -R %{_builddir}/%{buildsubdir}/%{name}/plugins $RPM_BUILD_ROOT/usr/local/%{name}/plugins
cp -R %{_builddir}/%{buildsubdir}/%{name}/scripts $RPM_BUILD_ROOT/usr/local/%{name}/scripts
cp %{_builddir}/%{buildsubdir}/dms/build/service-file-template/dms.systemd $RPM_BUILD_ROOT/usr/local/%{name}/scripts/dms.systemd
cp -R %{_builddir}/%{buildsubdir}/%{name}/static/* $RPM_BUILD_ROOT/usr/local/%{name}/static/


##########

%clean
#rm -rf $RPM_BUILD_ROOT

##########

%pre

#check directory
grep systemd /proc/1/comm 1>/dev/null 2>&1
if [ $? -eq 0 ]; then
    if [ ! -d "/lib/systemd/system" ];then
        mkdir -p /lib/systemd/system
        chmod 0755 /lib/systemd/system
    fi
fi

#create group & user
(which nologin 1>/dev/null 2>&1) || (echo "require nologin" && exit 11)
(which bash 1>/dev/null 2>&1) || (echo "require bash" && exit 12)
(which pkill 1>/dev/null 2>&1) || (echo "require pkill" && exit 13)
(getent group %{group_name} 1>/dev/null 2>&1) || groupadd -g 5700 %{group_name}
(id %{user_name} 1>/dev/null 2>&1) || (useradd -M -g %{group_name} -s $(which nologin) -u 5700 %{user_name} && chage -M 99999 %{user_name})


#check bash env
bash -c "" 2>&1 | grep -e 'warning' -e 'error'
if [ $? -eq 0 ]; then
    exit 14
fi

##########

%post

#service
grep systemd /proc/1/comm 1>/dev/null 2>&1
if [ $? -eq 0 ]; then
    sed -e "s|PIDFile=|PIDFile=$RPM_INSTALL_PREFIX\/sqled.pid|g" \
    -e "s|User=|User=actiontech-universe|g" \
    -e "s|ExecStart=|ExecStart=/bin/sh -c 'exec $RPM_INSTALL_PREFIX\/bin\/sqled --config $RPM_INSTALL_PREFIX\/etc\/config.yaml --pidfile=$RPM_INSTALL_PREFIX\/sqled.pid >>$RPM_INSTALL_PREFIX\/std.log 2>\&1'|g" \
    -e "s|WorkingDirectory=|WorkingDirectory=$RPM_INSTALL_PREFIX|g" \
    $RPM_INSTALL_PREFIX/scripts/sqled.systemd > /lib/systemd/system/sqled.service
    sed -e "s|PIDFile=|PIDFile=$RPM_INSTALL_PREFIX\/dms.pid|g" -e "s|User=|User=actiontech-universe|g" -e "s|ExecStart=|ExecStart=$RPM_INSTALL_PREFIX\/bin\/dms -conf $RPM_INSTALL_PREFIX\/etc\/config.yaml|g" -e "s|WorkingDirectory=|WorkingDirectory=$RPM_INSTALL_PREFIX|g" $RPM_INSTALL_PREFIX/scripts/dms.systemd > /lib/systemd/system/dms.service
    systemctl daemon-reload
    systemctl enable sqled.service
    systemctl enable dms.service
fi

mkdir -p $RPM_INSTALL_PREFIX/logs
# mkdir -p $RPM_INSTALL_PREFIX/etc


cat >> $RPM_INSTALL_PREFIX/etc/config.yaml<<EOF

sqle:
  id: 1
  dms_server_address: http://127.0.0.1:7601
  api:
    addr: 127.0.0.1
    port: 10000
    enable_https: false
    cert_file_path: './etc/cert.pem'
    key_file_path: './etc/key.pem'  
  secret_key:     
  service:
    auto_migrate_table: true   
    debug_log: false
    log_path: './logs'
    log_max_size_mb: 1024
    log_max_backup_number: 2
    plugin_path: './plugins'    
    enable_cluster_mode:
    database:
      mysql_host: '127.0.0.1'
      mysql_port: '3306'
      mysql_user: 'root'
      mysql_password: 'pass'
      mysql_schema: 'sqle'
EOF


cat > $RPM_INSTALL_PREFIX/etc/gh-ost.ini<<EOF
max_load=Threads_running=80,Threads_connected=1000
critical_load=Threads_running=160,Threads_connected=2000
chunk_size=1000
dml_batch_size=10
default_retries=120
exponential_backoff_max_interval=64
throttle_additional_flag_file=/tmp/gh-ost.throttle
cut_over_lock_timeout_seconds=3
max_lag_millis=1500
heartbeat_interval_millis=100
EOF

#chown
chown -R %{user_name}: $RPM_INSTALL_PREFIX

#chmod
find $RPM_INSTALL_PREFIX -type d -exec chmod 0750 {} \;
find $RPM_INSTALL_PREFIX -type f -exec chmod 0640 {} \;
chmod 0750 $RPM_INSTALL_PREFIX/bin/*
find $RPM_INSTALL_PREFIX/plugins -type f -exec chmod 0750 {} \;
chmod 0770 $RPM_INSTALL_PREFIX/etc

function box_out()
{
  local s=("$@") b w
  for l in "${s[@]}"; do
    ((w<${#l})) && { b="$l"; w="${#l}"; }
  done
  tput setaf 3
  echo " -${b//?/-}-
| ${b//?/ } |"
  for l in "${s[@]}"; do
    printf '| %s%*s%s |\n' "$(tput setaf 2)" "-$w" "$l" "$(tput setaf 3)"
  done
  echo "| ${b//?/ } |
 -${b//?/-}-"
  tput sgr 0
}

box_out "To start the service, please build $RPM_INSTALL_PREFIX/etc/config.yaml, then run script ./scripts/init_start.sh in the $RPM_INSTALL_PREFIX directory" "Example: 1: vi $RPM_INSTALL_PREFIX/etc/config.yaml 2: sh $RPM_INSTALL_PREFIX/scripts/init_start.sh"

##########

%preun

if [ "$1" = "0" ]; then
    grep systemd /proc/1/comm 1>/dev/null 2>&1
    if [ $? -eq 0 ]; then
        systemctl stop sqled.service || true
    else
        service sqled stop || true
    fi
fi

function kill_and_wait {
        pidfile=$1
        if [ -e $pidfile ]; then
                kill $(cat $pidfile) &>/dev/null
        fi
        for i in {1..60}; do
                if [ ! -e $pidfile ]; then
                        return 0
                fi
                kill -0 $(cat $pidfile) &>/dev/null
                if [ $? -ne 0 ]; then
                        return 0
                fi
                sleep 1
        done
        return 1
}

if [ "$1" = "0" ]; then
    kill_and_wait $RPM_INSTALL_PREFIX/sqled.pid
    if [ $? -ne 0 ]; then
    	(>&2 echo "wait sqled pid shutdown timeout")
    	exit 1
    fi
    kill_and_wait $RPM_INSTALL_PREFIX/dms.pid
    if [ $? -ne 0 ]; then
    	(>&2 echo "wait dms pid shutdown timeout")
    	exit 1
    fi    
fi

##########

%postun

if [ "$1" = "0" ]; then
    grep systemd /proc/1/comm 1>/dev/null 2>&1
    if [ $? -eq 0 ]; then
        systemctl disable sqled.service || true
        rm -f /lib/systemd/system/sqled.service || true
        systemctl disable dms.service || true
        rm -f /lib/systemd/system/dms.service || true
        systemctl daemon-reload
        systemctl reset-failed sqled.service || true
        systemctl reset-failed dms.service || true
    else
        chkconfig --del sqled || true
        rm -f /etc/init.d/sqled || true
        chkconfig --del dms || true
        rm -f /etc/init.d/dms || true        
    fi
fi

##########

%files
%defattr(-,root,root)
/usr/local/%{name}/bin/sqled
/usr/local/%{name}/bin/scannerd
/usr/local/%{name}/bin/dms
/usr/local/%{name}/plugins
/usr/local/%{name}/scripts/*
/usr/local/%{name}/static/* 
/usr/local/%{name}/etc/config.yaml


%config /usr/local/%{name}/etc/config.yaml
