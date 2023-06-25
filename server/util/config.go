package util

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Config stores all configuration for the application to init
// Values are read by viper using config file / environments variables
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmatricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	CookieDomainEndpoint string        `mapstructure:"COOKIE_DOMAIN_ENDPOINT"`
	AllowOriginEndpoint  string        `mapstructure:"ALLOW_ORIGIN_ENDPOINT"`
	OpenAIToken          string        `mapstructure:"OPENAI_TOKEN"`
	S3BucketName         string        `mapstructure:"S3_BUCKET_NAME"`
	AwsRegion            string        `mapstructure:"AWS_DEFUALT_REGION"`
}

type initConfig struct {
	ConfigFileName string
	ConfigFilePath string
}

type initConfigOption func(*initConfig)

func WithConfigFileName(configFileName string) initConfigOption {
	return func(c *initConfig) {
		c.ConfigFileName = configFileName
	}
}

func WithConfigFilePath(configFilePath string) initConfigOption {
	return func(c *initConfig) {
		c.ConfigFilePath = configFilePath
	}
}

func LoadConfig(options ...initConfigOption) (config Config, err error) {
	var initConfig = &initConfig{
		ConfigFilePath: ".",
		ConfigFileName: "conf",
	}

	for _, opt := range options {
		opt(initConfig)
	}

	viper.AddConfigPath(initConfig.ConfigFilePath)
	viper.SetConfigName(initConfig.ConfigFileName)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	enrichConfig(&config)

	return
}

// Validate makes sure each configuration is set correctly
func (c *Config) Validate() (errMsg error) {
	valueType := reflect.TypeOf(*c)
	value := reflect.ValueOf(*c)
	var missingConfigs []interface{}

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		fieldValue := value.Field(i)
		fieldTag := field.Tag.Get("mapstructure")

		if fieldValue.Interface() == reflect.Zero(fieldValue.Type()).Interface() {
			fmt.Println("here", fieldTag)
			missingConfigs = append(missingConfigs, fieldTag)
		}
	}
	if len(missingConfigs) > 0 {
		errMsg = buildEnvVarErrorMsg(missingConfigs)
	}

	return
}

// buildEnvVarErrorMsg get's the missing env vars config and build the error message
func buildEnvVarErrorMsg(missingConfigs []interface{}) error {
	var sb strings.Builder
	var i int

	for ; i < len(missingConfigs)-1; i++ {
		sb.WriteString(fmt.Sprintf("%v,", missingConfigs[i]))
	}

	sb.WriteString(fmt.Sprintf("%v", missingConfigs[i]))

	return errors.New(fmt.Sprintf("%v is required", sb.String()))

}

func enrichConfig(c *Config) error {
	valueType := reflect.TypeOf(c).Elem()
	value := reflect.ValueOf(c).Elem()

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		fieldValue := value.Field(i)
		fieldTag := field.Tag.Get("mapstructure")

		if fieldValue.Interface() == reflect.Zero(fieldValue.Type()).Interface() {
			val, exists := getEnvVar(fieldTag)
			if !exists {
				continue

			}

			err := setFieldValue(&fieldValue, val)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// getEnvVar returns the requested env var and if it's exists
func getEnvVar(fieldName string) (string, bool) {
	fieldName = strings.ToUpper(fieldName)
	value, exists := os.LookupEnv(fieldName)
	return value, exists
}

func setFieldValue(fieldValue *reflect.Value, val string) error {
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetInt(intVal)
	case reflect.Int64:
		if fieldValue.Type() == reflect.TypeOf(time.Duration(0)) {
			durationVal, err := time.ParseDuration(val)
			if err != nil {
				return err
			}
			fieldValue.SetInt(int64(durationVal))
		} else {
			intVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		fieldValue.SetBool(boolVal)
	case reflect.Struct:
		// Handle time.Time field type
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			timeVal, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(timeVal))
		}
	default:
		return errors.New("unsupported field type")
	}

	return nil
}
