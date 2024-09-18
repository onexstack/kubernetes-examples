FROM debian:trixie

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY aggregated-apiserver /opt/aggregated-apiserver/aggregated-apiserver

ENTRYPOINT ["/opt/aggregated-apiserver/aggregated-apiserver"]


