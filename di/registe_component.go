package di

import (
	"github.com/gw123/jobsvr"
	"github.com/spf13/viper"
)

var di = jobsvr.NewDI(viper.GetViper())

func RegisterComponent(c jobsvr.Component) {
	di.Register(c)
}

func GetComponent(name string) (jobsvr.Component, bool) {
	return di.Get(name)
}
