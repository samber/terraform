package terraform

import (
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/dot"
)

// GraphNodeConfigDataSource represents a data source within the config graph.
type GraphNodeConfigDataSource struct {
	DataSource *config.DataSource
	Path       []string
}

func (n *GraphNodeConfigDataSource) Copy() *GraphNodeConfigDataSource {
	ncds := &GraphNodeConfigDataSource{
		DataSource: n.DataSource.Copy(),
		Path:       make([]string, 0, len(n.Path)),
	}
	for _, p := range n.Path {
		ncds.Path = append(ncds.Path, p)
	}
	return ncds
}

func (n *GraphNodeConfigDataSource) ConfigType() GraphNodeConfigType {
	return GraphNodeConfigTypeDataSource
}

func (n *GraphNodeConfigDataSource) DependableName() []string {
	return []string{n.DataSource.Id()}
}

// GraphNodeDependent impl.
func (n *GraphNodeConfigDataSource) DependentOn() []string {
	result := make(
		[]string,
		len(n.DataSource.DependsOn),
		len(n.DataSource.RawConfig.Variables)+
			len(n.DataSource.DependsOn)*2,
	)

	copy(result, n.DataSource.DependsOn)

	for _, v := range n.DataSource.RawConfig.Variables {
		if vn := varNameForVar(v); vn != "" {
			result = append(result, vn)
		}
	}

	return result
}

// VarWalk calls a callback for all the variables that this data source
// depends on.
func (n *GraphNodeConfigDataSource) VarWalk(fn func(config.InterpolatedVariable)) {
	for _, v := range n.DataSource.RawConfig.Variables {
		fn(v)
	}
}

func (n *GraphNodeConfigDataSource) Name() string {
	return n.DataSource.Id()
}

// GraphNodeDotter impl.
func (n *GraphNodeConfigDataSource) DotNode(name string, opts *GraphDotOpts) *dot.Node {
	return dot.NewNode(name, map[string]string{
		"label": n.Name(),
		"shape": "cylinder",
	})
}

// GraphNodeProviderConsumer impl
func (n *GraphNodeConfigDataSource) ProvidedBy() []string {
	return []string{resourceProvider(n.DataSource.Type, n.DataSource.Provider)}
}
