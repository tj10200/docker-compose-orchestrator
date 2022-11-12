package orchestrator

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/tj10200/docker-compose-orchestrator/pkg/orchestrator/parser"
	"os"
	"os/exec"
	"path/filepath"
)

type Orchestrator struct {
	Fs           afero.Fs
	ComposeFiles []string
	ConfigFile   string
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

func NewOrchestrator(fs afero.Fs, composeFiles []string, configFile string) (*Orchestrator, error) {
	return &Orchestrator{
		Fs:           fs,
		ComposeFiles: composeFiles,
		ConfigFile:   configFile,
	}, nil
}

func (o *Orchestrator) Run() (err error) {
	if o.Fs == nil {
		o.Fs = afero.NewOsFs()
	}
	depGraph, err := (&parser.ComposeParser{
		Fs:    o.Fs,
		Files: o.ComposeFiles,
	}).ParseAsGraph()
	if err != nil {
		return fmt.Errorf("cannot parase compose files: %v", err)
	}

	cfg, err := (&parser.ConfigParser{
		Fs:   o.Fs,
		File: o.ConfigFile,
	}).Parse()
	if err != nil {
		return fmt.Errorf("cannot parse configuration file: %v", err)
	}

	order := depGraph.TopologySort()

	if err = o.runDownCmd(); err != nil {
		return err
	}

	for _, node := range order {
		cmd, has_cmd := cfg.Services[node.Name]
		if has_cmd {
			for _, c := range cmd.Commands {
				if err := runCommand(cmd.Name, c); err != nil {
					return err
				}
			}
		}

		path := filepath.Dir(node.File)
		if err := runDockerComposeUp(path, node.Name); err != nil {
			return err
		}

		if has_cmd {
			for _, c := range cmd.AfterRun {
				if err := runCommand(cmd.Name, c); err != nil {
					return err
				}
			}
		}
	}

	return err
}

func runCommand(name string, c parser.Command) error {
	if err := backoff.Retry(func() error {
		if c.Type == "docker" {
			return runDockerCmdAndWait(c.Image, c.Cmd, c.NetName, c.Env)
		} else if c.Type == "host" {
			return runHostCmdAndWait(c.Dir, c.Tool, c.Args, c.Env)
		}
		return nil
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10)); err != nil {
		return fmt.Errorf("cannot run dependency command: name(%s) image(%s) cmd(%s) err(%v)",
			name, c.Image, c.Cmd, err)
	}
	return nil
}

func runDockerCmdAndWait(image string, cmd string, net string, env map[string]string) error {
	cmdArgs := []string{
		"run",
		"--network", net,
	}

	for k, v := range env {
		cmdArgs = append(cmdArgs, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	cmdArgs = append(cmdArgs, "--entrypoint", "/bin/bash")
	cmdArgs = append(cmdArgs, image)
	cmdArgs = append(cmdArgs, "-c", cmd)
	cmdArgs = replaceArgVars(cmdArgs, env)
	out, err := exec.Command("docker", cmdArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot execute docker run %s %s: %v", image, cmd, err)
	} else {
		log.Infof(string(out))
	}
	return nil
}

func runHostCmdAndWait(dir, toolName string, args []string, env map[string]string) error {
	cmd := exec.Command(toolName, args...)
	cmd.Dir = dir
	cmd.Args = replaceArgVars(cmd.Args, env)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot execute host cmd %s %s %s: %v", dir, toolName, args, err)
	} else {
		log.Infof(string(out))
	}
	return nil
}

func (o *Orchestrator) runDownCmd() error {
	for _, file := range o.ComposeFiles {
		path := filepath.Dir(file)
		out, err := exec.Command(
			"cd", path, "&&", "docker", "compose", "down").CombinedOutput()
		if err != nil {
			return fmt.Errorf("cannot execute docker compose down. file(%s): err(%v)", file, err)
		} else {
			log.Infof(string(out))
		}
	}
	return nil
}

func runDockerComposeUp(dir, service string) error {
	c := exec.Command(
		"docker", "compose", "up", "-d", service)
	c.Dir = dir
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot execute docker compose up %s: output(%s) err(%v)", service, out, err)
	} else {
		log.Infof(string(out))
	}
	return nil
}

func replaceArgVars(args []string, env map[string]string) []string {
	for i, a := range args {
		os.Expand(a, func(k string) string {
			if v, ok := env[k]; ok {
				return v
			}
			return ""
		})
		args[i] = a
	}
	return args
}
