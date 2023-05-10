package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const requestTransactionsTable = "request_transactions"

type RequestTransactionsQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	updateBuilder sq.UpdateBuilder
	deleteBuilder sq.DeleteBuilder
}

func NewRequestTransactionsQ(db *pgdb.DB) data.RequestTransactions {
	return &RequestTransactionsQ{
		db:            db,
		selectBuilder: sq.Select(requestTransactionsTable + ".*").From(requestTransactionsTable),
		updateBuilder: sq.Update(requestTransactionsTable),
		deleteBuilder: sq.Delete(requestTransactionsTable),
	}
}

func (r RequestTransactionsQ) New() data.RequestTransactions {
	return NewRequestTransactionsQ(r.db)
}

func (r RequestTransactionsQ) FilterByIDs(ids ...string) data.RequestTransactions {
	stmt := sq.Eq{requestsTable + ".id": ids}
	r.selectBuilder = r.selectBuilder.Where(stmt)
	r.updateBuilder = r.updateBuilder.Where(stmt)
	r.deleteBuilder = r.deleteBuilder.Where(stmt)
	return r
}

func (r RequestTransactionsQ) Get() (*data.RequestTransaction, error) {
	var result data.RequestTransaction
	err := r.db.Get(&result, r.selectBuilder)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (r RequestTransactionsQ) Select() ([]data.RequestTransaction, error) {
	var result []data.RequestTransaction

	err := r.db.Select(&result, r.selectBuilder)

	return result, err
}

func (r RequestTransactionsQ) Insert(request data.RequestTransaction) error {
	insertStmt := sq.Insert(requestTransactionsTable).SetMap(structs.Map(request))

	return r.db.Exec(insertStmt)
}

func (r RequestTransactionsQ) Update(request data.RequestTransaction) error {
	updateStmt := r.updateBuilder.SetMap(structs.Map(request))

	return r.db.Exec(updateStmt)
}

func (r RequestTransactionsQ) Delete() error {
	var deleted []data.RequestTransaction

	err := r.db.Select(&deleted, r.deleteBuilder.Suffix("RETURNING *"))
	if err != nil {
		return err
	}

	if len(deleted) == 0 {
		return sql.ErrNoRows
	}

	return nil
}
