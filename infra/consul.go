package infra

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"nocai/gokit-demo/infra/configs"
	"nocai/gokit-demo/infra/constants"
	"os"
)

func ConsulApi(l log.Logger, consulAddr string) *api.Client {
	consulConfig := api.DefaultConfig()
	if len(consulAddr) > 0 {
		consulConfig.Address = consulAddr
	}

	consulApi, err := api.NewClient(consulConfig)
	if err != nil {
		log.With(l, "stage", "consul api client").Log("err", err)
		os.Exit(1)
	}
	return consulApi
}

func ConsulRegister(l log.Logger, consulClient consul.Client, port int) {
	if err := consulClient.Register(&api.AgentServiceRegistration{
		ID:   fmt.Sprintf("%v:%v:%d", constants.AppName, LocalIP(), port),
		Name: constants.AppName,
		Port: port,
		Tags: []string{constants.AppName, "urlprefix-/" + constants.AppName + " strip=/" + constants.AppName},
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/ping", LocalIP(), port),
			Timeout:  "5s",
			Interval: "5s",
		},
	}); err != nil {
		log.With(l, "stage", "register consul").Log("error", "register to consul error:"+err.Error())
		os.Exit(-1)
	}
}

func ConsulKv(l log.Logger, consulApi *api.Client, consulAddr, configPath string) {
	kvPair, meta, err := consulApi.KV().Get(configPath, nil)
	if err != nil {
		log.With(l, "stage", "config consul").Log("error", err)
		os.Exit(-1)
	}
	l.Log("meta", meta)
	if kvPair == nil {
		log.With(l, "stage", "config consul").Log("error", fmt.Sprintf("kvPair[%s] is nil", constants.AppName))
		os.Exit(-1)
	}
	if err := configs.Unmarshal(kvPair.Value); err != nil {
		log.With(l, "stage", "config consul").Log("error", err)
		os.Exit(-1)
	}

	go func() {
		if r := recover(); r != nil {
			l.Log("error", r)
		}

		watchConfig(l, consulAddr)
	}()

}

func watchConfig(l log.Logger, consulAddr string) {
	params := map[string]interface{}{
		"type": "key",
		"key":  constants.AppName,
	}
	plan, err := watch.Parse(params)
	if err != nil {
		log.With(l, "stage", "watchConsul").Log("error", err)
		panic(err)
	}

	plan.Handler = func(idx uint64, raw interface{}) {
		if raw == nil {
			return // ignore
		}
		v, ok := raw.(*api.KVPair)
		if !ok || v == nil {
			return // ignore
		}
		log.With(l, "stage", "watchConsul").Log("kvPair", v)

		if err := configs.Unmarshal(v.Value); err != nil {
			log.With(l, "stage", "watchConsul").Log("error", err)
			panic(err)
		}
	}

	log.With(l, "stage", "watchConsul").Log("info", "watchConsul")
	if err := plan.Run(consulAddr); err != nil {
		log.With(l, "stage", "watchConsul").Log("error", err)
	}
}
