package parser

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/tj10200/docker-compose-orchestrator/pkg/orchestrator/graph"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
)

type ComposeParser struct {
	Fs    afero.Fs
	Files []string
}

func (p *ComposeParser) ParseAsGraph() (gres *graph.Graph, err error) {
	gres = &graph.Graph{}
	for _, path := range p.Files {
		var file fs.File
		file, err = p.Fs.Open(path)
		if err != nil {
			return gres, fmt.Errorf("cannot open compose file(%s): %v", path, err)
		}

		var data []byte
		data, err = io.ReadAll(file)
		if err != nil {
			return gres, fmt.Errorf("cannot read file(%s): err(%v)", path, err)
		}

		cf := ComposeFile{}
		if err = yaml.Unmarshal(data, &cf); err != nil {
			return gres, fmt.Errorf("cannot unmarshal compose file(%s): err(%v)", path, err)
		}

		var g graph.Graph
		if g, err = cf.AsGraph(path); err != nil {
			return gres, fmt.Errorf("cannot create graph from compose file(%s): err(%v)", path, err)
		}

		if gres, err = gres.Append(&g); err != nil {
			return gres, fmt.Errorf("cannot append onto graph: %v", err)
		}
	}

	return
}
