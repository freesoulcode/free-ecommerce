package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load reads a YAML config file and overrides environment variables.
func Load(path string, out any) error {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}	
	if err := v.Unmarshal(out); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	
	return nil
}