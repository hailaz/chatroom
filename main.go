package main

import (
	_ "goframechat/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"goframechat/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
