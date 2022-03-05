package config

type Config struct {
	Workers         int64  `yaml:"workers"`
	GenesisFilePath string `yaml:"genesis_file_path,omitempty"`
	ParseGenesis    bool   `yaml:"parse_genesis"`
	FastSync        bool   `yaml:"fast_sync,omitempty"`
}

// NewParsingConfig allows to build a new Config instance
func NewParsingConfig(workers int64, parseGenesis bool, genesisFilePath string, fastSync bool) Config {
	return Config{
		Workers:         workers,
		ParseGenesis:    parseGenesis,
		GenesisFilePath: genesisFilePath,
		FastSync:        fastSync,
	}
}

// DefaultParsingConfig returns the default instance of Config
func DefaultParsingConfig() Config {
	return NewParsingConfig(1, true, "", false)
}
