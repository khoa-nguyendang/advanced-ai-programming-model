package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server                ServerConfig
	RabbitMQ              RabbitMQ
	Mysql                 MysqlConfig
	Redis                 RedisConfig
	RedisReplica          RedisReplicaConfig
	Cookie                Cookie
	Session               Session
	Metrics               Metrics
	Logger                Logger
	Jaeger                Jaeger
	AppEngine             AppEngine
	AppLightweightEngine  AppLightweightEngine
	AppFaiss              AppFaiss
	Mongodb               MongodbConfig
	MongoAtlas            MongoAtlasConfig
	UserService           string
	DeviceService         string
	LoggingService        string
	AuthenticationService string
	CompanyService        string
	AdministratorService  string
	VinHMSService         string
	PublisherServer       string
	MQTT                  Mqtt
	MinIO                 MinIO
	AndroidMinIO          AndroidMinIO
	Kafka                 Kafka
	ElasticSearch         ElasticSearch
	Aws                   Aws
}

// Server config struct
type ServerConfig struct {
	AppVersion               string
	Port                     string
	PprofPort                string
	Mode                     string
	JwtSecretKey             string
	JwtExpireInHour          int
	RefreshSecretKey         string
	RefreshTokenExpireInHour int
	CookieName               string
	ReadTimeout              time.Duration
	WriteTimeout             time.Duration
	SSL                      bool
	CtxDefaultTimeout        time.Duration
	CSRF                     bool
	Debug                    bool
	MaxConnectionIdle        time.Duration
	Timeout                  time.Duration
	MaxConnectionAge         time.Duration
	Time                     time.Duration
	CacheExpiryShort         time.Duration // default 10 mins
	CacheExpiryMedium        time.Duration // default 1 hours
	CacheExpiryLong          time.Duration // default 4 hours
	CacheExpiryDayLong       time.Duration // default 1 day
	HashKey                  string
	PassKey                  string
	IvKey                    string
}

// RabbitMQ
type RabbitMQ struct {
	Host           string
	Port           string
	User           string
	Password       string
	Exchange       string
	Queue          string
	RoutingKey     string
	ConsumerTag    string
	WorkerPoolSize int
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// MySQL config
type MysqlConfig struct {
	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPassword string
	MysqlDbname   string
	MysqlDriver   string
}

// MongoDB config
type MongodbConfig struct {
	MongodbHost     string
	MongodbPort     string
	MongodbUser     string
	MongodbPassword string
	MongodbDriver   string
	MongodbDbname   string
}

// MongoAtlas config
type MongoAtlasConfig struct {
	MongodbHost     string
	MongodbPort     string
	MongodbUser     string
	MongodbPassword string
	Mongodbname     string
	MongodbDriver   string
	MongodbDbname   string
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// Redis replica config
type RedisReplicaConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

// Cookie config
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Session config
type Session struct {
	Prefix string
	Name   string
	Expire int
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Jaeger
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// AppFaiss config
type AppFaiss struct {
	URL         string
	ServiceName string
}

// AppEngine config
type AppEngine struct {
	URL         string
	ServiceName string
}

type AppLightweightEngine struct {
	URL         string
	ServiceName string
}

// MinIO config
type MinIO struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

// MinIO config
type AndroidMinIO struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
}

// Kafka config
type Kafka struct {
	Server         string
	AuditLogServer string
	Username       string
	Password       string
}

// ElasticSearch config
type ElasticSearch struct {
	Server   string
	Username string
	Password string
}

type Mqtt struct {
	Broker   string
	Port     int32
	UserName string
	Password string
}

type Aws struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}

// Get config
func GetConfig(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config.dev"
}
