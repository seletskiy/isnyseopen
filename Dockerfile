FROM golang:1.13 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 go get -v

FROM scratch
LABEL authors="Stanislav Seletskiy <s.seletskiy@gmail.com>"

COPY --from=builder /go/bin/isnyseopen /bin/isnyseopen
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE 8080
ENTRYPOINT ["/bin/isnyseopen"]
