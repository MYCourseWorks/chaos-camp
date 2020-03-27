package settings

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// Settings comment
type Settings struct {
	ListOfURLs               []string `json:"ListOfURLs"`
	MaxRoutines              int      `json:"MaxRoutines"`
	MaxIndexingTimeInMs      int      `json:"MaxIndexingTimeInMs"`
	BaseResultFileName       string   `json:"BaseResultFileName"`
	MaxResultFileSizeInBytes int64    `json:"MaxResultFileSizeInBytes"`
	MaxResultFiles           int      `json:"MaxResultFiles"`
}

var instance *Settings

func init() {
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	instance, err = decodeSettings(bytes)
	if err != nil {
		panic(err)
	}
}

func decodeSettings(b []byte) (*Settings, error) {
	dec := json.NewDecoder(bytes.NewBuffer(b))
	s := new(Settings)
	dec.Decode(s)
	return s, nil
}

// GetSettings comment
func GetSettings() *Settings {
	return instance
}
