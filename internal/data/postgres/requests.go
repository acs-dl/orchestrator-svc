package postgres

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const requestsTable = "requests"

type requestsQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
	updateBuilder sq.UpdateBuilder
}

func NewRequestsQ(db *pgdb.DB) data.RequestQ {
	return &requestsQ{
		db:            db,
		selectBuilder: sq.Select(requestsTable + ".*").From(requestsTable),
		updateBuilder: sq.Update(requestsTable),
	}
}

func (r requestsQ) New() data.RequestQ {
	return NewRequestsQ(r.db)
}

func (r requestsQ) FilterByIDs(ids ...string) data.RequestQ {
	stmt := sq.Eq{requestsTable + ".id": ids}
	r.selectBuilder = r.selectBuilder.Where(stmt)
	r.updateBuilder = r.updateBuilder.Where(stmt)
	return r
}

func (r requestsQ) FilterByStatuses(statuses ...data.RequestStatus) data.RequestQ {
	stmt := sq.Eq{requestsTable + ".status": statuses}
	r.selectBuilder = r.selectBuilder.Where(stmt)
	r.updateBuilder = r.updateBuilder.Where(stmt)
	return r
}

func (r requestsQ) JoinsModule() data.RequestQ {
	r.selectBuilder = r.selectBuilder.
		LeftJoin(fmt.Sprint(modulesTable, " ON ", modulesTable, ".name = ", requestsTable, ".module_name"))
	return r
}

func (r requestsQ) Get() (*data.Request, error) {
	var result data.Request
	err := r.db.Get(&result, r.selectBuilder)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (r requestsQ) Select() ([]data.Request, error) {
	var result []data.Request
	err := r.db.Select(&result, r.selectBuilder)
	return result, errors.Wrap(err, "failed to select requests")
}

func (r requestsQ) Insert(request data.Request) error {
	insertStmt := sq.Insert(requestsTable).SetMap(structs.Map(request))
	err := r.db.Exec(insertStmt)
	return errors.Wrap(err, "failed to insert request")
}

func (r requestsQ) Update(request data.Request) error {
	updateStmt := r.updateBuilder.SetMap(structs.Map(request))
	err := r.db.Exec(updateStmt)
	return errors.Wrap(err, "failed to update request")
}

func (r requestsQ) SetStatus(status data.RequestStatus) error {
	updateStmt := r.updateBuilder.Set("status", status)
	err := r.db.Exec(updateStmt)
	return errors.Wrap(err, "failed to set status")
}

func (r requestsQ) SetStatusError(status data.RequestStatus, errorMsg string) error {
	updateStmt := r.updateBuilder.Set("status", status).Set("error", errorMsg)
	err := r.db.Exec(updateStmt)
	return errors.Wrap(err, "failed to set status and error")
}

func (r requestsQ) FilterByFromIds(ids ...int64) data.RequestQ {
	stmt := sq.Eq{requestsTable + ".from_user_id": ids}
	r.selectBuilder = r.selectBuilder.Where(stmt)
	r.updateBuilder = r.updateBuilder.Where(stmt)
	return r
}

func (r requestsQ) FilterByToIds(ids ...int64) data.RequestQ {
	stmt := sq.Eq{requestsTable + ".to_user_id": ids}
	r.selectBuilder = r.selectBuilder.Where(stmt)
	r.updateBuilder = r.updateBuilder.Where(stmt)
	return r
}

func (r requestsQ) FilterByActions(actions ...string) data.RequestQ {
	stmt := sq.Eq{requestsTable + ".payload->>'action'": actions}
	r.selectBuilder = r.selectBuilder.Where(stmt)
	r.updateBuilder = r.updateBuilder.Where(stmt)
	return r
}

func (r requestsQ) Count() data.RequestQ {
	r.selectBuilder = sq.Select("COUNT (*)").From(requestsTable)

	return r
}

func (r requestsQ) GetTotalCount() (int64, error) {
	var count int64

	err := r.db.Get(&count, r.selectBuilder)

	return count, err
}

func (r requestsQ) Page(pageParams pgdb.OffsetPageParams) data.RequestQ {
	r.selectBuilder = pageParams.ApplyTo(r.selectBuilder, "id")

	return r
}
