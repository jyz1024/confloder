package confloder

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/redis/go-redis/v9"
)

func buildHashConfigToJson(configMap map[string]string) string {
	builder := strings.Builder{}
	builder.Grow(1024)
	builder.Write([]byte{'{'})
	keyPairNum := len(configMap)
	rangePairNum := 0
	for cKey, cVal := range configMap {
		rangePairNum++
		if !json.Valid([]byte(cVal)) {
			continue
		}
		builder.WriteString("\"" + cKey + "\":" + cVal)
		if rangePairNum < keyPairNum {
			builder.Write([]byte{','})
		}
	}
	builder.Write([]byte{'}'})
	return builder.String()
}

func MakeStringConfigGetter(client *redis.Client, key string, data interface{}) LoadFunc {
	return func() (interface{}, error) {
		configStr, err := client.Get(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}
		configVal := reflect.New(reflect.TypeOf(data))
		err = json.Unmarshal([]byte(configStr), configVal.Interface())
		if err != nil {
			return nil, err
		}
		return configVal.Elem().Interface(), nil
	}
}

func MakeHashConfigGetter(client *redis.Client, key string, data interface{}) LoadFunc {
	return func() (interface{}, error) {
		configMap, err := client.HGetAll(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}
		configStr := buildHashConfigToJson(configMap)
		configVal := reflect.New(reflect.TypeOf(data))
		err = json.Unmarshal([]byte(configStr), configVal.Interface())
		if err != nil {
			return nil, err
		}
		return configVal.Elem().Interface(), nil
	}
}

func MakeHashFieldConfigGetter(client *redis.Client, key string, field string, data interface{}) LoadFunc {
	return func() (interface{}, error) {
		configStr, err := client.HGet(context.Background(), key, field).Result()
		if err != nil {
			return nil, err
		}
		configVal := reflect.New(reflect.TypeOf(data))
		err = json.Unmarshal([]byte(configStr), configVal.Interface())
		if err != nil {
			return nil, err
		}
		return configVal.Elem().Interface(), nil
	}
}
