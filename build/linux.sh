export GOOS="linux"
export GOARCH="amd64"
go get "gopkg.in/yaml.v2"
go build -o ../dist/linux/azmigrate ../main.go
chmod +x ../dist/linux/azmigrate