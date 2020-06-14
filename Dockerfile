FROM resin/raspberrypi3-golang:latest as builder

ENV INITSYSTEM on

WORKDIR /app
RUN go get -d -v github.com/ikester/blinkt 
RUN go get -d -v github.com/lucasb-eyer/go-colorful
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .    

FROM scratch  
WORKDIR /root/
COPY --from=builder /app .
ENTRYPOINT ["./app"]  


