package data

import "encoding/json"

type TransactionAction string

const (
	Single     TransactionAction = "single"
	DeleteUser TransactionAction = "delete_user"
)

type RequestTransactions interface {
	New() RequestTransactions

	Get() (*RequestTransaction, error)
	Select() ([]RequestTransaction, error)
	Insert(request RequestTransaction) error
	Update(request RequestTransaction) error
	Delete() error

	FilterByIDs(ids ...string) RequestTransactions
	FilterByRequestID(id string) RequestTransactions
}

type RequestTransaction struct {
	ID       string            `json:"id" db:"id" structs:"id"`
	Action   TransactionAction `json:"action" db:"action" structs:"action"`
	Requests json.RawMessage   `json:"requests" db:"requests" structs:"requests"`
}
