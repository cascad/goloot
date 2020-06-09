FROM golang:alpine as builder

RUN apk add --no-cache git gcc libc-dev

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

#COPY $PWD/ /src/app/
WORKDIR /src/app/
RUN go mod tidy

RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o /go/bin/runner

FROM scratch

COPY --from=builder /go/bin/runner /app/runner

WORKDIR /app

CMD ["/bin/true"]