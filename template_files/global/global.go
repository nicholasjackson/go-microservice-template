package global

import (
	"encoding/json"
	"os"
)

type ConfigStruct struct {
	Mykey      string
	RootFolder string
}

var Config ConfigStruct

func LoadConfig(config string, rootfolder string) {
	file, err := os.Open(config)
	if err != nil {
		panic("Unable to open config")
	}

	decoder := json.NewDecoder(file)
	Config = ConfigStruct{}
	err = decoder.Decode(&Config)
	Config.RootFolder = rootfolder
}
