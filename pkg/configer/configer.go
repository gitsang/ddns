package configer

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Options struct {
	template any

	envBind   bool
	envPrefix string
	envDelim  string

	command    *cobra.Command
	flagBind   bool
	flagPrefix string
	flagDelim  string
}

type Configer struct {
	Options
	viper *viper.Viper
}

type OptionFunc func(configer *Configer)

func WithTemplate(template any) OptionFunc {
	return func(c *Configer) {
		c.template = template
	}
}

func WithEnvBind(optfs ...OptionFunc) OptionFunc {
	return func(c *Configer) {
		c.envBind = true
		for _, apply := range optfs {
			apply(c)
		}
	}
}

func WithEnvPrefix(prefix string) OptionFunc {
	return func(c *Configer) {
		c.viper.SetEnvPrefix(prefix)
	}
}

func WithEnvDelim(delim string) OptionFunc {
	return func(c *Configer) {
		c.viper.SetEnvKeyReplacer(strings.NewReplacer(".", delim))
	}
}

func WithFlagBind(optfs ...OptionFunc) OptionFunc {
	return func(c *Configer) {
		c.flagBind = true
		for _, apply := range optfs {
			apply(c)
		}
	}
}

func WithCommand(command *cobra.Command) OptionFunc {
	return func(c *Configer) {
		c.command = command
	}
}

func WithFlagPrefix(prefix string) OptionFunc {
	return func(c *Configer) {
		c.flagPrefix = prefix
	}
}

func WithFlagDelim(delim string) OptionFunc {
	return func(c *Configer) {
		c.flagDelim = delim
	}
}

func (p *Configer) parseFlags(i interface{}, parents []string) {
	r := reflect.TypeOf(i)

	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		namespaces := append(parents, strings.ToLower(f.Name))

		if f.Type.Kind() == reflect.Struct {
			t := reflect.New(f.Type).Elem().Interface()
			p.parseFlags(t, namespaces)
			continue
		}

		// trim delim prefix to avoid empty prefix
		flagName := strings.TrimPrefix(strings.Join(namespaces, p.flagDelim), p.flagDelim)
		if flagTag := f.Tag.Get("flag"); flagTag != "" {
			flagName = flagTag
		}

		if f.Type.Kind() == reflect.Slice {
			p.command.Flags().StringSlice(flagName, nil, f.Tag.Get("usage"))
			continue
		}
		p.command.Flags().String(flagName, f.Tag.Get("default"), f.Tag.Get("usage"))

		// viperKey use dot to addressing (mapstructure default) and should exclude the prefix
		viperKey := strings.Join(namespaces[1:], ".")
		err := p.viper.BindPFlag(viperKey, p.command.Flags().Lookup(flagName))
		if err != nil {
			continue
		}
	}
}

func (p *Configer) parseEnv(i interface{}, parents []string) {
	t := reflect.TypeOf(i)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		namespaces := append(parents, strings.ToLower(f.Name))

		ft := f.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct {
			fi := reflect.New(ft).Elem().Interface()
			p.parseEnv(fi, namespaces)
			continue
		}

		viperKey := strings.Join(namespaces[1:], ".")
		if envTag := f.Tag.Get("env"); envTag != "" {
			if err := p.viper.BindEnv(viperKey, envTag); err != nil {
				continue
			}
		}
	}
}

func New(optfs ...OptionFunc) *Configer {
	c := &Configer{
		viper: viper.New(),
	}
	for _, apply := range optfs {
		apply(c)
	}

	if c.envBind {
		c.viper.AutomaticEnv()
		c.parseEnv(c.template, []string{c.envPrefix})
	}

	if c.flagBind {
		c.parseFlags(c.template, []string{c.flagPrefix})
	}

	return c
}

func (c *Configer) Load(config any, files ...string) error {
	for _, file := range files {
		c.viper.SetConfigFile(file)
		if err := c.viper.ReadInConfig(); err != nil {
			fmt.Printf("falied to read file %s: %s\n", file, err.Error())
			continue
		}
	}
	defaults.SetDefaults(config)
	return c.viper.Unmarshal(config)
}

func (c *Configer) Store(config any, file string) error {
	configYamlBytes, _ := yaml.Marshal(config)

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(configYamlBytes)
	if err != nil {
		return err
	}

	return nil
}
