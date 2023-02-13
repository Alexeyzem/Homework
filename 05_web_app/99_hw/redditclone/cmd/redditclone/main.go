package main

import (
	"gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/posts"
	"gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/session"
	"gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/user"
	"go.uber.org/zap"
	"html/template"
)

func main() {

	templates := template.Must(template.ParseGlob("./templates/*"))
	sm := session.NewSessionsManager()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	userRepo := user.NewMemoryRepo()
	postRepo := posts.NewMemoryRepo()
	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}
}
