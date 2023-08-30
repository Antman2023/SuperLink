.PHONY: all

all:
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-w -s" -o superlink.exe