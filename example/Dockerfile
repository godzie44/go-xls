FROM golang:1.16 as builder

WORKDIR /app
COPY . .


RUN wget https://github.com/libxls/libxls/releases/download/v1.6.2/libxls-1.6.2.tar.gz \
    && tar -xzvf libxls-1.6.2.tar.gz \
    && cd libxls-1.6.2 \
    && ./configure \
    && make \
    && make install \
    && ldconfig /etc/ld.so.conf.d

RUN CGO_LDFLAGS=-lxlsreader GOOSE=linux GO111MODULE=on go build -a -o main /app/example/main.go
CMD ["./main"]
