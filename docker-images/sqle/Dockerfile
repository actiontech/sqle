FROM centos:7

RUN (cd /lib/systemd/system/sysinit.target.wants/; for i in *; do [ $i == \
systemd-tmpfiles-setup.service ] || rm -f $i; done); \
rm -f /lib/systemd/system/multi-user.target.wants/*;\
rm -f /etc/systemd/system/*.wants/*;\
rm -f /lib/systemd/system/local-fs.target.wants/*; \
rm -f /lib/systemd/system/sockets.target.wants/*udev*; \
rm -f /lib/systemd/system/sockets.target.wants/*initctl*; \
rm -f /lib/systemd/system/basic.target.wants/*;\
rm -f /lib/systemd/system/anaconda.target.wants/*;

RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

RUN yum install -y which

COPY ./sqle.rpm /opt/sqle.rpm
COPY ./docker-images/sqle/start.sh /opt/

RUN rpm -ivh /opt/sqle.rpm --prefix=/opt/sqle
RUN rm -f /opt/sqle.rpm

RUN mkdir -p /opt/sqle/plugins
RUN chown actiontech-universe:actiontech /opt/sqle/plugins

RUN chown actiontech-universe:actiontech /opt/start.sh
RUN chmod +x /opt/start.sh

ENV MYSQL_HOST ""
ENV MYSQL_PORT 3306
ENV MYSQL_USER  ""
ENV MYSQL_PASSWORD ""
ENV MYSQL_DMS_SCHEMA ""
ENV MYSQL_SQLE_SCHEMA ""

ENV DEBUG false
ENV AUTO_MIGRATE_TABLE true

USER actiontech-universe
WORKDIR /opt
CMD ["sh", "-c", "/opt/start.sh"]
