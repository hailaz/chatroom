package main

import (
	_ "chatroom/internal/packed"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"github.com/gogf/gf/v2/os/gctx"

	"chatroom/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
