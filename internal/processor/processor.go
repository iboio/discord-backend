package processor

import (
	"discord/pkg/clickhouse"
	js "discord/pkg/jetstream"
)

type Handler struct {
	ch *clickhouse.ConnectionClickhouse
	js *js.ContextJetstream
}

func HandlerProcessor(JS *js.ContextJetstream, CH *clickhouse.ConnectionClickhouse) *Handler {
	return &Handler{js: JS, ch: CH}
}

func (h *Handler) Process() {
	go h.MessageProcessor()
	go h.VoiceProcessor()
}
