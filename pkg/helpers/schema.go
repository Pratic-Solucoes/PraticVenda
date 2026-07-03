package helpers

import (
	"context"
	"database/sql"
	"fmt"
)

func SetSchema(ctx context.Context, tx *sql.Tx) error {
	schema, ok := ctx.Value("schema").(string)
	if !ok || schema == "" {
		return fmt.Errorf("schema não encontrado no contexto do request")
	}

	query := fmt.Sprintf("SET LOCAL search_path TO %s", schema)
	_, err := tx.ExecContext(ctx, query)
	return err
}
