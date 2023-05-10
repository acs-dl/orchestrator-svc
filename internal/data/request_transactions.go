package data

type TransactionAction string

const (
	DeleteUser TransactionAction = "delete_user"
)

type RequestTransactions interface {
	New() RequestTransactions

	FilterByIDs(ids ...string) RequestTransactions

	Get() (*RequestTransaction, error)
	Select() ([]RequestTransaction, error)
	Insert(request RequestTransaction) error
	Update(request RequestTransaction) error
}

type RequestTransaction struct {
	ID       string            `json:"id" db:"id" structs:"id"`
	Action   TransactionAction `json:"action" db:"action" structs:"action"`
	Requests map[string]bool   `json:"requests" db:"requests" structs:"requests"`
}
