package json

import "github.com/json-iterator/go"

// 替代原生json，性能提升
var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	Unmarshal     = json.Unmarshal
)
