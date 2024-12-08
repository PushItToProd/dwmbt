// Config is loaded from a JSON config file and passed to the daemon.
// The config location can be overridden with the DWMBT_CONFIG_FILE environment variable.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileEnvVar = "DWMBT_CONFIG_FILE"
const DefaultConfigPath = "/etc/dwmbt/config.json"

type Config struct {
	ServeAddr string
	AuthKey   string `json:",omitempty"` // TODO: implement auth
	Peers     []struct {
		Addr        string
		DisplayName string `json:",omitempty"`
		AuthKey     string `json:",omitempty"`
	} `json:",omitempty"`
}

func GetConfigPath() string {
	if path := os.Getenv(ConfigFileEnvVar); path != "" {
		return path
	}

	if homedir, _ := os.UserHomeDir(); homedir != "" {
		return filepath.Join(homedir, ".config", "dwmbt", "config.json") // ~/.config/dwmbt/config.json
	}

	return DefaultConfigPath
}

func LoadConfigFile(path string) (Config, error) {
	// read the file
	j, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	// unmarshal the json
	var c Config
	err = json.Unmarshal(j, &c)
	if err != nil {
		return Config{}, err
	}

	// TODO: validate the config

	return c, nil
}

func setConfigDefaults(c *Config) {
	if c.ServeAddr == "" {
		c.ServeAddr = "localhost:11111"
	}
}

func LoadConfig() (Config, error) {
	// TODO: return LoadConfigFile(GetConfigPath())
	c := Config{}
	setConfigDefaults(&c)
	return c, nil
}
