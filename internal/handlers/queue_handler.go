package handlers

import (
	"cache/pkg/logger"
	"github.com/go-chi/render"
	"net/http"
)

type QueueService interface {
	SendMessage(message string) error
}

type queueBrokerHandler struct {
	service QueueService
	logger  logger.Logger
}

func NewBrokerHandler(service QueueService, logger logger.Logger) QueueBrokerHandler {
	return &queueBrokerHandler{
		service: service,
		logger:  logger,
	}
}

func (h queueBrokerHandler) SaveMessage(w http.ResponseWriter, r *http.Request) {
	var req MessageRequest
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		h.logger.Errorf("%s", err)
		render.JSON(w, r, NewRespError(err, http.StatusBadRequest))
		return
	}

	err = h.service.SendMessage(req.Message)
	if err != nil {
		h.logger.Errorf("%s", err)
		render.JSON(w, r, NewRespError(err, http.StatusBadRequest))
		return
	}

	render.JSON(w, r, err)
}
