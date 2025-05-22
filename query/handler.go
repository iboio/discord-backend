package query

import (
	"discord/context"
)

type Query struct {
	db *context.DB
}

func NewQuery(appCtx context.AppContext) *Query {
	return &Query{
		db: appCtx.DB(),
	}
}
