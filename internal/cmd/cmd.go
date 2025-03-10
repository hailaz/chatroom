package cmd

import (
	"context"
	"goframechat/internal/controller/chat"
	"goframechat/internal/controller/chatroom"
	"goframechat/internal/controller/user"
	"goframechat/internal/dao"
	"goframechat/internal/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// Initialize database
			if err := dao.InitDatabase(ctx); err != nil {
				return err
			}

			s := g.Server()

			// Global middleware
			s.Use(ghttp.MiddlewareHandlerResponse)

			// Register controllers
			s.Group("/", func(group *ghttp.RouterGroup) {
				// Public routes
				group.Group("/api", func(group *ghttp.RouterGroup) {
					// User routes
					userController := user.NewController()
					group.Bind(
						// Only register non-auth routes here
						userController.Register,
						userController.Login,
					)

					// Protected routes
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)

						// User profile routes
						group.Bind(
							userController.Profile,
							userController.UpdateProfile,
						)

						// Chat room routes
						roomController := chatroom.NewController()
						group.Bind(
							roomController,
						)

						// Chat routes
						chatController := chat.NewController()
						group.Bind(
							chatController.GetHistory,
							chatController.GetRoomMembers,
						)
					})
				})

				// WebSocket routes
				group.Group("/ws", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.Auth)
					chatController := chat.NewController()
					group.GET("/chat", chatController.Connect)
				})

				// Static files
				s.SetIndexFolder(true)
				s.SetServerRoot("resource/public")
			})

			s.Run()
			return nil
		},
	}
)
