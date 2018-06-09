#FROM golang:latest 
#RUN mkdir /app 
#ADD . /app/ 
#WORKDIR /app 
#RUN go build -o main . 
#CMD ["/app/main"]


FROM resin/raspberrypi3-golang as builder
WORKDIR /app
RUN go get -d -v golang.org/x/net/html  
COPY app.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .    

FROM scratch  
WORKDIR /root/
COPY --from=builder /app .
CMD ["./app"]  
