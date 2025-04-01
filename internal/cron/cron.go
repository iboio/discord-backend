package cron

import (
	"discord/pkg/clickhouse"
)

type Handler struct {
	ch *clickhouse.ConnectionClickhouse
}

func HandlerProcessor(CH *clickhouse.ConnectionClickhouse) *Handler {
	return &Handler{ch: CH}
}

func (h *Handler) StartCron() {
	h.VoiceLogScheduler()
}
