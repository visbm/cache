package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

type Handler struct {
	articleHandler     ArticleHandler
	queueBrokerHandler QueueBrokerHandler
}

func NewHandler(articleHandler ArticleHandler, queueBrokerHandler QueueBrokerHandler) *Handler {
	return &Handler{
		articleHandler:     articleHandler,
		queueBrokerHandler: queueBrokerHandler,
	}
}

type ArticleHandler interface {
	SaveArticle(w http.ResponseWriter, r *http.Request)
	GetArticle(w http.ResponseWriter, r *http.Request)
	GetArticles(w http.ResponseWriter, r *http.Request)
}

type QueueBrokerHandler interface {
	SaveMessage(w http.ResponseWriter, r *http.Request)
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/{id}", h.articleHandler.GetArticle)
	r.Get("/", h.articleHandler.GetArticles)
	r.Post("/", h.articleHandler.SaveArticle)

	r.Post("/message", h.queueBrokerHandler.SaveMessage)

	return r
}
