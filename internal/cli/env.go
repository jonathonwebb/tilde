package cli

import (
	"io"
	"os"
	"strings"
)

type Env struct {
	Build          string
	Stderr, Stdout io.Writer
	Args           []string
	Vars           map[string]string
	ConfigFile     map[string]string
	LoadConfigFile func(name string) (map[string]string, error)
}

func DefaultEnv(build string) *Env {
	env := make(map[string]string)
	for _, variable := range os.Environ() {
		parts := strings.SplitN(variable, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	return &Env{
		Build:          build,
		Stderr:         os.Stderr,
		Stdout:         os.Stdout,
		Args:           os.Args,
		Vars:           env,
		LoadConfigFile: DefaultConfigFileLoader,
	}
}

func DefaultConfigFileLoader(name string) (map[string]string, error) {
	// TODO: implement
	return map[string]string{}, nil
}

func (e *Env) lookupVar(key string) (val string, found bool) {
	if e.Vars != nil {
		val, found = e.Vars[key]
	}
	return
}

func (e *Env) lookupConfig(key string) (val string, found bool) {
	if e.ConfigFile != nil {
		val, found = e.ConfigFile[key]
	}
	return
}

func (e *Env) initConfig(name string) error {
	if e.LoadConfigFile != nil {
		config, err := e.LoadConfigFile(name)
		if err != nil {
			return err
		}
		e.ConfigFile = config
	}
	return nil
}
