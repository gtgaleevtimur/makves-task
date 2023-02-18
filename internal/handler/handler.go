package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"makves-task/internal/repository"
)

const (
	ContentTypeApplicationJSON = "application/json"
)

func NewRouter(d repository.Databaser) chi.Router {
	router := chi.NewRouter()

	controller := newController(d)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/get-items", controller.GetItems)

	router.NotFound(NotFound())
	router.MethodNotAllowed(NotAllowed())

	return router
}

type Controller struct {
	Database repository.Databaser
}

// newController - функция-конструктор контролера хэндлера.
func newController(s repository.Databaser) *Controller {
	return &Controller{Database: s}
}

func (c *Controller) GetItems(w http.ResponseWriter, r *http.Request) {
	ids := sliceOfIDs(r.URL.RawQuery)
	result := c.Database.GetItems(ids)
	if len(result) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := json.Marshal(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", ContentTypeApplicationJSON)
	w.Write(body)
}

// NotFound - обработчик неподдерживаемых маршрутов.
func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("route does not exist"))
	}
}

// NotAllowed - обработчик неподдерживаемых методов.
func NotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("method does not allowed"))
	}
}

func sliceOfIDs(value string) []string {
	re := regexp.MustCompile("[0-9]+")
	result := re.FindAllString(value, -1)
	return result
}
