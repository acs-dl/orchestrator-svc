package data

type Module struct {
	Topic    string `json:"topic" structs:"topic" db:"topic"`
	Link     string `json:"link" structs:"link" db:"link"`
	Name     string `json:"name" structs:"name" db:"name"`
	Title    string `json:"title" structs:"title" db:"title"`
	Prefix   string `json:"prefix" structs:"prefix" db:"prefix"`
	IsModule bool   `json:"is_module" structs:"is_module" db:"is_module"`
}

type ModuleQ interface {
	New() ModuleQ

	FilterByNames(names ...string) ModuleQ
	FilterByIsModule(isModule bool) ModuleQ

	Get() (*Module, error)
	Select() ([]Module, error)

	Insert(module Module) error
	Delete(name string) error
}

var ModuleIcon = map[string]string{
	"gitlab":   "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/GitLab_logo.svg",
	"github":   "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/GitHub_logo.svg",
	"telegram": "https://gist.githack.com/mhrynenko/bc7966e91e28ac17dbe809503855e60a/raw/fdbaf990e1f1a1fdc84712af801413b34a709959/Telegram_logo.svg",
	"docker":   "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/DockerHub_logo.svg",
	"click_up": "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/ClickUp_logo.svg",
	"slack":    "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/Slack_logo.svg",
	"mail":     "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/17bcfd78386a8a438b3abc0739d31c306b0ba33a/DL_logo.svg",
}
