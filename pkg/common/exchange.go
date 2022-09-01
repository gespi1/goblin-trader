package common

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func CheckForToken(v *viper.Viper, tokenName string) {
	if v.GetString(tokenName) == "" {
		log.Warnf("token %v is not set", tokenName)
	}
}
