FROM golang
RUN mkdir -p /app
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o app

#Change the timezone to Asia
RUN apt-get update \
    &&  DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends tzdata
RUN TZ=Asia/Taipei \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone \
    && dpkg-reconfigure -f noninteractive tzdata

# GET AB for test
# RUN apt-get install apache2-utils

ENTRYPOINT ["./app"]