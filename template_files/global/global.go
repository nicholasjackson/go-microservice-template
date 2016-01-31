package global

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigStruct struct {
	StatsDServerIP string `json:"stats_d_server"`
	RootFolder     string
}

var Config ConfigStruct

func LoadConfig(config string, rootfolder string) error {
	fmt.Println("Loading Config: ", config)

	file, err := os.Open(config)
	if err != nil {
		return fmt.Errorf("Unable to open config")
	}

	decoder := json.NewDecoder(file)
	Config = ConfigStruct{}
	err = decoder.Decode(&Config)
	Config.RootFolder = rootfolder

	fmt.Println(Config)

	return nil
}
