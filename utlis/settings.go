package utlis

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

type Setting struct {
	Spider struct {
		Name     string `yaml:"SpiderName"`
		Worker   int    `yaml:"WorkerNumber"`
		TLS      bool   `yaml:"TLS"`
		LOGLEVEL string `yaml:"LOGLEVEL"`
	} `yaml:"Spider"`
	Headers map[string]string `yaml:"Headers"`
	Cookies map[string]string `yaml:"Cookies"`
}

type SettingManager struct {
	mu       sync.RWMutex
	Settings map[string]string
}

func NewSettingManager() *SettingManager {
	return &SettingManager{
		Settings: make(map[string]string),
	}
}
func (setting *SettingManager) GetSetting(key string) (string, bool) {
	setting.mu.RLock()
	defer setting.mu.RUnlock()
	value, ok := setting.Settings[key]
	return value, ok
}

func (setting *SettingManager) SetSetting(key string, value string) {
	setting.mu.Lock()
	defer setting.mu.Unlock()
	setting.Settings[key] = value
}

func (setting *SettingManager) GetInt(FlagName string, Flagdefault int) int {
	value, ok := setting.GetSetting(FlagName)
	if ok {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return Flagdefault
		}
		return intVal
	}
	return Flagdefault
}
func (setting *SettingManager) GetBool(FlagName string) bool {
	value, ok := setting.GetSetting(FlagName)
	if ok {
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return false
		}
		return boolVal
	}
	return false
}

func (setting *SettingManager) LoadFromSetting(s Setting) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		val := v.Field(i)

		prefix := field.Name

		switch val.Kind() {
		case reflect.Struct:
			for j := 0; j < val.NumField(); j++ {
				subField := val.Type().Field(j)
				subVal := val.Field(j)

				key := prefix + "." + subField.Name

				var strVal string
				switch subVal.Kind() {
				case reflect.String:
					strVal = subVal.String()
				case reflect.Bool:
					strVal = strconv.FormatBool(subVal.Bool())
				case reflect.Int, reflect.Int64:
					strVal = strconv.FormatInt(subVal.Int(), 10)
				default:
					strVal = fmt.Sprintf("%v", subVal.Interface())
				}

				setting.SetSetting(key, strVal)
			}
		case reflect.Map:
			for _, key := range val.MapKeys() {
				mapVal := val.MapIndex(key)
				if key.Kind() == reflect.String && mapVal.Kind() == reflect.String {
					fullKey := prefix + "." + key.String()
					setting.SetSetting(fullKey, mapVal.String())
				}
			}
		default:
			key := prefix
			strVal := fmt.Sprintf("%v", val.Interface())
			setting.SetSetting(key, strVal)
		}
	}
}
