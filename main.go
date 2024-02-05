package main

import (
	"log"
	"os"
	"time"

	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	ld "github.com/launchdarkly/go-server-sdk/v7"
	"github.com/launchdarkly/go-server-sdk/v7/ldcomponents"
	"github.com/pkg/errors"
)

func DefaultConfig(appId string, appVersion string) *ld.Config {
	var config ld.Config
	config.Logging = ldcomponents.NoLogging()
	config.ApplicationInfo.ApplicationID = appId
	config.ApplicationInfo.ApplicationVersion = appVersion

	return &config
}

func AnonymousContext(val string) ldcontext.Context {
	context := ldcontext.NewBuilder(val).Anonymous(true).Build()

	return context
}

func EnvContext(env string) ldcontext.Context {
	var b ldcontext.Builder
	context := b.Kind("env").Key(env).Name(env).Build()

	return context
}

func ldClient(key string, config *ld.Config, timeout time.Duration) (*ld.LDClient, error) {
	ldclient, err := ld.MakeCustomClient(
		key,
		*config,
		timeout,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot initiate LaunchDarkly instance")
	}

	return ldclient, nil
}

func main() {
	var context ldcontext.Context

	env := os.Args[1]
	if env == "user" {
		context = AnonymousContext("61b89b779bf2cf2b5b1b73f7")
	} else {
		context = EnvContext(env)
	}

	flag := "enable-scheduler-shardmanager"
	defaultVal := false
	val := defaultVal
	apiKey := os.Getenv("LD_SDK_KEY")

	ldc, err := ldClient(apiKey, DefaultConfig("recart-test", "1.0.0"), 50*time.Second)
	if err == nil {
		val, err = ldc.BoolVariation(flag, context, defaultVal)
	}
	if err != nil {
		log.Println("FeatureFlags Error:", flag, err)
	}
	log.Println(context)
	log.Println(flag, "FS for", env, "is", val)

}
