package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rancher/apiserver/pkg/server"
	"github.com/rancher/apiserver/pkg/store/apiroot"
	"github.com/rancher/apiserver/pkg/types"

	// pprof
	"net/http/pprof"
)

// Start starts daemon.
func Start() http.Handler {
	s := server.DefaultAPIServer()

	// register schemas.
	initGeoIP(s.Schemas)

	apiroot.Register(s.Schemas, []string{"v1"})
	router := mux.NewRouter()
	router.UseEncodedPath()
	router.StrictSlash(true)

	router.Path("/").HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, "/", http.StatusFound)
	})

	// profiling handlers for pprof under /debug/pprof.
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)

	// manually add support for paths linked to by index page at /debug/pprof/.
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

	router.Path("/{prefix}/{type}").Handler(s)
	router.Path("/{prefix}/{type}/{name}").Queries("link", "{link}").Handler(s)
	router.Path("/{prefix}/{type}/{name}").Queries("action", "{action}").Handler(s)
	router.Path("/{prefix}/{type}/{name}").Handler(s)

	router.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		s.Handle(&types.APIRequest{
			Request:   r,
			Response:  rw,
			Type:      "apiRoot",
			URLPrefix: "v1",
		})
	})

	return router
}
