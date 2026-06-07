package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Encrypt  EncryptConfig  `mapstructure:"encrypt"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// DSN 返回 MySQL 连接字符串
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name)
}

type JWTConfig struct {
	Secret       string `mapstructure:"secret"`
	ExpireHours  int    `mapstructure:"expire_hours"`
	RefreshHours int    `mapstructure:"refresh_hours"`
}

type EncryptConfig struct {
	Key string `mapstructure:"key"`
}

type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Addr 返回 Redis 地址
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type CORSConfig struct {
	Origins []string `mapstructure:"origins"`
}

// LoadConfig 加载配置
// 优先级：环境变量 > config.yaml > 默认值
func LoadConfig() (*Config, error) {
	v := viper.New()

	// 设置默认值
	v.SetDefault("app.name", "campus_collab")
	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", 8080)
	v.SetDefault("app.host", "0.0.0.0")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("jwt.expire_hours", 24)
	v.SetDefault("jwt.refresh_hours", 168)
	v.SetDefault("redis.enabled", true)
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("log.level", "debug")
	v.SetDefault("log.format", "console")
	v.SetDefault("cors.origins", []string{"http://localhost:5173"})

	// 读取 YAML 配置文件（可选）
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	_ = v.ReadInConfig() // 忽略文件不存在的错误

	// 环境变量覆盖（前缀 + 下划线映射）
	v.SetEnvPrefix("")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 手动绑定关键环境变量（Viper 的 AutomaticEnv 对嵌套结构支持有限）
	bindEnvVars(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 校验必填项
	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func bindEnvVars(v *viper.Viper) {
	envBindings := map[string]string{
		"app.name":             "APP_NAME",
		"app.env":              "APP_ENV",
		"app.port":             "APP_PORT",
		"app.host":             "APP_HOST",
		"database.host":        "DB_HOST",
		"database.port":        "DB_PORT",
		"database.name":        "DB_NAME",
		"database.user":        "DB_USER",
		"database.password":    "DB_PASSWORD",
		"database.max_idle_conns": "DB_MAX_IDLE_CONNS",
		"database.max_open_conns": "DB_MAX_OPEN_CONNS",
		"jwt.secret":           "JWT_SECRET",
		"jwt.expire_hours":     "JWT_EXPIRE_HOURS",
		"jwt.refresh_hours":    "JWT_REFRESH_HOURS",
		"encrypt.key":          "ENCRYPT_KEY",
		"redis.enabled":        "REDIS_ENABLED",
		"redis.host":           "REDIS_HOST",
		"redis.port":           "REDIS_PORT",
		"redis.password":       "REDIS_PASSWORD",
		"redis.db":             "REDIS_DB",
		"log.level":            "LOG_LEVEL",
		"log.format":           "LOG_FORMAT",
	}

	for key, env := range envBindings {
		_ = v.BindEnv(key, env)
	}

	// CORS_ORIGINS 特殊处理（逗号分隔字符串 → 切片）
	if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
		v.Set("cors.origins", strings.Split(origins, ","))
	}
}

func validate(cfg *Config) error {
	if cfg.Database.Name == "" {
		return fmt.Errorf("DB_NAME 不能为空")
	}
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET 不能为空")
	}
	if cfg.Encrypt.Key == "" {
		return fmt.Errorf("ENCRYPT_KEY 不能为空")
	}
	if len(cfg.Encrypt.Key) < 32 {
		return fmt.Errorf("ENCRYPT_KEY 长度不能少于 32 字节")
	}
	return nil
}
