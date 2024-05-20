package postgres

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
)

type DBLogger struct {
}

func (d DBLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d DBLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

func New(opts *pg.Options) *pg.DB {
	return pg.Connect(opts)
}
