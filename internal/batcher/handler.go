package batcher

import (
	"discord/context"
	"discord/query"
)

type Batcher struct {
	query          *query.Query
	MessageRecords []string
	VoiceRecords   []string
}

func NewBatcher(appCtx context.AppContext) *Batcher {
	return &Batcher{
		query:          query.NewQuery(appCtx),
		MessageRecords: make([]string, 0),
		VoiceRecords:   make([]string, 0),
	}
}
