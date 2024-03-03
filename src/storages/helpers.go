package storages

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func scanUniqReactions(ctx context.Context, rows pgx.Rows) (map[string]struct{}, error) {
	defer rows.Close()
	res := make(map[string]struct{})
	for rows.Next() {
		var reactionId string
		err := rows.Scan(&reactionId)
		if err != nil {
			return nil, err
		}
		res[reactionId] = struct{}{}
	}
	return res, nil
}

func scanUniqReactionsStrict(ctx context.Context, rows pgx.Rows) map[string]struct{} {
	res, err := scanUniqReactions(ctx, rows)
	if err != nil {
		log.Fatalf("failed to scan uniq reactions: %w", err)
	}
	return res
}

func stopOnError(fns ...func() error) error {
	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}
	return nil
}
