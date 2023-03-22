package postgres

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const modulesTable = "modules"

type modulesQ struct {
	db            *pgdb.DB
	selectBuilder sq.SelectBuilder
}

func NewModuleQ(db *pgdb.DB) data.ModuleQ {
	return &modulesQ{
		db:            db,
		selectBuilder: sq.Select(modulesTable + ".*").From(modulesTable),
	}
}

func (m modulesQ) New() data.ModuleQ {
	return NewModuleQ(m.db)
}

func (m modulesQ) FilterByNames(names ...string) data.ModuleQ {
	stmt := sq.Eq{modulesTable + ".name": names}
	m.selectBuilder = m.selectBuilder.Where(stmt)
	return m
}

func (m modulesQ) FilterByIsModule(isModule bool) data.ModuleQ {
	stmt := sq.Eq{modulesTable + ".is_module": isModule}
	m.selectBuilder = m.selectBuilder.Where(stmt)
	return m
}

func (m modulesQ) Get() (*data.Module, error) {
	var result data.Module
	err := m.db.Get(&result, m.selectBuilder)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (m modulesQ) Select() ([]data.Module, error) {
	var result []data.Module
	err := m.db.Select(&result, m.selectBuilder)
	return result, errors.Wrap(err, "failed to select modules")
}

func (m modulesQ) Insert(module data.Module) error {
	insertStmt := sq.Insert(modulesTable).SetMap(structs.Map(module)).Suffix("ON CONFLICT (name) DO NOTHING")
	err := m.db.Exec(insertStmt)
	return errors.Wrap(err, "failed to insert module")
}

func (m modulesQ) Delete(name string) error {
	query := sq.Delete(modulesTable).Where(sq.Eq{"name": name})

	result, err := m.db.ExecWithResult(query)
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows == 0 {
		return errors.New("no module with such name")
	}

	return nil
}
