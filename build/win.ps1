$Env:GOOS="windows"
$Env:GOARCH="amd64"
go get "gopkg.in/yaml.v2"
go build -o azmigrate.exe ../main.go