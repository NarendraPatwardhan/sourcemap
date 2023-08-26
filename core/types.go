package core

type Repository []*Commit

type Commit struct {
	Hash    string `json:"hash"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Data    *Data  `json:"data"`
}

type Data struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Size     int64   `json:"size"`
	Repr     string  `json:"repr,omitempty"`
	Children []*Data `json:"children,omitempty"`
}

type ParseOpts struct {
	Limit        int
	ExcludeGlobs []string
	ExcludePaths []string
}

type WalkFunc func(data *Data, args ...interface{}) error
