package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config アプリケーション設定
type Config struct {
	// 環境設定
	Environment string `json:"environment"`
	
	// AWS設定
	AWS AWSConfig `json:"aws"`
	
	// テーブル名
	Tables TableConfig `json:"tables"`
	
	// リトライ設定
	Retry RetryConfig `json:"retry"`
	
	// サーバー設定
	Server ServerConfig `json:"server"`
	
	// ログ設定
	Logging LoggingConfig `json:"logging"`
}

// AWSConfig AWS関連の設定
type AWSConfig struct {
	Region           string `json:"region"`
	DynamoDBEndpoint string `json:"dynamodb_endpoint"`
	Profile          string `json:"profile"`
	AccessKeyID      string `json:"access_key_id"`
	SecretAccessKey  string `json:"secret_access_key"`
}

// TableConfig テーブル名の設定
type TableConfig struct {
	Achievements   string `json:"achievements"`
	Rewards        string `json:"rewards"`
	CurrentPoints  string `json:"current_points"`
	RewardHistory  string `json:"reward_history"`
}

// RetryConfig リトライ設定
type RetryConfig struct {
	MaxRetries int `json:"max_retries"`
	BackoffMs  int `json:"backoff_ms"`
}

// ServerConfig サーバー設定
type ServerConfig struct {
	Port         string `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// LoggingConfig ログ設定
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// LoadConfig 設定ファイルと環境変数から設定を読み込み
func LoadConfig() (*Config, error) {
	// デフォルト設定
	config := getDefaultConfig()
	
	// 環境を取得
	env := getEnv("ENVIRONMENT", "development")
	config.Environment = env
	
	// 設定ファイルから読み込み
	if err := loadConfigFile(config, env); err != nil {
		// 設定ファイルが見つからない場合は警告のみ
		fmt.Printf("Warning: Could not load config file for environment '%s': %v\n", env, err)
	}
	
	// 環境変数で上書き
	overrideWithEnvVars(config)
	
	// 設定値の検証
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	return config, nil
}

// getDefaultConfig デフォルト設定を取得
func getDefaultConfig() *Config {
	return &Config{
		Environment: "development",
		AWS: AWSConfig{
			Region:           "us-east-1",
			DynamoDBEndpoint: "",
			Profile:          "",
		},
		Tables: TableConfig{
			Achievements:  "achievements",
			Rewards:       "rewards",
			CurrentPoints: "current_points",
			RewardHistory: "reward_history",
		},
		Retry: RetryConfig{
			MaxRetries: 3,
			BackoffMs:  100,
		},
		Server: ServerConfig{
			Port:         "8080",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}
}

// loadConfigFile 設定ファイルから設定を読み込み
func loadConfigFile(config *Config, env string) error {
	// 設定ファイルのパスを決定
	configPaths := []string{
		fmt.Sprintf("config/%s.json", env),
		fmt.Sprintf("configs/%s.json", env),
		fmt.Sprintf("%s.json", env),
	}
	
	var configData []byte
	var err error
	
	for _, path := range configPaths {
		if configData, err = os.ReadFile(path); err == nil {
			break
		}
	}
	
	if err != nil {
		return fmt.Errorf("config file not found for environment '%s'", env)
	}
	
	if err := json.Unmarshal(configData, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return nil
}

// overrideWithEnvVars 環境変数で設定を上書き
func overrideWithEnvVars(config *Config) {
	// 環境設定
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = env
	}
	
	// AWS設定
	if region := os.Getenv("AWS_REGION"); region != "" {
		config.AWS.Region = region
	}
	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		config.AWS.DynamoDBEndpoint = endpoint
	}
	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		config.AWS.Profile = profile
	}
	if accessKey := os.Getenv("AWS_ACCESS_KEY_ID"); accessKey != "" {
		config.AWS.AccessKeyID = accessKey
	}
	if secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY"); secretKey != "" {
		config.AWS.SecretAccessKey = secretKey
	}
	
	// テーブル名
	if table := os.Getenv("ACHIEVEMENTS_TABLE"); table != "" {
		config.Tables.Achievements = table
	}
	if table := os.Getenv("REWARDS_TABLE"); table != "" {
		config.Tables.Rewards = table
	}
	if table := os.Getenv("CURRENT_POINTS_TABLE"); table != "" {
		config.Tables.CurrentPoints = table
	}
	if table := os.Getenv("REWARD_HISTORY_TABLE"); table != "" {
		config.Tables.RewardHistory = table
	}
	
	// リトライ設定
	if retries := getEnvAsInt("MAX_RETRIES", 0); retries > 0 {
		config.Retry.MaxRetries = retries
	}
	if backoff := getEnvAsInt("RETRY_BACKOFF_MS", 0); backoff > 0 {
		config.Retry.BackoffMs = backoff
	}
	
	// サーバー設定
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if timeout := getEnvAsInt("SERVER_READ_TIMEOUT", 0); timeout > 0 {
		config.Server.ReadTimeout = timeout
	}
	if timeout := getEnvAsInt("SERVER_WRITE_TIMEOUT", 0); timeout > 0 {
		config.Server.WriteTimeout = timeout
	}
	
	// ログ設定
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.Logging.Level = level
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.Logging.Format = format
	}
	if output := os.Getenv("LOG_OUTPUT"); output != "" {
		config.Logging.Output = output
	}
}

// validateConfig 設定値の検証
func validateConfig(config *Config) error {
	var errors []string
	
	// 環境の検証
	validEnvs := []string{"development", "staging", "production"}
	if !contains(validEnvs, config.Environment) {
		errors = append(errors, fmt.Sprintf("invalid environment: %s (must be one of: %s)", 
			config.Environment, strings.Join(validEnvs, ", ")))
	}
	
	// AWS設定の検証
	if config.AWS.Region == "" {
		errors = append(errors, "AWS region is required")
	}
	
	// テーブル名の検証
	if config.Tables.Achievements == "" {
		errors = append(errors, "achievements table name is required")
	}
	if config.Tables.Rewards == "" {
		errors = append(errors, "rewards table name is required")
	}
	if config.Tables.CurrentPoints == "" {
		errors = append(errors, "current points table name is required")
	}
	if config.Tables.RewardHistory == "" {
		errors = append(errors, "reward history table name is required")
	}
	
	// リトライ設定の検証
	if config.Retry.MaxRetries < 0 {
		errors = append(errors, "max retries must be non-negative")
	}
	if config.Retry.BackoffMs < 0 {
		errors = append(errors, "backoff milliseconds must be non-negative")
	}
	
	// サーバー設定の検証
	if config.Server.Port == "" {
		errors = append(errors, "server port is required")
	}
	if config.Server.ReadTimeout <= 0 {
		errors = append(errors, "server read timeout must be positive")
	}
	if config.Server.WriteTimeout <= 0 {
		errors = append(errors, "server write timeout must be positive")
	}
	
	// ログ設定の検証
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, config.Logging.Level) {
		errors = append(errors, fmt.Sprintf("invalid log level: %s (must be one of: %s)", 
			config.Logging.Level, strings.Join(validLogLevels, ", ")))
	}
	
	validLogFormats := []string{"json", "text"}
	if !contains(validLogFormats, config.Logging.Format) {
		errors = append(errors, fmt.Sprintf("invalid log format: %s (must be one of: %s)", 
			config.Logging.Format, strings.Join(validLogFormats, ", ")))
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, "; "))
	}
	
	return nil
}

// contains スライスに要素が含まれているかチェック
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetConfigPath 設定ファイルのパスを取得
func GetConfigPath(env string) string {
	// 設定ファイルのパスを決定
	configPaths := []string{
		fmt.Sprintf("config/%s.json", env),
		fmt.Sprintf("configs/%s.json", env),
		fmt.Sprintf("%s.json", env),
	}
	
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	return fmt.Sprintf("config/%s.json", env)
}

// CreateConfigFile 設定ファイルを作成
func CreateConfigFile(env string) error {
	config := getDefaultConfig()
	config.Environment = env
	
	// 環境別の設定調整
	switch env {
	case "production":
		config.Logging.Level = "warn"
		config.Tables.Achievements = "prod-achievements"
		config.Tables.Rewards = "prod-rewards"
		config.Tables.CurrentPoints = "prod-current-points"
		config.Tables.RewardHistory = "prod-reward-history"
	case "staging":
		config.Logging.Level = "info"
		config.Tables.Achievements = "staging-achievements"
		config.Tables.Rewards = "staging-rewards"
		config.Tables.CurrentPoints = "staging-current-points"
		config.Tables.RewardHistory = "staging-reward-history"
	}
	
	configPath := GetConfigPath(env)
	
	// ディレクトリを作成
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// JSON形式で保存
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// getEnv 環境変数を取得（デフォルト値付き）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 環境変数を整数として取得（デフォルト値付き）
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}