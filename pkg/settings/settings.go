package settings

import (
	"time"
)

const (
	Debug   string = "debug"
	Release string = "release"
	Test    string = "test"
)

type Settings struct {
	Global struct {
		Name              string        `yaml:"name" env:"GLOBAL_NAME" env-default:"maani" env-description:"Instance Name"`
		ReadTimeout       time.Duration `yaml:"readTimeout" env:"GLOBAL_READ_TIMEOUT" env-default:"2m" env-description:"Read timeout of http server"`
		ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" env:"GLOBAL_READ_HEADER_TIMEOUT" env-default:"2m" env-description:"Read header timeout of http server"`
		WriteTimeout      time.Duration `yaml:"writeTimeout" env:"GLOBAL_WRITE_TIMEOUT" env-default:"2m" env-description:"Write timeout of http server"`
		IdleTimeout       time.Duration `yaml:"idleTimeout" env:"GLOBAL_IDLE_TIMEOUT" env-default:"2m" env-description:"Idle timeout of http server"`
		MaxHeaderBytes    int           `yaml:"maxHeaderBytes" env:"GLOBAL_MAX_HEADER_BYTES" env-default:"8196" env-description:"Max header bytes of http server"`
		BackendPort       int           `yaml:"backendPort" env:"GLOBAL_BACKEND_PORT" env-default:"9000" env-description:"Port of backend server"`
		AdminPort         int           `yaml:"adminPort" env:"GLOBAL_ADMIN_PORT" env-default:"9001" env-description:"Port of admin server"`
		GatewayPort       int           `yaml:"gatewayPort" env:"GLOBAL_GATEWAY_PORT" env-default:"8000" env-description:"Port of gateway server"`
		Environment       string        `yaml:"environment" env:"CONFIG_MODE" env-default:"file" env-description:"Execution mode of Gin framework"`
	} `yaml:"global"`
	Database struct {
		Type                string        `yaml:"type" env:"CONFIG_DB_TYPE" env-default:"pgsql" env-description:"Postgres connection mode"`
		SSLMode             string        `yaml:"sslMode" env:"CONFIG_DB_SSL_MODE" env-default:"disable" env-description:"Database connection mode"`
		Host                string        `yaml:"host" env:"CONFIG_DB_HOST" env-default:"127.0.0.1" env-description:"Database connection host"`
		User                string        `yaml:"user" env:"CONFIG_DB_USER" env-default:"postgres" env-description:"Database connection user"`
		DatabaseName        string        `yaml:"databaseName" env:"CONFIG_DB_DATABASE" env-default:"postgres" env-description:"Database name"`
		Password            string        `yaml:"password" env:"CONFIG_DB_PASSWORD" env-description:"Database connection password"`
		Port                int           `yaml:"port" env:"CONFIG_DB_PORT" env-default:"5432" env-description:"Database connection port"`
		QueryTimeout        time.Duration `yaml:"queryTimeout" env:"CONFIG_DB_QUERY_TIMEOUT" env-default:"5s" env-description:"Database query timeout"`
		SyncInterval        time.Duration `yaml:"syncInterval" env:"CONFIG_DB_SYNC_INTERVAL" env-default:"60s" env-description:"Sync interval of Configs in db mode"`
		MaxOpenConnections  int           `yaml:"maxOpenConnections" env:"CONFIG_DB_MAX_OPEN_CONNECTIONS" env-default:"50" env-description:"The maximum number of open connections to database"`
		MaxIdleConnections  int           `yaml:"maxIdleConnections" env:"CONFIG_DB_MAX_IDLE_CONNECTIONS" env-default:"20" env-description:"The maximum number of idle connections to database"`
		ConnMaxLifetime     time.Duration `yaml:"connMaxLifetime" env:"CONFIG_DB_CONN_MAX_LIFETIME" env-default:"60s" env-description:"The maximum amount of time a connection may be reused"`
		ConnMaxIdleTime     time.Duration `yaml:"connMaxIdleTime" env:"CONFIG_DB_CONN_MAX_IDLE_TIME" env-default:"3s" env-description:"The maximum amount of time a connection can be idle"`
		StatusCheckInterval time.Duration `yaml:"statusCheckInterval" env:"CONFIG_DBMETRIC_CHECK_INTERVAL" env-default:"5s" env-description:"Interval for checking db connection status"`
	} `yaml:"database"`
	GatewayServer struct {
		StoreHost           string        `yaml:"storeHost" env:"STORE_HOST" env-default:"http://store" env-description:"Host for request to store servers"`
		SecretKey           string        `env:"GATEWAY_API_SECRET_KEY" env-default:"gatewaySecret" env-description:"Secret key for gateway server api authentication"`
		TokenTimeout        time.Duration `yaml:"tokenTimeout" env:"API_TOKEN_TIMEOUT" env-default:"1h" env-description:"Timeout of token for api authentication"`
		RefreshTokenTimeout time.Duration `yaml:"refreshTokenTimeout" env:"API_REFRESH_TOKEN_TIMEOUT" env-default:"3h" env-description:"Timeout of refresh token for api authentication"`
		UserIdHeaderKey     string        `yaml:"userIdHeaderKey" env:"USER_ID_HEADER_KEY" env-default:"X-MAANI-USER" env-description:"Header key to set user id and pass it throw reequest"`
	} `yaml:"retreival"`
	BackendServer struct {
		EncryptKey       string `yaml:"encryptKey" env:"ENCRYPT_KEY" env-default:"files-secret-key"  env-description:"Key for encrypting file"`
		FilePath         string `yaml:"filePath" env:"FILE_PATH" env-default:"/opt/files" env-description:"Path for new file to save"`
		MaxFilesSizeByte int    `yaml:"maxFilesSizeByte" env:"MAX_FilES_SIZE_BYTE" env-default:"100000000" env-description:"Maximum limitation of files size in byte"`
	} `yaml:"store"`
}

func (settings Settings) IsValid() (bool, error) {
	if settings.Global.Name == "" {
		return false, ErrSettingNameEmpty
	}
	// Check Ports duplications
	var hasDuplicatedPort bool
	duplicatedPorts := make(map[int]bool)
	for _, item := range []int{
		settings.Global.BackendPort,
		settings.Global.AdminPort,
	} {
		_, exist := duplicatedPorts[item]
		if exist {
			hasDuplicatedPort = true
			duplicatedPorts[item] = true
		} else {
			duplicatedPorts[item] = false
		}
	}
	if hasDuplicatedPort {
		return false, ErrSettingDuplicatedServerPorts
	}

	if settings.Global.Environment != Debug && settings.Global.Environment != Release && settings.Global.Environment != Test {
		return false, ErrSettingInvalidEnvironment
	}
	return true, nil
}
