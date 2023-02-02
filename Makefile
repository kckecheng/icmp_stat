files := main.go
name := icmp_stat

all: x86_64 arm windows

x86_64:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $(name).x86_64 $(files)

arm:
	CGO_ENABLED=0 GOARCH=arm GOOS=linux go build -o $(name).arm $(files)

windows:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o $(name).exe $(files)
