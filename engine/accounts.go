package engine

import (
	"time"

	"github.com/pkg/errors"
)

//go:generate reform

//reform:acca.currencies
type Currency struct {
	CurrencyID int64   `reform:"curr_id,pk"`
	Key        string  `reform:"key"`
	Meta       *[]byte `reform:"meta"`
}

//reform:acca.accounts
type Account struct {
	AccountID  int64     `reform:"acc_id,pk"`
	CurrencyID int64     `reform:"curr_id"`
	Key        string    `reform:"key"`
	Balance    int64     `reform:"balance"`
	Meta       *[]byte   `reform:"meta"`
	UpdatedAt  time.Time `reform:"updated_at"`
	CreatedAt  time.Time `reform:"created_at"`
}

func (a *Account) BeforeInsert() error {
	a.UpdatedAt = time.Now()
	a.CreatedAt = time.Now()
	if a.Balance != 0 {
		return errors.New("new account can't be with not zero balance")
	}
	return nil
}

func (a *Account) BeforeUpdate() error {
	a.UpdatedAt = time.Now()
	return nil
}