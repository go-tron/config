package config

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var configPflag = pflag.StringP("config", "c", "development", "specify the config file")
var configFlag = flag.String("config", "development", "specify the config file")
var configShortFlag = flag.String("c", "development", "specify the config file")

func New() *Config {
	v := viper.New()
	v.SetDefault("application.env", "development")
	v.SetDefault("cluster.namespace", "")
	v.SetDefault("cluster.nodeName", "node-0")
	v.SetDefault("cluster.nodeIP", "127.0.0.1")
	v.SetDefault("cluster.podName", "localhost-0")
	v.SetDefault("cluster.podIP", "127.0.0.1")
	v.SetDefault("cluster.dnsSuffix", "")

	v.SetDefault("logging.path", "./logs")
	v.SetDefault("logging.maxSize", 128)
	v.SetDefault("logging.maxBackups", 30)
	v.SetDefault("logging.maxAge", 7)
	v.SetDefault("logging.compress", true)

	v.SetDefault("database.maxIdleConns", 5)
	v.SetDefault("database.maxOpenConns", 0)

	v.SetDefault("redis.poolSize", 10)
	v.SetDefault("redis.minIdleConns", 2)

	v.SetEnvPrefix("GO_OPT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.BindEnv("cluster.namespace", "CLUSTER_NAMESPACE")
	v.BindEnv("cluster.nodeName", "CLUSTER_NODENAME")
	v.BindEnv("cluster.nodeIP", "CLUSTER_NODEIP")
	v.BindEnv("cluster.podName", "HOSTNAME")
	v.BindEnv("cluster.podIP", "CLUSTER_PODIP")

	flagSet := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	flagSet.StringP("config", "c", "development", "specify the config file")
	flagSet.Parse(os.Args[1:])
	v.BindPFlags(flagSet)

	v.AddConfigPath(".")
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err != nil {
		log.Println("config read error", err)
	}

	if v.GetString("config") != "" {
		v.SetConfigName("config." + v.GetString("config"))
	}

	if err := v.MergeInConfig(); err != nil {
		log.Println("config merge error", err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config file changed:", e.Name)
	})

	if ns := v.GetString("cluster.namespace"); ns != "" {
		v.Set("logging.path", fmt.Sprintf("%s/%s/%s", v.GetString("logging.path"), ns, v.GetString("application.name")))
	}

	nodes := v.GetStringMapString("cluster.nodes")
	if len(nodes) > 0 {
		for node, ip := range nodes {
			if node == v.GetString("cluster.nodeName") {
				v.SetDefault("cluster.hostIP", ip)
			}
		}
	}

	log.Println("config application.env", v.GetString("application.env"))
	log.Println("config cluster.namespace", v.GetString("cluster.namespace"))
	log.Println("config cluster.nodeName", v.GetString("cluster.nodeName"))
	log.Println("config cluster.nodeIP", v.GetString("cluster.nodeIP"))
	log.Println("config cluster.hostIP", v.GetString("cluster.hostIP"))
	log.Println("config cluster.podName", v.GetString("cluster.podName"))
	log.Println("config cluster.podIP", v.GetString("cluster.podIP"))
	log.Println("config cluster.dnsSuffix", v.GetString("cluster.dnsSuffix"))
	log.Println("config logging.path", v.GetString("logging.path"))

	return &Config{v}
}

type Config struct {
	*viper.Viper
}

//func (c *Config) Get(key string) interface{} {
//	value := c.Viper.Get(key)
//	if value == nil {
//		return value
//	}
//	if reflect.TypeOf(value).Kind() == reflect.String && strings.Contains(value.(string), "${") {
//		envValue := ReadEnv(value.(string))
//		if envValue != "" {
//			return envValue
//		} else {
//			return value
//		}
//	} else {
//		return value
//	}
//}
//
//func (c *Config) GetString(key string) string {
//	return cast.ToString(c.Get(key))
//}
//
//func (c *Config) GetBool(key string) bool {
//	return cast.ToBool(c.Get(key))
//}
//
//func (c *Config) GetInt(key string) int {
//	return cast.ToInt(c.Get(key))
//}
//
//func (c *Config) GetInt32(key string) int32 {
//	return cast.ToInt32(c.Get(key))
//}
//
//func (c *Config) GetInt64(key string) int64 {
//	return cast.ToInt64(c.Get(key))
//}
//
//func (c *Config) GetUint(key string) uint {
//	return cast.ToUint(c.Get(key))
//}
//
//func (c *Config) GetUint32(key string) uint32 {
//	return cast.ToUint32(c.Get(key))
//}
//
//func (c *Config) GetUint64(key string) uint64 {
//	return cast.ToUint64(c.Get(key))
//}
//
//func (c *Config) GetFloat64(key string) float64 {
//	return cast.ToFloat64(c.Get(key))
//}
//
//func (c *Config) GetTime(key string) time.Time {
//	return cast.ToTime(c.Get(key))
//}
//
//func (c *Config) GetDuration(key string) time.Duration {
//	return cast.ToDuration(c.Get(key))
//}
