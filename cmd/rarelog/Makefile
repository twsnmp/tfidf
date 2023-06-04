.PHONY: all clean zip hash

### バージョンの定義
VERSION     := "v1.0.0"
COMMIT      := $(shell git rev-parse --short HEAD)

### コマンドの定義
GO          = go
GO_BUILD    = $(GO) build
GO_LDFLAGS  = -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"
ZIP          = zip

### ターゲットパラメータ
DIST = ./dist

### PHONY ターゲットのビルドルール
all: zip
clean:
	rm -rf ./dist
zip: $(DIST)/rarelog $(DIST)/rarelog.app $(DIST)/rarelog.exe 
	cd $(DIST) && $(ZIP) rarelog.zip rarelog*

hash:
	cd ./dist && shasum -a 256 rarelog.zip

### 実行ファイルのビルドルール
$(DIST)/rarelog.exe: main.go $(DIST)
	env GOOS=windows GOARCH=amd64 $(GO_BUILD) $(GO_LDFLAGS) -o $@
$(DIST)/rarelog.app: main.go $(DIST)
	env GOOS=darwin GOARCH=amd64 $(GO_BUILD) $(GO_LDFLAGS) -o $@
$(DIST)/rarelog: main.go $(DIST)
	env GOOS=linux GOARCH=amd64 $(GO_BUILD) $(GO_LDFLAGS) -o $@

$(DIST):
	mkdir -p $(DIST)
