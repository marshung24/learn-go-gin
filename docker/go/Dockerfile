# 基底映像：Go 1.22（官方映像）
FROM golang:1.22

# 安裝開發常用工具
#   git, curl, unzip       — 基本開發工具
#   tini                   — 容器內的 init process，確保子程序能正常結束
#   iproute2, net-tools    — 網路除錯（ip, ifconfig, netstat 等指令）
#   vim, bash-completion   — 終端機編輯與自動補全
#   default-mysql-client   — 可在容器內直接用 mysql 指令連資料庫除錯
#   redis-tools            — 可在容器內直接用 redis-cli 指令連 Redis 除錯
RUN apt-get update && apt-get install -y --no-install-recommends \
    git curl ca-certificates unzip bash tini sudo \
    iproute2 iputils-ping net-tools \
    vim \
    bash-completion \
    default-mysql-client redis-tools \
  && rm -rf /var/lib/apt/lists/* \
  && echo 'alias ll="ls -l"' >> /etc/bash.bashrc \
  && echo 'alias la="ls -la"' >> /etc/bash.bashrc

# 設定容器時區為台北
RUN ln -sf /usr/share/zoneinfo/Asia/Taipei /etc/localtime \
  && echo "Asia/Taipei" > /etc/timezone

# 本機專案目錄會透過 bind mount 掛載進容器，這行讓 git 信任掛載進來的目錄
RUN git config --global --add safe.directory '*'

# 安裝常用的 Go 開發工具
#   gopls        — Go 語言伺服器（VS Code 智慧提示必要）
#   dlv          — Go 除錯器
#   golangci-lint — 程式碼檢查工具
#   air          — 熱重載工具（修改程式碼自動重新編譯執行）
RUN go install golang.org/x/tools/gopls@latest \
  && go install github.com/go-delve/delve/cmd/dlv@latest \
  && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
  && go install github.com/air-verse/air@latest

WORKDIR /workspaces

# tini 作為 PID 1，負責處理信號轉發與回收子程序
ENTRYPOINT ["/usr/bin/tini","--"]
CMD ["bash","-lc","sleep infinity"]
