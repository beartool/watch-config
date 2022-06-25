PROJECTNAME = go

## linux: 编译打包linux
.PHONY: linux
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(RACE) -o watch-file ./main.go

## win: 编译打包win
.PHONY: win
win:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(RACE) -o watch-file.exe ./main.go
