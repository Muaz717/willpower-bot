package pullupsStorage

import (
	"context"
	"fmt"

	"github.com/Muaz717/willpower-bot/internal/lib/gym"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pullUpsTable      = "pull_ups"
	usersTable        = "users"
	usersPullUpsTable = "users_pullups"
)

type PullUpsStorage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PullUpsStorage {
	return &PullUpsStorage{
		db: db,
	}
}

func (s *PullUpsStorage) Create(ctx context.Context, chatId int, pullUps gym.PullUps) (int, error) {
	const op = "storage.pull-ups.Create"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	query := fmt.Sprintf("INSERT INTO %s (quantity, date) VALUES ($1, $2) RETURNING id", pullUpsTable)

	var pullUpId int

	err = tx.QueryRow(ctx, query, pullUps.Quantity, pullUps.Date).Scan(&pullUpId)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	usersPullUpsQuery := fmt.Sprintf("INSERT INTO %s (chat_id, pull_up_id) VALUES ($1, $2)", usersPullUpsTable)

	_, err = tx.Exec(ctx, usersPullUpsQuery, chatId, pullUpId)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("%s: failed to scan row: %w", op, err)
	}

	return pullUpId, tx.Commit(ctx)
}

func (s *PullUpsStorage) GetAll(ctx context.Context, chatId int) ([]gym.PullUps, error) {
	const op = "storage.pull-ups.GetAll"

	query := fmt.Sprintf(`SELECT ROW_NUMBER() OVER (ORDER BY pt.id), pt.quantity, pt.date, pt.id FROM %s pt 
						INNER JOIN %s up ON pt.id = up.pull_up_id WHERE up.chat_id = $1`,
		pullUpsTable, usersPullUpsTable)

	rows, err := s.db.Query(ctx, query, chatId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[gym.PullUps])
}

func (s *PullUpsStorage) Delete(ctx context.Context, chatId, pullUpsId int) error {
	const op = "storage.pull-ups.Delete"

	query := fmt.Sprintf(`DELETE FROM %s pt USING %s up WHERE pt.id = up.pull_up_id
							AND up.chat_id = $1 AND up.pull_up_id = $2`,
		pullUpsTable, usersPullUpsTable)

	_, err := s.db.Exec(ctx, query, chatId, pullUpsId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
