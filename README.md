go-docker

# Build
## Build for pi zero 
```
GOARCH=arm GOARM=6 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app . 
```
