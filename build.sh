CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-w -s" -o voice-server-linux-arm7
CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -o voice-server-linux-mipsle
CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -o voice-server-linux-mips
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-w -s" -o voice-server-linux-x86
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o voice-server-linux-x64
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-w -s" -o voice-server-windows-x86.exe
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o voice-server-windows-x64.exe
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o voice-server-darwin-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o voice-server-darwin-arm64
