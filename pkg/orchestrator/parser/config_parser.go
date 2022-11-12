package parser

import (
	"fmt"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
)

type ConfigParser struct {
	Fs   afero.Fs
	File string
}

func (p *ConfigParser) Parse() (cf ConfigFile, err error) {
	var file fs.File
	file, err = p.Fs.Open(p.File)
	if err != nil {
		return cf, fmt.Errorf("cannot open compose file(%s): %v", p.File, err)
	}

	var data []byte
	data, err = io.ReadAll(file)
	if err != nil {
		return cf, fmt.Errorf("cannot read file(%s): err(%v)", p.File, err)
	}

	if err = yaml.Unmarshal(data, &cf); err != nil {
		return cf, fmt.Errorf("cannot unmarshal config file(%s): err(%v)", p.File, err)
	}

	for k, v := range cf.Services {
		for i, c := range v.Commands {
			newEnv := make(map[string]string, len(c.Env))
			for ek, ev := range c.Env {
				if ev[0] == '$' {
					ev = os.ExpandEnv(ev)
				}
				newEnv[ek] = ev
			}
			c.Env = newEnv
			v.Commands[i] = c
		}
		for i, c := range v.AfterRun {
			newEnv := make(map[string]string, len(c.Env))
			for ek, ev := range c.Env {
				if ev[0] == '$' {
					ev = os.ExpandEnv(ev)
				}
				newEnv[ek] = ev
			}
			c.Env = newEnv
			v.AfterRun[i] = c
		}
		cf.Services[k] = v
	}
	return
}
