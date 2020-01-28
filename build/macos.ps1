$Env:GOOS="darwin"
$Env:GOARCH="amd64"
go get "gopkg.in/yaml.v2"
go build -o ../dist/macos/azmigrate ../main.go