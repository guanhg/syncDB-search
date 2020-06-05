FROM golang:1.14.1
WORKDIR /usr/src/statement-search-go
COPY . /usr/src/statement-search-go
ENV GOPROXY=https://goproxy.cn,direct
ENV ACTIVE=prd
ENV GIN_MODE=release
RUN export GO111MODULE=on
RUN  go build  -o web ./main/webApplication.go
RUN  go build  -o sync2db ./main/syncDb2Mq.go
RUN  go build  -o sync4mq ./main/syncElastic4Mq.go
EXPOSE 8080
ENTRYPOINT ["./web"]