package manager

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	FilePath string `json:"-"`
	Groups   []HostsGroup
}

func createDefaultConfigFilePath() (string, error) {
	var defaultConfigDirName = "hosts-manager"
	var defaultConfigFileName = "hosts-manager.json"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(homeDir, ".config", defaultConfigDirName)
	if err := os.MkdirAll(path, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(path, defaultConfigFileName), nil
}

func InitializeConfig(path string) (*Config, error) {

	var configFilePath string = path
	var err error
	if path == "" {
		configFilePath, err = createDefaultConfigFilePath()
		if err != nil {
			return nil, err
		}

		if _, err := os.Stat(configFilePath); err != nil {
			config := NewConfig()
			config.FilePath = configFilePath
			return config, nil
		}
	}

	config, err := LoadConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func LoadConfig(filepath string) (*Config, error) {

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	config.FilePath = filepath
	return &config, nil
}

func SaveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if config.FilePath == "" {
		configFilePath, err := createDefaultConfigFilePath()
		if err != nil {
			return err
		}
		config.FilePath = configFilePath
	}

	return ioutil.WriteFile(config.FilePath, data, 0644)
}

func NewConfig() *Config {
	return &Config{
		Groups: []HostsGroup{},
	}
}

func (config *Config) ExistsHostsGroup(name string) bool {
	for _, group := range config.Groups {
		if group.Name == name {
			return true
		}
	}
	return false
}

func (config *Config) GetHostsGroup(name string) *HostsGroup {
	for _, group := range config.Groups {
		if group.Name == name {
			var ret HostsGroup = group
			return &ret
		}
	}
	return nil
}

func (config *Config) SetHostsGroup(group HostsGroup) {
	for i := 0; i < len(config.Groups); i++ {
		if config.Groups[i].Name == group.Name {
			config.Groups[i] = group
			return
		}
	}
	config.Groups = append(config.Groups, group)
}

func (config *Config) RemoveHostsGroup(name string) {
	for i := 0; i < len(config.Groups); i++ {
		if config.Groups[i].Name == name {
			if i == len(config.Groups)-1 {
				config.Groups = config.Groups[:i]
			} else {
				config.Groups = append(config.Groups[:i], config.Groups[i+1:]...)
			}
			return
		}
	}
}

type HostsGroup struct {
	Name     string
	Body     string
	IsActive bool
}
