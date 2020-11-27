Summary: Actiontech Sqle_sqlserver
Name: sqle_sqlserver
Version: 4.20.11.0_%{commit}
%if %{?_with_qa:1}%{!?_with_qa:0}
Release: qa.%{os_version}
%else
Release: rel.%{os_version}
%endif
Source0: %{name}.tar.gz
License: Commercial
Group: Actiontech
Prefix: /usr/local/sqle

%description
Acitontech Sqle_sqlserver

%define debug_package %{nil}
%define _source_filedigest_algorithm md5
%define _binary_filedigest_algorithm md5
%define _source_payload w0.gzdio
%define _binary_payload w0.gzdio
%{!?_dotnet_target: %define _dotnet_target centos.7-x64}
%{!?_git_version: %define _git_version master}
##########

%prep
%setup -q

##########

%build
echo "build sqle_sqlserver..."
cd %{_builddir}/%{buildsubdir}/sqle
make build_sqlserver

##########

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/usr/local/sqle/libs
cp -R %{_builddir}/%{buildsubdir}/sqle/scripts $RPM_BUILD_ROOT/usr/local/sqle/scripts
cp -R %{_builddir}/%{buildsubdir}/sqle/sqle/sqlserver/SqlserverProtoServer/bin/Release/netcoreapp2.1/%{_dotnet_target}/* $RPM_BUILD_ROOT/usr/local/sqle/libs
echo "%{_git_version}" > $RPM_BUILD_ROOT/usr/local/sqle/VERSION

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
(getent group %{group_name} 1>/dev/null 2>&1) || groupadd -g 5700 %{group_name}
(id %{user_name} 1>/dev/null 2>&1) || useradd -M -g %{group_name} -s $(which nologin) -u 5700 %{user_name}


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
    sed -e "s|PIDFile=|PIDFile=$RPM_INSTALL_PREFIX\/sqle_sqlserver.pid|g" \
    -e "s|User=|User=actiontech-universe|g" \
    -e "s|ExecStart=|ExecStart=/usr/bin/dotnet $RPM_INSTALL_PREFIX\/libs\/publish\/SqlserverProtoServer.dll --config $RPM_INSTALL_PREFIX\/etc\/sqle_sqlserver.cnf --pidfile=$RPM_INSTALL_PREFIX\/sqle_sqlserver.pid|g" \
    -e "s|WorkingDirectory=|WorkingDirectory=$RPM_INSTALL_PREFIX|g" \
    $RPM_INSTALL_PREFIX/scripts/sqled.systemd > /lib/systemd/system/sqle_sqlserver.service
    systemctl daemon-reload
    systemctl enable sqle_sqlserver.service

else
    sed -e "s|PROJECT_PATH=|PROJECT_PATH=$RPM_INSTALL_PREFIX|g" \
    $RPM_INSTALL_PREFIX/scripts/sqled.initd > /etc/init.d/sqle_sqlserver
    chmod 755 /etc/init.d/sqle_sqlserver
    chkconfig --add sqle_sqlserver
fi

mkdir -p $RPM_INSTALL_PREFIX/etc

cat > $RPM_INSTALL_PREFIX/etc/sqle_sqlserver.cnf.template<<EOF
[server]
port=10001
EOF

#chown
chown -R actiontech-universe: $RPM_INSTALL_PREFIX

#chmod
find $RPM_INSTALL_PREFIX -type d -exec chmod 0750 {} \;
find $RPM_INSTALL_PREFIX -type f -exec chmod 0640 {} \;
chmod 0770 $RPM_INSTALL_PREFIX/etc

##########

%preun

if [ "$1" = "0" ]; then
    grep systemd /proc/1/comm 1>/dev/null 2>&1
    if [ $? -eq 0 ]; then
        systemctl stop sqle_sqlserver.service || true
    else
        service sqle_sqlserver stop || true
    fi
fi

##########

%postun

if [ "$1" = "0" ]; then
    grep systemd /proc/1/comm 1>/dev/null 2>&1
    if [ $? -eq 0 ]; then
        systemctl disable sqle_sqlserver.service || true
        rm -f /lib/systemd/system/sqle_sqlserver.service || true
        systemctl daemon-reload
        systemctl reset-failed sqle_sqlserver.service || true
    else
        chkconfig --del sqle_sqlserver || true
        rm -f /etc/init.d/sqle_sqlserver || true
    fi
fi


##########

%files
%defattr(-,root,root)
/usr/local/sqle/libs
/usr/local/sqle/scripts
/usr/local/sqle/VERSION
