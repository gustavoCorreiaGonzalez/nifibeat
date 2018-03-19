// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	URL    string        `config:"url"`
	Method string        `config:"method"`
	Period time.Duration `config:"period"`
}

var DefaultConfig = Config{
	URL:    "",
	Method: "",
	Period: 1 * time.Second,
}
