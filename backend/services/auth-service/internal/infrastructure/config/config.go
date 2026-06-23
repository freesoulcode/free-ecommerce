package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultServiceName = "auth-service"
	defaultEnv         = "local"
	defaultHTTPAddr    = ":8081"
	defaultLogLevel    = "info"
)

type Config struct {
	ServiceName string
	Env         string
	HTTPAddr    string
	GRPCAddr    string
	LogLevel    string
	MySQL       MySQLConfig
	JWT         JWTConfig
	Snowflake   SnowflakeConfig
}

type MySQLConfig struct {
	DSN string
}

type JWTConfig struct {
	Issuer           string
	Audience         string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	RSAPrivateKeyPEM string
}

type SnowflakeConfig struct {
	Node int64
}

func Load() Config {
	return Config{
		ServiceName: getEnv("AUTH_SERVICE_NAME", defaultServiceName),
		Env:         getEnv("AUTH_SERVICE_ENV", defaultEnv),
		HTTPAddr:    getEnv("AUTH_SERVICE_HTTP_ADDR", defaultHTTPAddr),
		GRPCAddr:    getEnv("AUTH_SERVICE_GRPC_ADDR", ":9081"),
		LogLevel:    getEnv("AUTH_SERVICE_LOG_LEVEL", defaultLogLevel),
		MySQL: MySQLConfig{
			DSN: getEnv("AUTH_SERVICE_MYSQL_DSN", "root:password@tcp(127.0.0.1:30306)/auth_service?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		JWT: JWTConfig{
			Issuer:           getEnv("AUTH_SERVICE_JWT_ISSUER", "auth-service"),
			Audience:         getEnv("AUTH_SERVICE_JWT_AUDIENCE", "buyer-api"),
			AccessTokenTTL:   getEnvDurationMinutes("AUTH_SERVICE_ACCESS_TOKEN_TTL_MINUTES", 15),
			RefreshTokenTTL:  getEnvDurationHours("AUTH_SERVICE_REFRESH_TOKEN_TTL_HOURS", 24*30),
			RSAPrivateKeyPEM: getEnv("AUTH_SERVICE_JWT_PRIVATE_KEY_PEM", defaultRSAPrivateKeyPEM),
		},
		Snowflake: SnowflakeConfig{
			Node: getEnvInt64("AUTH_SERVICE_SNOWFLAKE_NODE", 2),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvDurationMinutes(key string, fallbackMinutes int64) time.Duration {
	return time.Duration(getEnvInt64(key, fallbackMinutes)) * time.Minute
}

func getEnvDurationHours(key string, fallbackHours int64) time.Duration {
	return time.Duration(getEnvInt64(key, fallbackHours)) * time.Hour
}

const defaultRSAPrivateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAoMLiE28DGl8hebX3cDUZXxMnnmRLFkSI4VO5Lg/LJwe6AhkN
bwurVKzRIhFtjyMKPHemYL7+2XvVnwwUIcI+90OypPrSfIP1akBpoVMN+UT4aRDx
TfNb1rQuwnaaGpUuLoSUzWa3T6UbQG3jk16hmZlFyNjOIC8pqEeKK7Y1p0F3m693
b8+LNlKEw630wyLLgwn+AtBTg6mY0r477yxkrLxe+CwBCZDz4btR1nV6RZ7IBI2P
ZERz7So2rQENolZtLdy2smSZrAOZDnwpwgKeVpAzEF9Zo/ZwmbA0okwmdP7AruS5
7xXRVimd7ReEKXampl5sIHMzsnJkVsXxSQW5owIDAQABAoIBAQCMPLUeos6gKLB5
DgXF+mwhhgIfp/ngePS3K2P1DI35hEH9JoGThyh0ezUMdQuPu89oJDAdYT/L1Lzr
O4wsTtjCtmmWhb8sI6jogTwkIOGlu0a/0KnPiCVrTE8mEHQqEEzzA3ETJTFv5uW4
9KN7oSdzaEN6C7b0WHAMfivIKfDv70fayWkgOUUTgpri3t+JPGUehsKXCeFCE8AH
OMMxgKcbFBCvmYjF1+v0he/+uejNGc8AOH2gISnVZ5krgrfykTUV73m5MRHwESqE
iCyC+16CeOFti9LBFsIyZ30fBVR0s468d5j+m7HbCn+nC+TiK5QiSKfLgWCCtBNe
BA/Z44dRAoGBANUhOR5nG2hrWo+WYSfif1lpSZS22nMWXXhs1Y8PVRQ0RmTSipDX
xL7eEIYmLwcU5UBd02C0ixmyHqBXX2QGTe2pjhM8znvb6Tq0wctqFMrin7rG7qnx
92V7EvqtJ42PHZjr9/XD51tuppUwErmC6e3Kum8X2b9xtWGe5M2dsjKtAoGBAMEZ
Bu4DtF80btENRDEnjOS/eCIpgW7HWhH3TUhDq+Tnf3ZbEY23J5MTFJsAjg/t5OuG
j32zP8Kt0NErNWhsf/glRnakS8dZWpQhkDxsQlNi7RgQfnMvqlxxBr7ASLxs3pHS
/7pew3SKLrQnQ0RnVoVDuIFz2SNq/3P2mV5bQXePAoGAYRBjahQ9KD4UHWa4UqjV
pMvNpfvs2xMpeIngbOnnrm7sTEiSsMqDoQWTcvT63/fFPJ4+gUFYRFiZmB6SpAQ3
A3D/8oTz6PbLbmAaDmD+nTO+2Rp2YVGAgWgeyamIZPDz4sw8vmH9AOgQ18rwDCqy
DQkSBTxQf97yY0YxH++c03UCgYEAgHKwfGW0d1w+lwt3ICeJ/qQrOrZXZiRwEuFp
5Dc3wiYIUOfVbmq2hYw8ubsNxSTfkZjKHLi/IjZTYMCYX2VFXwEUtVknG22h5kXJ
V5hAKo3033whUWgUsDdzYDIycD0PdPthp0zgQcaluKshgQAouq9IrbwtZfUIBtC0
RuL3UpsCgYEAriU6hP8NhDdTB9n8/vu7LPBhs9jCTtZkeqagjDosFYpIbUPnAwKG
2W0ZfUDXbwuD1emsMEamvLBMZ3ggYykAkaXqQ9MaweRtEPtd94boDMmftyyf+5j3
BvQzTPw1KXcrCC5GsJ/TLOr5pwEcO52tpslfLFcGByxHOu1z5ig+nBE=
-----END RSA PRIVATE KEY-----`
