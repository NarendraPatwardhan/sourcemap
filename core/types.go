package core

type Repository []*Commit

type Commit struct {
	Hash      string `json:"hash"`
	Author    string `json:"author"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Data      *Data  `json:"data"`
}

type Data struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	size     int64
	Size     string  `json:"size"`
	Changes  Changes `json:"changes"`
	Repr     string  `json:"repr,omitempty"`
	Children []*Data `json:"children,omitempty"`
}

type Changes struct {
	Addition int `json:"addition"`
	Deletion int `json:"deletion"`
}

type ParseOpts struct {
	Limit        int
	ExcludeGlobs []string
	ExcludePaths []string
}
