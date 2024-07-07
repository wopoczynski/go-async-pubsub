package server

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/wopoczynski/playground/internal/queue"
)

type Handler struct {
	ampq      queue.AMQP
	queueName string
}

func NewHandler(q queue.AMQP, qN string) *Handler {
	return &Handler{
		ampq:      q,
		queueName: qN,
	}
}

// Ping godoc
//
//	@Tags       Health
//	@Produce    plain
//	@Success    200 string  pong
//	@Router     /ping [get]
func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

// Message godoc
//
// @Tags	App
// @Produce	json
// @Param	request	body	database.PersistingStruct	true	"Request body"
// @Success	201	nil
// @Router	/messages	[post]
func (h *Handler) message(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid payload")
	}
	h.ampq.Publish(&queue.Message{
		Body:      string(body),
		QueueName: h.queueName,
	})

	return c.JSON(http.StatusCreated, nil)
}
