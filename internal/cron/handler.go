package cron

import (
	"discord/context"
	"discord/internal/batcher"
)

type Cron struct {
	batch *batcher.Batcher
}

func NewCron(appCtx context.AppContext) *Cron {
	return &Cron{
		batch: batcher.NewBatcher(appCtx),
	}
}

func (cr *Cron) StartCron() {
	cr.VoiceLogScheduler()
}
