package configuration

import (
	"fmt"

	"github.com/spf13/viper"
)

/*
Bind bind config to struct

	@param key: key of config
	@param path: path to config file
	@param file: name of config file
	@param ext: extension of config file

	@return T: config struct
	@return error: error if any
*/
func Bind[T any](key, path, file, ext string) (T, error) {
	var result T

	viper.SetConfigName(file)
	viper.AddConfigPath(path)
	viper.SetConfigType(ext)

	if err := viper.ReadInConfig(); err != nil {
		return result, fmt.Errorf("failed to read config: %v", err)
	}

	if key == "" {
		if err := viper.Unmarshal(&result); err != nil {
			return result, fmt.Errorf("failed to unmarshal config: %v", err)
		}
	} else {
		if err := viper.UnmarshalKey(key, &result); err != nil {
			return result, fmt.Errorf("failed to unmarshal config by key: %v", err)
		}
	}

	return result, nil
}
