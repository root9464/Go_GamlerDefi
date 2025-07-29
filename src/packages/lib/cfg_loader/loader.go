package cfgloader

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	baseConfigsPath = "../../configs"
)

func NewLoader(opts ...Option) *UniversalConfig {
	v := viper.New()
	return &UniversalConfig{
		viper:       v,
		stopWatcher: make(chan struct{}),
	}
}

type ConfigFile struct {
	Name    string
	Path    string
	CfgType string
}

func (uc *UniversalConfig) LoadConfigs(cfgFiles []ConfigFile, dest any, env string) error {
	uc.viper.AutomaticEnv()
	uc.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	merged := make(map[string]any)

	for _, t := range cfgFiles {
		filePath := fmt.Sprintf("%s/%s.%s", t.Path, t.Name, t.CfgType)
		current, err := loadConfigWithEnvSubstitution(filePath)
		if err != nil {
			return err
		}

		merged = deepMerge(merged, current)
	}

	return UnmarshalWithCustomMap(dest, merged)
}

func loadConfigWithEnvSubstitution(filePath string) (map[string]any, error) {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file %s: %v", filePath, err)
	}

	content := resolveEnvVarsInString(string(contentBytes))

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(strings.NewReader(content)); err != nil {
		return nil, fmt.Errorf("cannot parse config after env substitution: %v", err)
	}

	return v.AllSettings(), nil
}

var envPattern = regexp.MustCompile(`\$\{(\w+)\}`)

func resolveEnvVarsInString(input string) string {
	return envPattern.ReplaceAllStringFunc(input, func(s string) string {
		key := envPattern.FindStringSubmatch(s)[1]
		return os.Getenv(key)
	})
}

func deepMerge(dst, src map[string]any) map[string]any {
	for k, v := range src {
		if vMap, ok := v.(map[string]any); ok {
			if dstMap, exists := dst[k].(map[string]any); exists {
				dst[k] = deepMerge(dstMap, vMap)
			} else {
				dst[k] = deepMerge(make(map[string]any), vMap)
			}
		} else {
			dst[k] = v
		}
	}
	return dst
}

func UnmarshalWithCustomMap(dest any, customMap map[string]any) error {
	decoderConfig := &mapstructure.DecoderConfig{
		Result:  dest,
		TagName: "mapstructure",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			stringToIntHookFunc(),
		),
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %v", err)
	}

	return decoder.Decode(customMap)
}

func stringToIntHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() == reflect.Int {
			return strconv.Atoi(data.(string))
		}
		return data, nil
	}
}
