package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/max-messenger/max-bot-example-todolist/pkg/http/middlewares"
)

// Controller api controller.
type Controller interface {
	Register(r chi.Router)
}
type Controllers []Controller

// 	@title          TODOLIST API
// 	@version        1.0
// 	@description    This is a TODOLIST API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func NewRouter(logger *zap.Logger, controllers Controllers) (http.Handler, error) {
	router := chi.NewRouter()

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	})

	router.Route("/", func(r chi.Router) {
		r.Use(
			middlewares.Metrics(),
			middlewares.AccessLog(logger),
			middlewares.Trace(),
		)
		for _, ctrl := range controllers {
			ctrl.Register(r)
		}
	})

	return router, nil
}
