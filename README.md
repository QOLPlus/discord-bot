Discord Bot
===========

## install dependencies

### Private repos:
```bash
git config --global url.ssh://git@github.com:.insteadOf https://github.com/
GOPRIVATE=github.com/QOLPlus go get github.com/QOLPlus/core
```

### Public repos:
Public packages would be installed when go build.

## build & run
```bash
go build
./discord-bot
```