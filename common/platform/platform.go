package platform

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type EnvFlag struct {
	Name    string
	AltName string
}

func (f EnvFlag) GetValue(defaultValue func() string) string {
	if v, found := os.LookupEnv(f.Name); found {
		return v
	}
	if len(f.AltName) > 0 {
		if v, found := os.LookupEnv(f.AltName); found {
			return v
		}
	}

	return defaultValue()
}

func (f EnvFlag) GetValueAsInt(defaultValue int) int {
	useDefaultValue := false
	s := f.GetValue(func() string {
		useDefaultValue = true
		return ""
	})
	if useDefaultValue {
		return defaultValue
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int(v)
}

func NormalizeEnvName(name string) string {
	return strings.Replace(strings.ToUpper(strings.TrimSpace(name)), ".", "_", -1)
}

func GetAssetLocation(file string) string {
	const name = "v2ray.location.asset"
	assetPath := EnvFlag{Name: name, AltName: NormalizeEnvName(name)}.GetValue(func() string {
		exec, err := os.Executable()
		if err != nil {
			return ""
		}
		return filepath.Dir(exec)
	})
	return filepath.Join(assetPath, file)
}
