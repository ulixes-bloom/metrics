// Certgenerator is a tool that generates random RSA private and public keys.
//
// Usage:
//
//	go build .
//	./certgenerator -c="public.pem" -k="private.pem"
package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/ulixes-bloom/ya-metrics/internal/pkg/rsa"
)

type Config struct {
	PublicKeyFile  string
	PrivateKeyFile string
}

func GetDefaultConfig() *Config {
	return &Config{
		PublicKeyFile:  "public.pem",
		PrivateKeyFile: "private.pem",
	}
}

func ParseConfig() (*Config, error) {
	var conf Config
	defaultValues := GetDefaultConfig()

	flag.StringVar(&conf.PrivateKeyFile, "k", defaultValues.PrivateKeyFile, "path to save generated private key")
	flag.StringVar(&conf.PublicKeyFile, "c", defaultValues.PublicKeyFile, "path to save generated public key")
	flag.Parse()

	return &conf, nil
}

func main() {
	conf, err := ParseConfig()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	publicKeyPEM, privateKeyPEM, err := rsa.GenerateKeyPair()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = os.WriteFile(conf.PrivateKeyFile, privateKeyPEM, 0644)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = os.WriteFile(conf.PublicKeyFile, publicKeyPEM, 0644)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
