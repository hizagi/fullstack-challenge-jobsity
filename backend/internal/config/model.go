package config

import (
	"fmt"
	"time"
)

const (
	defaultServerPort   = 8080
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
	defaultIdleTimeout  = 5 * time.Second
)

type ServiceConfig struct {
	Server *HTTPServer `koanf:"server"`
	DB     *DB         `koanf:"db"`
}

type HTTPServer struct {
	Port         int           `koanf:"port"`
	ReadTimeout  time.Duration `koanf:"readTimeout"`
	WriteTimeout time.Duration `koanf:"writeTimeout"`
	IdleTimeout  time.Duration `koanf:"idleTimeout"`
}

type DB struct {
	URI         string        `koanf:"uri"`
	User        string        `koanf:"user"`
	Password    string        `koanf:"password"`
	Database    string        `koanf:"database"`
	Timeout     time.Duration `koanf:"timeout"`
	MaxPoolSize uint64        `koanf:"maxPoolSize"`
}

func (c *ServiceConfig) HTTPServerConfig() HTTPServer {
	if c.Server == nil {
		return HTTPServer{
			Port:         defaultServerPort,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			IdleTimeout:  defaultIdleTimeout,
		}
	}

	return *c.Server
}

func (c *ServiceConfig) DBConfig() (DB, error) {
	if c.Server == nil {
		return DB{}, fmt.Errorf("db config is mandatory")
	}

	return *c.DB, nil
}
