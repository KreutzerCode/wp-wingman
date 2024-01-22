package types

type PluginData struct {
	Name    string
	Version string
	Found   bool
}

type PluginInfo struct {
	Info struct {
		Pages int `json:"pages"`
	} `json:"info"`
	Plugins []struct {
		Slug string `json:"slug"`
	} `json:"plugins"`
}

type VersionNumber struct {
	Number  string
	FoundBy string
}