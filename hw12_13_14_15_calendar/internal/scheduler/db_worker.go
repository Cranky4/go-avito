package internalscheduler

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBWorker struct {
	DB   *sql.DB
	Logg Logger
	Conf DBWorkerConf
}

func NewDBWorker(db *sql.DB, conf DBWorkerConf, logg Logger) *DBWorker {
	return &DBWorker{
		DB:   db,
		Conf: conf,
		Logg: logg,
	}
}

func (w *DBWorker) Run(ctx context.Context, ch chan Notification) {
	w.Logg.Info("ready to work...")

	go func(ctx context.Context, ch chan Notification) {
		w.runScanner(ctx, ch)
	}(ctx, ch)

	go func(ctx context.Context) {
		w.runCleaner(ctx)
	}(ctx)

	<-ctx.Done()

	w.Logg.Info("db worker stopped...")
}

func (w *DBWorker) runScanner(ctx context.Context, ch chan Notification) {
	selectStmt, err := w.DB.Prepare(
		"SELECT id, title, starts_at FROM events WHERE notify_after < $1 AND notified_at IS NULL LIMIT 10;",
	)
	if err != nil {
		w.Logg.Error(err.Error())
		return
	}
	defer selectStmt.Close()

	updateStmt, err := w.DB.Prepare(
		"UPDATE events SET notified_at = NOW() WHERE id = $1",
	)
	if err != nil {
		w.Logg.Error(err.Error())
		return
	}
	defer updateStmt.Close()

	scanPeriod, err := time.ParseDuration(w.Conf.ScanPeriod)
	if err != nil {
		w.Logg.Error(err.Error())
		return
	}

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case <-time.After(scanPeriod):
			w.Logg.Debug("start scanning new notifications...")

			rows, err := selectStmt.QueryContext(ctx, time.Now())
			if err != nil {
				w.Logg.Error(err.Error())
				break L
			}
			if rows.Err() != nil {
				w.Logg.Error(err.Error())
				break L
			}

			var id, title string
			var startsAt time.Time

			for rows.Next() {
				err = rows.Scan(&id, &title, &startsAt)
				if err != nil {
					w.Logg.Error(err.Error())
				}
				defer rows.Close()

				if id != "" {
					w.Logg.Info(fmt.Sprintf("Event '%s' will starts at '%s'", title, startsAt.Format(time.RFC822)))

					select {
					case <-ctx.Done():
						break L
					case ch <- Notification{ID: id, Title: title, StartedAt: startsAt}:
						_, err := updateStmt.ExecContext(ctx, id)
						if err != nil {
							w.Logg.Error(err.Error())
						}

						w.Logg.Debug(fmt.Sprintf("event updated %s", id))
					}
				}
			}
		}
	}
}

func (w *DBWorker) runCleaner(ctx context.Context) {
	stmt, err := w.DB.Prepare(
		"DELETE FROM events WHERE ends_at < $1",
	)
	if err != nil {
		w.Logg.Error(err.Error())
		return
	}
	defer stmt.Close()

	scanPeriod, err := time.ParseDuration(w.Conf.ScanPeriod)
	if err != nil {
		w.Logg.Error(err.Error())
		return
	}

	threshold := time.Unix(time.Now().Unix()-int64(w.Conf.ClearPeriodDays)*3600*24, 0)

L:
	for {
		select {
		case <-ctx.Done():
			break L
		case <-time.After(scanPeriod):
			w.Logg.Debug(fmt.Sprintf("start clearing notifications older than %s...", threshold.Format(time.RFC1123Z)))

			res, err := stmt.ExecContext(ctx, threshold)
			if err != nil {
				w.Logg.Error(err.Error())
				break L
			}
			rows, err := res.RowsAffected()
			if err != nil {
				w.Logg.Error(err.Error())
				break L
			}

			w.Logg.Info(fmt.Sprintf("cleared %d", rows))
		}
	}
}
