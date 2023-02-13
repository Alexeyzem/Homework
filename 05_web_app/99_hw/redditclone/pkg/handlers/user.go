package handlers

import (
	"gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/session"
	"gitlab.com/mailru-go/lectures-2022-1/05_web_app/99_hw/redditclone/pkg/user"
	"go.uber.org/zap"
	"html/template"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo user.UserRepo
	Sessions *session.SessionsManager
}
