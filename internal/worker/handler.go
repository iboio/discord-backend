package worker

import (
	"discord/context"
	"discord/internal/batcher"
	"discord/query"
)

type Worker struct {
	stream  *context.Stream
	query   *query.Query
	batcher *batcher.Batcher
}

func NewWorker(appCtx context.AppContext) *Worker {
	return &Worker{
		stream:  appCtx.Stream(),
		query:   query.NewQuery(appCtx),
		batcher: batcher.NewBatcher(appCtx),
	}
}

func (w *Worker) StartWorkers() {
	go w.MessageProcessor()
	go w.VoiceProcessor()
	go w.UserProcessor()
}
