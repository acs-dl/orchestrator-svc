package data

type Module struct {
	Endpoint *string `json:"endpoint" structs:"endpoint"`
	Link     *string `json:"link" structs:"link"`
	Name     string  `json:"name" structs:"name"`
}

type ModuleQ interface {
	New() ModuleQ

	FilterByNames(names ...string) ModuleQ

	Get() (*Module, error)
	Select() ([]Module, error)

	Insert(module Module) error
	Delete(name string) error
}
