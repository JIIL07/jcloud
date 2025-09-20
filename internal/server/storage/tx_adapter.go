package storage

import (
	"github.com/JIIL07/jcloud/internal/server/types"
	"github.com/jmoiron/sqlx"
)

type TxAdapter struct {
	*sqlx.Tx
}

func NewTxAdapter(tx *sqlx.Tx) types.Tx {
	return &TxAdapter{Tx: tx}
}

func (ta *TxAdapter) Exec(query string, args ...interface{}) (interface{}, error) {
	return ta.Tx.Exec(query, args...)
}
