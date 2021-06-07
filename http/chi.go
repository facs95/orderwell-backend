package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

type chiRouter struct{}

var (
	chiDispatcher = chi.NewRouter()
)

func NewChiRouter() Router {
	return &chiRouter{}
}

func (*chiRouter) Get(uri string, handlerFn http.HandlerFunc) {
	chiDispatcher.Get(uri, handlerFn)
}

func (*chiRouter) Post(uri string, handlerFn http.HandlerFunc) {
	chiDispatcher.Post(uri, handlerFn)
}

func (*chiRouter) Serve(port string) error {
	fmt.Printf("Chi HTTP server running on port %v", port)
	// handler := cors.Default().Handler(chiDispatcher)
	return http.ListenAndServe(port, chiDispatcher)
}

func (*chiRouter) GetUrlParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (*chiRouter) UseCors() {
	chiDispatcher.Use(cors.AllowAll().Handler)
}

func (*chiRouter) PrivateGet(uri string, authMiddleware func(http.Handler) http.Handler, handlerFn http.HandlerFunc) {
	chiDispatcher.With(authMiddleware).Get(uri, handlerFn)
}
