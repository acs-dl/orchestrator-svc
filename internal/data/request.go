package data

import (
	"encoding/json"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type RequestStatus string

const (
	CREATED  RequestStatus = "pending"
	PENDING  RequestStatus = "in progress"
	FINISHED RequestStatus = "success"
	INVITED  RequestStatus = "invited"
	FAILED   RequestStatus = "failed"
)

type Request struct {
	ID          string          `json:"id" db:"id" structs:"id"`
	FromUserID  int64           `json:"from_user_id" db:"from_user_id" structs:"from_user_id,omitempty"`
	ToUserID    int64           `json:"to_user_id" db:"to_user_id" structs:"to_user_id,omitempty"`
	Payload     json.RawMessage `json:"payload" db:"payload" structs:"payload"`
	Status      RequestStatus   `json:"status" db:"status" structs:"status,omitempty"`
	Error       *string         `json:"error,omitempty" db:"error,omitempty" structs:"error,omitempty"`
	Description *string         `json:"description,omitempty" db:"description,omitempty" structs:"description,omitempty"`
	Group       *[]string       `json:"group,omitempty" db:"group,omitempty" structs:"group,omitempty"`
	CreatedAt   string          `json:"created_at" db:"created_at" structs:"created_at,omitempty"`
	Module      *Module         `json:"module,omitempty" structs:"module,omitempty"`
	ModuleName  string          `json:"module_name" db:"module_name" structs:"module_name,omitempty"`
}

type RequestQ interface {
	New() RequestQ

	FilterByIDs(ids ...string) RequestQ
	FilterByStatuses(statuses ...RequestStatus) RequestQ
	FilterByFromIds(ids ...int64) RequestQ
	FilterByToIds(ids ...int64) RequestQ
	FilterNotByActions(actions ...string) RequestQ
	JoinsModule() RequestQ

	Get() (*Request, error)
	Select() ([]Request, error)

	Insert(request Request) error
	Update(request Request) error

	SetStatus(status RequestStatus) error
	SetStatusError(status RequestStatus, errorMsg string) error

	Count() RequestQ
	GetTotalCount() (int64, error)

	Page(pageParams pgdb.OffsetPageParams) RequestQ
}
