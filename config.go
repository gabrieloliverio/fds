package fds

type Config struct {
	Flags map[string]bool
}

func NewConfig() Config {
	return Config{
		Flags: map[string]bool{"confirm": false, "insensitive": false, "literal": false, "verbose": false},
	}
}
