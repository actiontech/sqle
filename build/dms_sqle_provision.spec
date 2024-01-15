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
mkdir -p $RPM_BUILD_ROOT/usr/local/%{name}/plugins
cp -R %{_builddir}/%{buildsubdir}/%{name}/builddir/config $RPM_BUILD_ROOT/usr/local/%{name}/etc
cp -R %{_builddir}/%{buildsubdir}/%{name}/builddir/bin $RPM_BUILD_ROOT/usr/local/%{name}/bin
cp -R %{_builddir}/%{buildsubdir}/%{name}/builddir/scripts $RPM_BUILD_ROOT/usr/local/%{name}/scripts
cp -R %{_builddir}/%{buildsubdir}/%{name}/builddir/static $RPM_BUILD_ROOT/usr/local/%{name}/static
cp -R %{_builddir}/%{buildsubdir}/%{name}/builddir/neo4j-community $RPM_BUILD_ROOT/usr/local/%{name}/neo4j-community
cp -R %{_builddir}/%{buildsubdir}/%{name}/builddir/lib $RPM_BUILD_ROOT/usr/local/%{name}/lib

##########

%files
%defattr(-,root,root)
/usr/local/%{name}/bin/*
/usr/local/%{name}/plugins
/usr/local/%{name}/scripts/*
/usr/local/%{name}/static/*
/usr/local/%{name}/neo4j-community/*
/usr/local/%{name}/lib/*
/usr/local/%{name}/etc/config.yaml
%config /usr/local/%{name}/etc/config.yaml

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

#create ulimit
if [ ! -d "/etc/security/limits.d" ];then
    mkdir /etc/security/limits.d
    chmod 0755 /etc/security/limits.d
fi

cat > /etc/security/limits.d/%{name}.conf <<EOF
%{user_name}     soft    nofile    65535
%{user_name}     hard    nofile    65535
%{user_name}     soft    nproc     65535
%{user_name}     hard    nproc     65535
EOF
chown root: /etc/security/limits.d/%{name}.conf
chmod 440 /etc/security/limits.d/%{name}.conf

#check bash env
bash -c "" 2>&1 | grep -e 'warning' -e 'error'
if [ $? -eq 0 ]; then
    exit 14
fi

##########

%post

cat >> $RPM_INSTALL_PREFIX/neo4j-community/conf/neo4j.conf <<EOF
dbms.security.auth_enabled=false
dbms.default_listen_address=0.0.0.0
EOF

#service
grep systemd /proc/1/comm 1>/dev/null 2>&1
if [ $? -eq 0 ]; then
    sed -e "s|PIDFile=|PIDFile=$RPM_INSTALL_PREFIX\/sqled.pid|g" \
    -e "s|User=|User=actiontech-universe|g" \
    -e "s|ExecStart=|ExecStart=/bin/sh -c 'exec $RPM_INSTALL_PREFIX\/bin\/sqled --config $RPM_INSTALL_PREFIX\/etc\/config.yaml --pidfile=$RPM_INSTALL_PREFIX\/sqled.pid >>$RPM_INSTALL_PREFIX\/std.log 2>\&1'|g" \
    -e "s|WorkingDirectory=|WorkingDirectory=$RPM_INSTALL_PREFIX|g" \
    $RPM_INSTALL_PREFIX/scripts/sqled.systemd > /lib/systemd/system/sqled.service
    sed -e "s|PIDFile=|PIDFile=$RPM_INSTALL_PREFIX\/dms.pid|g" \
    -e "s|User=|User=actiontech-universe|g" \
    -e "s|ExecStart=|ExecStart=$RPM_INSTALL_PREFIX\/bin\/dms -conf $RPM_INSTALL_PREFIX\/etc\/config.yaml|g" \
    -e "s|WorkingDirectory=|WorkingDirectory=$RPM_INSTALL_PREFIX|g" \
    $RPM_INSTALL_PREFIX/scripts/dms.systemd > /lib/systemd/system/dms.service
    sed -e "s|PIDFile=|PIDFile=$RPM_INSTALL_PREFIX\/provision.pid|g" \
    -e "s|ExecStart=|ExecStart=$RPM_INSTALL_PREFIX\/bin\/provision -conf $RPM_INSTALL_PREFIX\/etc\/config.yaml|g" \
    -e "s|WorkingDirectory=|WorkingDirectory=$RPM_INSTALL_PREFIX|g" \
    $RPM_INSTALL_PREFIX/scripts/provision.systemd > /lib/systemd/system/provision.service
    sed -e "s|PIDFile=|PIDFile=$RPM_INSTALL_PREFIX\/neo4j-community/run/neo4j.pid|g" \
    -e "s|ExecStart=|ExecStart=$RPM_INSTALL_PREFIX\/neo4j-community/bin/neo4j start|g" \
    -e "s|ExecStop=|ExecStop=$RPM_INSTALL_PREFIX\/neo4j-community/bin/neo4j stop|g" \
    -e "s|ExecReload=|ExecReload=$RPM_INSTALL_PREFIX\/neo4j-community/bin/neo4j restart|g" \
    -e "s|WorkingDirectory=|WorkingDirectory=$RPM_INSTALL_PREFIX/neo4j-community|g" \
    $RPM_INSTALL_PREFIX/scripts/neo4j.systemd > /lib/systemd/system/neo4j.service
    systemctl daemon-reload
    systemctl enable provision.service
    systemctl enable neo4j.service
    systemctl enable sqled.service
    systemctl enable dms.service
fi

mkdir -p $RPM_INSTALL_PREFIX/logs
# mkdir -p $RPM_INSTALL_PREFIX/etc

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
find $RPM_INSTALL_PREFIX/plugins -type f -exec chmod 0750 {} \;
chmod 0750 $RPM_INSTALL_PREFIX/bin/*
chmod 0750 $RPM_INSTALL_PREFIX/neo4j-community/bin/neo4j
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
        systemctl stop dms.service || true
        systemctl stop neo4j.service || true
        systemctl stop provision.service || true
    fi
fi

##########

%postun

if [ "$1" = "0" ]; then
    grep systemd /proc/1/comm 1>/dev/null 2>&1
    if [ $? -eq 0 ]; then
        rm -f /lib/systemd/system/sqled.service || true
        rm -f /lib/systemd/system/dms.service || true
        rm -f /lib/systemd/system/neo4j.service || true
        rm -f /lib/systemd/system/provision.service || true
    fi
fi


