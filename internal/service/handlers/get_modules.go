package handlers

import (
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/ape"
	"net/http"
)

func GetModules(w http.ResponseWriter, r *http.Request) {
	ape.Render(w, newGetModulesResponse())
}

func newGetModulesResponse() resources.ModuleIconListResponse {
	var moduleIcon = map[string]string{
		"gitlab": "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/GitLab_logo.svg",
		//"github":          "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/GitHub_logo.svg",
		//"telegram":        "https://gist.githack.com/mhrynenko/bc7966e91e28ac17dbe809503855e60a/raw/fdbaf990e1f1a1fdc84712af801413b34a709959/Telegram_logo.svg",
		//"docker":          "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/DockerHub_logo.svg",
		//"click_up":        "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/ClickUp_logo.svg",
		//"slack":           "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/Slack_logo.svg",
		//"distributed_lab": "https://gist.githack.com/mhrynenko/422743c0c88341f79542a4e381d239f5/raw/13b466475d9096da01f3958f956457b6a4a30498/DL_logo.svg",
	}

	var result []resources.ModuleIcon
	for module, icon := range moduleIcon {
		result = append(result, resources.ModuleIcon{
			Key: resources.Key{
				ID:   module,
				Type: resources.MODULES,
			},
			Attributes: resources.ModuleIconAttributes{
				Icon: icon,
				Name: module,
			},
		})
	}
	return resources.ModuleIconListResponse{
		Data: result,
	}
}
