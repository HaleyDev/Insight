package autoload

type LoggerConfig struct {
	FileName        string       `mapstructure:"file_name"`
	DefaultDivision string       `mapstructure:"default_division"`
	DivisionTime    DivisionTime `mapstructure:"division_time"`
	DivisionSize    DivisionSize `mapstructure:"division_size"`
}

type DivisionTime struct {
	MaxAge       int64 `mapstructure:"max_age"`
	RotationTime int64 `mapstructure:"rotation_time"`
}

type DivisionSize struct {
	MaxSize    int  `mapstructure:"max_size"`
	MaxBackups int  `mapstructure:"max_backups"`
	MaxAge     int  `mapstructure:"max_age"`
	Compress   bool `mapstructure:"compress"`
}
