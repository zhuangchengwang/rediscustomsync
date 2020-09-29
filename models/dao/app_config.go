package dao

type AppConfig struct {
	Title string `toml:"title"`
	Description string `toml:"description"`
	Interval    int `toml:"interval"`
	EnableTTL   bool `toml:"enablettl"`
	Source RedisNode `toml:"source"`
	Destination RedisNode `toml:"destination"`
	Keysfile string `toml:"keysfile"`
	Patterns []string `toml:"patterns"`
}

type RedisNode struct {
	Host string `toml:"host"`
	Port int `toml:"port"`
	Db   int `toml:"db"`
	Auth string `toml:"auth"`
}

