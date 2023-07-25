package router

import (
	"net/http"
	"strings"

	"github.com/SergeyTyurin/banner_rotation/database"
	"github.com/SergeyTyurin/banner_rotation/handlers"
	"github.com/SergeyTyurin/banner_rotation/message_broker"
)

type Router interface {
	Mux() *http.ServeMux
}

type routerImpl struct {
	handlers handlers.Handlers
	mux      *http.ServeMux
}

func NewRouter(db database.Database, broker message_broker.MessageBroker) Router {
	var r routerImpl
	r.mux = http.NewServeMux()
	r.mux.HandleFunc("/banner", r.handleBannersFunc)
	r.mux.HandleFunc("/slot", r.handleSlotsFunc)
	r.mux.HandleFunc("/group", r.handleGroupsFunc)
	r.mux.HandleFunc("/rotation", r.handleRotationFunc)
	r.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimRight(r.URL.Path, "/") != "" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte("Rotation service is running"))
	})

	r.handlers = handlers.NewHandlers(db, broker)
	return &r
}

func (r *routerImpl) Mux() *http.ServeMux {
	return r.mux
}
