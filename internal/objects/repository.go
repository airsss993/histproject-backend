package objects

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Repository — доступ к данным объектов и типов событий.
type Repository interface {
	GetByID(ctx context.Context, id int) (*SingleObjectInfo, error)
	List(ctx context.Context, eventTypeIDs []int, dateFrom, dateTo string) ([]ObjectInfo, error)
	ListEventTypes(ctx context.Context) ([]EventTypeInfo, error)
}

type repository struct {
	db *sqlx.DB
}

// NewRepository создаёт реализацию Repository для PostgreSQL.
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id int) (*SingleObjectInfo, error) {
	var obj SingleObjectInfo
	query := `
		SELECT id, title, description, event_date, preview_image_url, site_url
		FROM objects WHERE id = $1`
	err := r.db.GetContext(ctx, &obj, query, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения информации об объекте %v с БД: %w", id, err)
	}
	return &obj, nil
}

func (r *repository) List(ctx context.Context, eventTypeIDs []int, dateFrom, dateTo string) ([]ObjectInfo, error) {
	query := `
		SELECT id, title, description, latitude, longitude, event_date, event_type_id, preview_image_url
		FROM objects`
	var args []interface{}
	argPos := 1

	if len(eventTypeIDs) > 0 {
		query += fmt.Sprintf(" WHERE event_type_id = ANY($%d)", argPos)
		args = append(args, pq.Array(eventTypeIDs))
		argPos++
	}
	if dateFrom != "" {
		if argPos > 1 {
			query += fmt.Sprintf(" AND event_date >= $%d", argPos)
		} else {
			query += fmt.Sprintf(" WHERE event_date >= $%d", argPos)
		}
		args = append(args, dateFrom)
		argPos++
	}
	if dateTo != "" {
		if argPos > 1 {
			query += fmt.Sprintf(" AND event_date <= $%d", argPos)
		} else {
			query += fmt.Sprintf(" WHERE event_date <= $%d", argPos)
		}
		args = append(args, dateTo)
		argPos++
	}

	var list []ObjectInfo
	err := r.db.SelectContext(ctx, &list, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения объектов с БД: %w", err)
	}
	if list == nil {
		return []ObjectInfo{}, nil
	}
	return list, nil
}

func (r *repository) ListEventTypes(ctx context.Context) ([]EventTypeInfo, error) {
	var list []EventTypeInfo
	query := `SELECT id, name, description FROM event_types`
	err := r.db.SelectContext(ctx, &list, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения типов событий с БД: %w", err)
	}
	return list, nil
}
