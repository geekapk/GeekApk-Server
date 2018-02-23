package modelmap

import (
	"github.com/gorilla/sessions"
)

type RequestContext struct {
	Session *sessions.Session
}
