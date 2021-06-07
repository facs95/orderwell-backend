package router

import "net/http"

type Router interface {
	Get(uri string, handlerFn http.HandlerFunc)
	PrivateGet(uri string, authMiddleware func(http.Handler) http.Handler, handlerFn http.HandlerFunc)
	Post(uri string, handlerFn http.HandlerFunc)
	Serve(port string) error

	//Helpers
	GetUrlParam(r *http.Request, key string) string
	UseCors()
}
