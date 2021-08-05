FROM alpine:3.12
ENV TZ Asia/Shanghai

RUN apk add --no-cache \
        tzdata \
        && rm -f /etc/localtime \
        && ln -s /usr/share/zoneinfo/$TZ /etc/localtime \
        && echo $TZ > /etc/timezone

RUN apk add --update --no-cache \
    ca-certificates \
    && rm -rf /var/cache/apk/*

COPY jobsvr /usr/local/bin/jobsvr
RUN  chmod +x /usr/local/bin/jobsvr

WORKDIR /usr/local/var/jobsvr
EXPOSE 80