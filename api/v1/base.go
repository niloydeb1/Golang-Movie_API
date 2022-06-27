package v1

import "github.com/labstack/echo/v4"

// Router api/v1 base router
func Router(g *echo.Group) {
	UserRouter(g.Group("/users"))
	OauthRouter(g.Group("/oauth"))
	MovieRouter(g.Group("/movies"))
	ReviewRouter(g.Group("/reviews"))
	CommentRouter(g.Group("/comments"))
}
