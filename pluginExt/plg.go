// Package plg provides extensions and utilities for the standard plugin library.
// It aims to simplify plugin loading, symbol lookup, and error handling when
// working with Go plugins.
package pluginExt

import (
	"fmt"
	"plugin"
	"reflect"
	"sync"
)

// PluginCache provides a thread-safe cache for loaded plugins
type PluginCache struct {
	mu      sync.RWMutex
	plugins map[string]*plugin.Plugin
}

// NewPluginCache creates a new plugin cache
func NewPluginCache() *PluginCache {
	return &PluginCache{
		plugins: make(map[string]*plugin.Plugin),
	}
}

// Load loads a plugin from the given path, caching the result
func (pc *PluginCache) Load(path string) (*plugin.Plugin, error) {
	pc.mu.RLock()
	if p, ok := pc.plugins[path]; ok {
		pc.mu.RUnlock()
		return p, nil
	}
	pc.mu.RUnlock()

	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin %s: %w", path, err)
	}

	pc.mu.Lock()
	pc.plugins[path] = p
	pc.mu.Unlock()

	return p, nil
}

// LoadOrError loads a plugin and returns an error if it cannot be loaded
func LoadOrError(path string) (*plugin.Plugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin %s: %w", path, err)
	}
	return p, nil
}

// LookupSymbol looks up a symbol in a plugin with proper error handling
func LookupSymbol(p *plugin.Plugin, symbolName string) (plugin.Symbol, error) {
	sym, err := p.Lookup(symbolName)
	if err != nil {
		return nil, fmt.Errorf("symbol %q not found in plugin: %w", symbolName, err)
	}
	return sym, nil
}

// LookupFunc looks up a symbol and ensures it's a function matching the expected type
func LookupFunc(p *plugin.Plugin, symbolName string, expectedType interface{}) (reflect.Value, error) {
	sym, err := LookupSymbol(p, symbolName)
	if err != nil {
		return reflect.Value{}, err
	}

	expectedTypeVal := reflect.TypeOf(expectedType)
	actualTypeVal := reflect.TypeOf(sym)

	if expectedTypeVal != actualTypeVal {
		return reflect.Value{}, fmt.Errorf("symbol %q has type %v, expected %v",
			symbolName, actualTypeVal, expectedTypeVal)
	}

	return reflect.ValueOf(sym), nil
}

// LoadAndLookup combines loading a plugin and looking up a symbol
func LoadAndLookup(path, symbolName string) (plugin.Symbol, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin %s: %w", path, err)
	}

	return LookupSymbol(p, symbolName)
}

// LoadAll loads multiple plugins from the given paths
func LoadAll(paths []string) (map[string]*plugin.Plugin, error) {
	result := make(map[string]*plugin.Plugin, len(paths))
	errors := []error{}

	for _, path := range paths {
		p, err := plugin.Open(path)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to load plugin %s: %w", path, err))
			continue
		}
		result[path] = p
	}

	if len(errors) > 0 {
		return result, fmt.Errorf("failed to load some plugins: %v", errors)
	}

	return result, nil
}
