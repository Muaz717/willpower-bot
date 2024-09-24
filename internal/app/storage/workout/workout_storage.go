package workoutStorage

import (
	"context"
	"fmt"

	"github.com/Muaz717/willpower-bot/internal/lib/gym"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	usersTable         = "users"
	workoutsTable      = "workouts"
	usersWorkoutsTable = "users_workouts"
)

type WorkoutStorage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *WorkoutStorage {
	return &WorkoutStorage{
		db: db,
	}
}

func (s *WorkoutStorage) SaveUser(ctx context.Context, chatId int, userName string) (int, error) {
	const op = "storage.workout.SaveUser"

	query := fmt.Sprintf("INSERT INTO %s (ch_id, username) VALUES ($1, $2) RETURNING id", usersTable)

	var userId int

	err := s.db.QueryRow(ctx, query, chatId, userName).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (s *WorkoutStorage) Create(ctx context.Context, chatId int, workout gym.Workout) (int, error) {
	const op = "storage.workout.Create"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	query := fmt.Sprintf("INSERT INTO %s (workout_date, weight) VALUES ($1, $2) RETURNING id", workoutsTable)

	var workoutId int

	err = tx.QueryRow(ctx, query, workout.Date, workout.Weight).Scan(&workoutId)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("%s: failed to scan row(workoutId): %w", op, err)
	}

	usersWorkoutsQuery := fmt.Sprintf("INSERT INTO %s (chat_id, workout_id) VALUES ($1, $2)", usersWorkoutsTable)

	_, err = tx.Exec(ctx, usersWorkoutsQuery, chatId, workoutId)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return workoutId, tx.Commit(ctx)
}

func (s *WorkoutStorage) GetAll(ctx context.Context, chatId int) ([]gym.Workout, error) {
	const op = "storage.workout.GetAll"

	query := fmt.Sprintf(`SELECT ROW_NUMBER() OVER (ORDER BY wt.id), wt.workout_date, wt.weight, wt.id FROM %s wt 
						INNER JOIN %s uw ON wt.id = uw.workout_id WHERE uw.chat_id = $1`,
		workoutsTable, usersWorkoutsTable)

	rows, err := s.db.Query(ctx, query, chatId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[gym.Workout])
}

func (s *WorkoutStorage) Delete(ctx context.Context, chatId, workoutId int) error {
	const op = "storage.workout.Delete"

	query := fmt.Sprintf(`DELETE FROM %s wt USING %s uw WHERE wt.id = uw.workout_id 
						AND uw.chat_id=$1 AND uw.workout_id=$2`,
		workoutsTable, usersWorkoutsTable)

	_, err := s.db.Exec(ctx, query, chatId, workoutId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
