package data

import "encoding/json"

type RequestStatus string

const (
	CREATED  RequestStatus = "created"
	PENDING  RequestStatus = "pending"
	FINISHED RequestStatus = "finished"
	FAILED   RequestStatus = "failed"
)

type Request struct {
	ID         string          `json:"id" db:"id" structs:"-"`
	FromUserID string          `json:"from_user_id" db:"from_user_id"`
	ToUserID   string          `json:"to_user_id" db:"to_user_id"`
	Payload    json.RawMessage `json:"payload" db:"payload"`
	Status     RequestStatus   `json:"status" db:"status"`
	Error      *string         `json:"error,omitempty" db:"error,omitempty"`
	CreatedAt  string          `json:"created_at"`
	Module     *Module         `json:"module,omitempty" structs:"module,omitempty"`
}

type RequestQ interface {
	New() RequestQ

	FilterByIDs(ids ...string) RequestQ
	FilterByStatuses(statuses ...RequestStatus) RequestQ
	JoinsModule() RequestQ

	Get() (*Request, error)
	Select() ([]Request, error)

	Insert(request Request) error
	Update(request Request) error

	SetStatus(status RequestStatus) error
}
