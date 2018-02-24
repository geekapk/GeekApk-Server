package modelmap

import (
	"github.com/gorilla/sessions"
)

// A RequestContext includes the context info of the request, like
// the current Session.
type RequestContext struct {
	Session *sessions.Session
}
