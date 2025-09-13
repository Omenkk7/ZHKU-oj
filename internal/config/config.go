package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	MongoDB  MongoDBConfig  `yaml:"mongodb"`
	Redis    RedisConfig    `yaml:"redis"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	Judge    JudgeConfig    `yaml:"judge"`
	JWT      JWTConfig      `yaml:"jwt"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string        `yaml:"port"`
	Mode         string        `yaml:"mode"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI            string        `yaml:"uri"`
	Database       string        `yaml:"database"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	MaxPoolSize    uint64        `yaml:"max_pool_size"`
	MinPoolSize    uint64        `yaml:"min_pool_size"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	VHost    string `yaml:"vhost"`
}

// JudgeConfig 判题配置
type JudgeConfig struct {
	Sandboxes      []SandboxConfig      `yaml:"sandboxes"`
	Compile        CompileConfig        `yaml:"compile"`
	Runtime        RuntimeConfig        `yaml:"runtime"`
	FileManagement FileManagementConfig `yaml:"file_management"`
}

// SandboxConfig 沙箱配置
type SandboxConfig struct {
	URL                 string        `yaml:"url"`
	Weight              int           `yaml:"weight"`
	MaxConcurrent       int           `yaml:"max_concurrent"`
	Timeout             time.Duration `yaml:"timeout"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
}

// CompileConfig 编译配置
type CompileConfig struct {
	Java JavaCompileConfig `yaml:"java"`
}

// JavaCompileConfig Java编译配置
type JavaCompileConfig struct {
	Command     []string `yaml:"command"`
	Env         []string `yaml:"env"`
	CPULimit    int64    `yaml:"cpu_limit"`
	MemoryLimit int64    `yaml:"memory_limit"`
	ProcLimit   int      `yaml:"proc_limit"`
	FileTimeout int      `yaml:"file_timeout"`
}

// RuntimeConfig 运行时配置
type RuntimeConfig struct {
	Java JavaRuntimeConfig `yaml:"java"`
}

// JavaRuntimeConfig Java运行时配置
type JavaRuntimeConfig struct {
	Command     []string `yaml:"command"`
	Env         []string `yaml:"env"`
	CPULimit    int64    `yaml:"cpu_limit"`
	MemoryLimit int64    `yaml:"memory_limit"`
	ProcLimit   int      `yaml:"proc_limit"`
	OutputLimit int      `yaml:"output_limit"`
}

// FileManagementConfig 文件管理配置
type FileManagementConfig struct {
	CleanupInterval time.Duration `yaml:"cleanup_interval"`
	MaxCacheSize    string        `yaml:"max_cache_size"`
	AutoCleanup     bool          `yaml:"auto_cleanup"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string     `yaml:"level"`
	Format string     `yaml:"format"`
	Output string     `yaml:"output"`
	File   FileConfig `yaml:"file"`
}

// FileConfig 文件日志配置
type FileConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// Load 加载配置文件
func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         ":8080",
			Mode:         "debug",
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
		},
		MongoDB: MongoDBConfig{
			URI:            "mongodb://localhost:27017",
			Database:       "campus_oj",
			ConnectTimeout: 10 * time.Second,
			MaxPoolSize:    50,
			MinPoolSize:    5,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
			PoolSize: 20,
		},
		RabbitMQ: RabbitMQConfig{
			Host:     "localhost",
			Port:     5672,
			Username: "guest",
			Password: "guest",
			VHost:    "/",
		},
		Judge: JudgeConfig{
			Sandboxes: []SandboxConfig{
				{
					URL:                 "http://localhost:5050",
					Weight:              1,
					MaxConcurrent:       10,
					Timeout:             30 * time.Second,
					HealthCheckInterval: 10 * time.Second,
				},
			},
			Compile: CompileConfig{
				Java: JavaCompileConfig{
					Command:     []string{"/usr/bin/javac"},
					Env:         []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
					CPULimit:    10000000000, // 10秒
					MemoryLimit: 268435456,   // 256MB
					ProcLimit:   50,
					FileTimeout: 60,
				},
			},
			Runtime: RuntimeConfig{
				Java: JavaRuntimeConfig{
					Command:     []string{"/usr/bin/java"},
					Env:         []string{"PATH=/usr/bin:/bin", "JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64"},
					CPULimit:    5000000000, // 5秒
					MemoryLimit: 134217728,  // 128MB
					ProcLimit:   1,
					OutputLimit: 10240, // 10KB
				},
			},
			FileManagement: FileManagementConfig{
				CleanupInterval: 5 * time.Minute,
				MaxCacheSize:    "1GB",
				AutoCleanup:     true,
			},
		},
		JWT: JWTConfig{
			Secret: "your-secret-key-change-in-production",
			Expire: 24 * time.Hour,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
			File: FileConfig{
				Enabled:    true,
				Path:       "logs/app.log",
				MaxSize:    100,
				MaxBackups: 10,
				MaxAge:     30,
			},
		},
	}
}
