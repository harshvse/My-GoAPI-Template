package mailer

import "embed"

const (
	FromName            = "GoAPITemplate"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandBox bool) (int,error)
}
