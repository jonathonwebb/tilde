package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/spf13/pflag"
)

type ExitStatus int

const (
	ExitSuccess    ExitStatus = 0
	ExitFailure    ExitStatus = 1
	ExitUsageError ExitStatus = 2
)

var (
	ErrMissingCommand = errors.New("missing command")
	ErrUnknownCommand = errors.New("missing command")
)

type Command struct {
	Name, Description, Usage, Help string
	Error                          func(env *Env, err error) ExitStatus
	Action                         func(context.Context, *Env) ExitStatus

	configFlagName string
	flags          *pflag.FlagSet
	envMap         map[string]string
	configMap      map[string]string
	commands       []*Command
}

func (c *Command) AddCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

func (c *Command) ensureFlagSet() {
	if c.flags == nil {
		c.flags = pflag.NewFlagSet(c.Name, pflag.ContinueOnError)
	}
}

func (c *Command) initFlagSet() {
	c.ensureFlagSet()
	c.flags.Usage = func() {}
	c.flags.SetOutput(io.Discard)
	c.flags.SetInterspersed(false)
	c.flags.SortFlags = false
}

func (c *Command) SetFlags(fn func(*pflag.FlagSet)) {
	c.ensureFlagSet()
	fn(c.flags)
	c.initFlagSet()
}

func (c *Command) setEnvMapping(key, value string) {
	if c.envMap == nil {
		c.envMap = make(map[string]string, 1)
	}
	c.envMap[key] = value
}

func (c *Command) lookupEnvMapping(key string) (string, bool) {
	if c.envMap == nil {
		return "", false
	}
	val, found := c.envMap[key]
	return val, found
}

func (c *Command) MapEnvVar(name string, envvar string) {
	f := c.flags.Lookup(name)
	if f != nil {
		c.setEnvMapping(name, envvar)
	}
}

func (c *Command) setConfigMapping(key, value string) {
	if c.configMap == nil {
		c.configMap = make(map[string]string, 1)
	}
	c.configMap[key] = value
}

func (c *Command) lookupConfigMapping(key string) (string, bool) {
	if c.configMap == nil {
		return "", false
	}
	val, found := c.configMap[key]
	return val, found
}

func (c *Command) MapConfigKey(name string, configkey string) {
	f := c.flags.Lookup(name)
	if f != nil {
		c.setConfigMapping(name, configkey)
	}
}

func (c *Command) SetConfigFlag(name string) {
	if f := c.flags.Lookup(name); f != nil {
		c.configFlagName = name
	}
}

func (c *Command) error(env *Env, err error) ExitStatus {
	if c.Error != nil {
		return c.Error(env, err)
	} else {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(env.Stdout, "%s\n", c.Help)
			return ExitSuccess
		}
		if errors.Is(err, ErrUnknownCommand) {
			fmt.Fprintf(env.Stderr, "unknown command %s\n", env.Args[0])
		} else if errors.Is(err, ErrMissingCommand) {
			fmt.Fprintf(env.Stderr, "missing command\n")
		} else {
			fmt.Fprintf(env.Stderr, "%v\n", err)
		}
		fmt.Fprintf(env.Stderr, "%s\n", c.Usage)
		return ExitUsageError
	}
}

func (c *Command) Execute(ctx context.Context, env *Env) ExitStatus {
	c.initFlagSet()

	help := c.flags.BoolP("help", "h", false, "show this help and exit")

	if err := c.flags.Parse(env.Args[1:]); err != nil {
		return c.error(env, err)
	}

	if *help {
		return c.error(env, flag.ErrHelp)
	}

	if c.configFlagName != "" {
		var file string
		if configFlag := c.flags.Lookup(c.configFlagName); configFlag != nil {
			file = configFlag.Value.String()
			if file != "" {
				if err := env.initConfig(file); err != nil {
					return c.error(env, err)
				}
			}
		}
	}

	var flagErr error
	c.flags.VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			envvar, found := c.lookupEnvMapping(f.Name)
			if found {
				if val, found := env.lookupVar(envvar); found {
					err := f.Value.Set(val)
					if err != nil {
						if flagErr != nil {
							flagErr = err
						}
					}
					return
				}
			}

			configkey, found := c.lookupConfigMapping(f.Name)
			if found {
				if val, found := env.lookupConfig(configkey); found {
					err := f.Value.Set(val)
					if err != nil {
						if flagErr != nil {
							flagErr = err
						}
					}
					return
				}
			}
		}
	})
	if flagErr != nil {
		return c.error(env, flagErr)
	}

	env.Args = c.flags.Args()

	if c.Action != nil {
		return c.Action(ctx, env)
	}
	if len(env.Args) == 0 {
		return c.error(env, ErrMissingCommand)
	}
	for _, cmd := range c.commands {
		if cmd.Name == env.Args[0] {
			return cmd.Execute(ctx, env)
		}
	}
	return c.error(env, ErrUnknownCommand)
}
