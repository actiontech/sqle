Summary: Actiontech Sqle
Name: sqle
Version: 9.9.9.9
%if %{?_with_qa:1}%{!?_with_qa:0}
Release: qa
%else
Release: rel
%endif
Source0: %{name}.tar.gz
License: Commercial
Group: Actiontech
Prefix: /usr/local/sqle

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

%build
echo "build sqle..."
export GOPATH=%{_builddir}/%{buildsubdir}
cd %{_builddir}/%{buildsubdir}/src/sqle
make build

##########

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/usr/local/sqle/bin
cp %{_builddir}/%{buildsubdir}/src/sqle/sqled $RPM_BUILD_ROOT/usr/local/sqle/bin/sqled
cp -R %{_builddir}/%{buildsubdir}/src/sqle/scripts $RPM_BUILD_ROOT/usr/local/sqle/scripts

##########

%clean
#rm -rf $RPM_BUILD_ROOT

##########

%pre
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
(getent group actiontech 1>/dev/null 2>&1) || groupadd -g 5700 actiontech
(id actiontech-sqle 1>/dev/null 2>&1) || useradd -M -g actiontech -s $(which nologin) actiontech-sqle

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
    -e "s|User=|User=actiontech-sqle|g" \
    -e "s|ExecStart=|ExecStart=$RPM_INSTALL_PREFIX\/bin\/sqled --config $RPM_INSTALL_PREFIX\/etc\/sqled.cnf --pidfile=$RPM_INSTALL_PREFIX\/sqled.pid|g" \
    -e "s|WorkingDirectory=|WorkingDirectory=$RPM_INSTALL_PREFIX|g" \
    $RPM_INSTALL_PREFIX/scripts/sqled.systemd > /lib/systemd/system/sqled.service
    systemctl daemon-reload
    systemctl enable sqled.service

else
    sed -e "s|PROJECT_PATH=|PROJECT_PATH=$RPM_INSTALL_PREFIX|g" \
    $RPM_INSTALL_PREFIX/scripts/sqled.initd > /etc/init.d/sqled
    chmod 755 /etc/init.d/sqled
    chkconfig --add sqled
fi

mkdir -p $RPM_INSTALL_PREFIX/logs
mkdir -p $RPM_INSTALL_PREFIX/etc

cat > $RPM_INSTALL_PREFIX/etc/sqled.cnf.template<<EOF
[server]
port=
mysql_host=
mysql_port=
mysql_user=
mysql_password=
mysql_schema=
log_path=./logs
#
auto_migrate_table=true
debug=false
EOF

#chown
chown -R actiontech-sqle: $RPM_INSTALL_PREFIX

#chmod
find $RPM_INSTALL_PREFIX -type d -exec chmod 0750 {} \;
find $RPM_INSTALL_PREFIX -type f -exec chmod 0640 {} \;
chmod 0750 $RPM_INSTALL_PREFIX/bin/*
chmod 0770 $RPM_INSTALL_PREFIX/etc

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

##########

%postun

if [ "$1" = "0" ]; then
    grep systemd /proc/1/comm 1>/dev/null 2>&1
    if [ $? -eq 0 ]; then
        systemctl disable sqled.service || true
        rm -f /lib/systemd/system/sqled.service || true
        systemctl daemon-reload
        systemctl reset-failed sqled.service || true
    else
        chkconfig --del sqled || true
        rm -f /etc/init.d/sqled || true
    fi
fi


##########

%files
%defattr(-,root,root)
/usr/local/sqle/bin/sqled
/usr/local/sqle/scripts/sqled.systemd
/usr/local/sqle/scripts/sqled.initd
