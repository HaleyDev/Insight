package autoload

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	PrintSql     bool   `mapstructure:"print_sql"`
	LogLevel     string `mapstructure:"log_level"`
	TablePrefix  string `mapstructure:"table_prefix"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxLifetime  int    `mapstructure:"max_life_time"`
	Enable       bool   `mapstructure:"enable"`
}
