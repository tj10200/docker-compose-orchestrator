package orchestrator

import (
	log "github.com/sirupsen/logrus"
)

type Orchestrator struct {
	ComposeFiles []string
	ConfigFile   string
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

func NewOrchestrator() (*Orchestrator, error) {
	return &Orchestrator{}, nil
}

func (o *Orchestrator) Run() (err error) {
	return err
}
