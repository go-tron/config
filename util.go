package config

import (
	"os"
	"reflect"
	"regexp"
	"strings"
)

func ReadEnv(value string) string {
	path := "\\${(.*?)}"
	exp := regexp.MustCompile(path)
	result := exp.FindAllStringSubmatch(value, -1)
	if len(result) > 0 {
		for _, v := range result {
			if env := os.Getenv(v[1]); env != "" {
				value = strings.Replace(value, v[0], env, 1)
			}
		}
		return value
	}
	return ""
}

func ReplaceEnv(s interface{}) {
	if reflect.TypeOf(s).Kind() != reflect.Ptr {
		panic("s must be a pointer")
	}

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for k := 0; k < t.NumField(); k++ {
		if t.Field(k).Type.Kind() == reflect.Struct {
			ReplaceEnv(v.Field(k).Addr().Interface())
		} else {
			if t.Field(k).Type.Kind() == reflect.String {
				value := v.Field(k).Interface().(string)
				env := ReadEnv(value)
				if env != "" {
					v.Field(k).Set(reflect.ValueOf(env))
				}
			}
			if t.Field(k).Type.Kind() == reflect.Map {
				mv := reflect.ValueOf(v.Field(k).Interface())
				keys := mv.MapKeys()
				for _, mk := range keys {
					value := mv.MapIndex(mk)
					env := ReadEnv(value.Interface().(string))
					if env != "" {
						mv.SetMapIndex(mk, reflect.ValueOf(env))
					}
				}
			}
			//fmt.Printf("[%s]%s:%v\n", t.Field(k).Type.Kind(), t.Field(k).Name, v.Field(k).Interface())
		}
	}
}
