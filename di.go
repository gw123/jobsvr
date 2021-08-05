package jobsvr

import (
	"errors"

	"github.com/gw123/glog"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type InitComponent func() Component

type Component interface {
	Name() string                        // 组件名称
	Init() error                         // 组件初始化方式
	Close() error                        // 组件关闭方式
	Reload()                             // 当必要配置发生变化时候组件重新加载
	GetWatchConfigKeys() WatchConfigKeys // 组件所需要的配置实时更新组件状态
	OnConfigChange(key, val string)      // 配置发生变化时候调用
}

type DI struct {
	data              map[string]Component
	watchConfig       map[string][]Component
	watchReloadConfig map[string][]Component
	viper             *viper.Viper
}

func NewDI(v *viper.Viper) *DI {
	return &DI{
		viper:             v,
		watchReloadConfig: make(map[string][]Component),
		watchConfig:       make(map[string][]Component),
		data:              make(map[string]Component),
	}
}

func (di *DI) Register(c Component) error {
	if _, ok := di.data[c.Name()]; ok {
		return errors.New("component name already register")
	}
	di.data[c.Name()] = c

	wk := c.GetWatchConfigKeys()

	for _, key := range wk.NeedUpdateConfigKeys {
		di.watchConfig[key] = append(di.watchConfig[key], c)
	}

	for _, key := range wk.NeedReloadKeys {
		di.watchReloadConfig[key] = append(di.watchReloadConfig[key], c)
	}

	return nil
}

func (di *DI) Get(name string) (Component, bool) {
	c, ok := di.data[name]
	if !ok {
		glog.DefaultLogger().Infof("component not register : %s", name)
	}
	return c, ok
}

func (di *DI) Watch() {
	di.viper.OnConfigChange(func(in fsnotify.Event) {
		glog.DefaultLogger().Infof("viper on config change %+v", in)
		key := in.String()
		if _, ok := di.watchConfig[key]; ok {
			for _, component := range di.watchConfig[key] {
				component.OnConfigChange(key, viper.GetString(key))
			}
		}

		if _, ok := di.watchReloadConfig[key]; ok {
			for _, component := range di.watchReloadConfig[key] {
				component.Reload()
			}
		}
	})
}

type WatchConfigKeys struct {
	NeedUpdateConfigKeys []string
	NeedReloadKeys       []string
}
