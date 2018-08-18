package repository

import (
	"fmt"

	"github.com/go-gorp/gorp"
)

type (
	Transaction interface {
		Insert(i interface{}) Transaction
		Update(i interface{}) Transaction
		Commit() error
	}

	transactionImpl struct {
		tx  *gorp.Transaction
		err error
	}
)

func (t *transactionImpl) Insert(i interface{}) Transaction {
	if t.err == nil {
		t.err = t.tx.Insert(i)
	}
	return t
}

func (t *transactionImpl) Update(i interface{}) Transaction {
	if t.err == nil {
		if cnt, err := t.tx.Update(i); err != nil {
			t.err = err
		} else if cnt <= 0 {
			t.err = fmt.Errorf("no target to update")
		}
	}
	return t
}

func (t *transactionImpl) Commit() error {
	// already has error, not commit
	if t.err != nil {
		return t.err
	}
	if err := t.tx.Commit(); err != nil {
		return err
	}
	return nil
}
