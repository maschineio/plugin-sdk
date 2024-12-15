package pluginsdk

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/samber/lo"
	"maschine.io/core/context"
)

var _ (ResourceManager) = (*manager)(nil)
var lock = &sync.Mutex{}

var managerInstance ResourceManager

type LambdaFn func(*context.Context) (any, error)

type ResourceManager interface {
	GetFn(resourceName string) LambdaFn
	RegisterLambdaFn(rn string, fn LambdaFn) error
	ResourceNames() []string
	LoadPlugins(pluginDir string) error
}

type manager struct {
	searchPathes []string
	lambdas      map[string]LambdaFn
	plugins      map[string]plugin.Plugin
}

func GetResourceManager() ResourceManager {
	if managerInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if managerInstance == nil {

			// only needed for external plugins
			searchPathes := lo.Uniq(filepath.SplitList(os.Getenv("PATH")))
			managerInstance = &manager{
				searchPathes: searchPathes,
				lambdas:      make(map[string]LambdaFn, 0),
				plugins:      make(map[string]plugin.Plugin, 0),
			}

		}
	}
	return managerInstance
}

func (m *manager) GetFn(resourceName string) (f LambdaFn) {
	if f, found := m.lambdas[resourceName]; found {
		return f
	}
	return nil
}

func (m *manager) RegisterLambdaFn(rn string, fn LambdaFn) error {
	if _, found := m.lambdas[rn]; found {
		return fmt.Errorf("lambda function already registered: %v", rn)
	}
	m.lambdas[rn] = fn
	return nil
}

func (m *manager) ResourceNames() (result []string) {
	result = make([]string, len(m.lambdas))
	i := 0
	for n := range m.lambdas {
		result[i] = n
		i++
	}
	sort.Strings(result)
	return
}

func (m *manager) LoadPlugins(pluginDir string) error {
	pluginFiles, err := filepath.Glob(filepath.Join(pluginDir, "*.so"))
	if err != nil {
		return err
	}

	for _, pluginFile := range pluginFiles {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  1,
				MagicCookieKey:   "PLUGIN_MAGIC_COOKIE",
				MagicCookieValue: "plugin",
			},
			Plugins: map[string]plugin.Plugin{
				// "lambda": &LambdaPluginImpl{},
			},
			Cmd:              exec.Command(pluginFile),
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		})

		rpcClient, err := client.Client()
		if err != nil {
			return err
		}

		raw, err := rpcClient.Dispense("lambda")
		if err != nil {
			return err
		}

		lambdaPlugin := raw.(LambdaPlugin)
		for name, fn := range lambdaPlugin.Functions() {
			m.RegisterLambdaFn(name, fn)
		}
	}

	return nil
}

type LambdaPlugin interface {
	Functions() map[string]LambdaFn
}

/*
type LambdaPluginImpl struct {
	plugin.Plugin
	Impl LambdaPlugin
}

func (p *LambdaPluginImpl) Functions() map[string]LambdaFn {
	return p.Impl.Functions()
}
*/
