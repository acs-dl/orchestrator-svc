package data

type Module struct {
	Endpoint *string `json:"endpoint"`
	Name     string  `json:"name"`
}

type ModuleQ interface {
	New() ModuleQ

	FilterByNames(names ...string) ModuleQ

	Get() (*Module, error)
	Select() ([]Module, error)

	Insert(module Module) error
}
