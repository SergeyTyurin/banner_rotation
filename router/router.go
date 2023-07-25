package router

import (
	"net/http"
	"strings"

	"github.com/SergeyTyurin/banner-rotation/database"
	"github.com/SergeyTyurin/banner-rotation/handlers"
	"github.com/SergeyTyurin/banner-rotation/messagebroker"
)

type Router interface {
	CustomMux() *http.ServeMux
}

type routerImpl struct {
	handlers handlers.Handlers
	mux      *http.ServeMux
}

func NewRouter(db database.Database, broker messagebroker.MessageBroker) Router {
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

func (r *routerImpl) CustomMux() *http.ServeMux { //nolint:stylecheck
	return r.mux
}
