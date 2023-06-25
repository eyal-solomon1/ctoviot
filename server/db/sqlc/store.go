package db

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/net/context"
)

type Store interface {
	VideoTx(ctx context.Context, arg VideoTxParam) (VideoTxResult, error)
	Querier
}

// SQLStore provides all functions to execute SQL queries or transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{db: db, Queries: New(db)}
}

// execTx executes a function within a db transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbError := tx.Rollback(); rbError != nil {
			return fmt.Errorf("tx error: %v,rb error: %v", err, rbError)

		}
		return err
	}
	return tx.Commit()
}

// VideoTxParam contains the input param for a new video transaction
type VideoTxParam struct {
	Username string            `json:"username"`
	Video    CreateVideoParams `json:"video_info"`
}

// VideoTxResult is the result of the new video creation
type VideoTxResult struct {
	Video Video `json:"video"`
	User  User  `json:"user"`
	Entry Entry `json:"entry"`
}

// VideoTx attaches a new video to the provided username
// It creates a new video recored, add user entries, and updates the users's balanc within a single db transaction
func (store *SQLStore) VideoTx(ctx context.Context, arg VideoTxParam) (VideoTxResult, error) {
	var result VideoTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		user, err := q.GetUser(ctx, arg.Username)

		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("user not found")
			}
			return err
		}

		if user.Balance-arg.Video.VideoLength <= 0 {
			return errors.New("user current balance isn't enough")
		}

		entriesCount, err := q.GetUsersVideosCount(ctx, arg.Username)

		if err != nil {
			return err
		}

		if entriesCount >= 3 {
			return errors.New("users reached videos limit (3)")
		}

		if err != nil {
			return err
		}
		result.Video, err = q.CreateVideo(ctx, CreateVideoParams{
			Owner:           arg.Username,
			VideoLength:     arg.Video.VideoLength,
			VideoIdentifier: arg.Video.VideoIdentifier,
			VideoName:       arg.Video.VideoName,
			VideoRemotePath: arg.Video.VideoRemotePath,
			VideoDecs:       arg.Video.VideoDecs,
		})

		if err != nil {
			return err
		}

		result.Entry, err = q.CreateEntry(ctx, CreateEntryParams{
			Username:  arg.Username,
			VideoName: arg.Video.VideoName,
			Amount:    -arg.Video.VideoLength,
		})

		if err != nil {
			return err
		}

		result.User, err = q.UpdateUserBalance(ctx, UpdateUserBalanceParams{
			Balance:  arg.Video.VideoLength,
			Username: arg.Username,
		})

		if err != nil {
			return err
		}

		return nil

	})

	return result, err
}
