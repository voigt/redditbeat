// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	Reddit struct {
		Subs      *[]string
		Username  *string `config:"username"`
		Password  *string `config:"password"`
		Useragent *string `config:"useragent"`
		Limit     *int    `config:"limit"`
	}
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}
