FROM golang:alpine as builder

RUN apk add --no-cache git gcc libc-dev

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go get "github.com/zegl/goriak/v3"
RUN go get "github.com/bkaradzic/go-lz4"
RUN go get "gopkg.in/mgo.v2"
RUN go get "github.com/montanaflynn/stats"
RUN go get "golang.org/x/perf/internal/stats"
RUN go get "github.com/go-redis/redis"
RUN go get "github.com/shamaton/msgpack"
#COPY $PWD/ /src/app/
WORKDIR /src/app/

#RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o /go/bin/runner

#FROM scratch

#COPY --from=builder /go/bin/runner /app/runner

#WORKDIR /app

CMD ["/bin/true"]