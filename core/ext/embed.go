package ext

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"
)

//go:embed colors.json
var colors []byte

var (
	syncOnce sync.Once
	Map      map[string]string
)

func Get(name string) string {
	syncOnce.Do(func() {
		Map = make(map[string]string)
		json.Unmarshal(colors, &Map)
	})
	frags := strings.Split(strings.ToLower(name), ".")
	ext := frags[len(frags)-1]

	return Map[ext]
}
