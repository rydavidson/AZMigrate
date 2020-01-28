export GOOS="windows"
export GOARCH="amd64"
go get "gopkg.in/yaml.v2"
go build -o azmigrate ../main.go