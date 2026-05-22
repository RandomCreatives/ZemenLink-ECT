package modules

import (
	"encoding/json"
	"os"
)

type Manifest struct {
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Enabled     bool            `json:"enabled"`
	Config      json.RawMessage `json:"config"`
	Permissions []string        `json:"permissions"`
}

func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// LoadModules walks the modules directory and loads all manifests
func LoadModules(basePath string) ([]*Manifest, error) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	var manifests []*Manifest
	for _, entry := range entries {
		if entry.IsDir() {
			manifestPath := basePath + "/" + entry.Name() + "/manifest.json"
			m, err := LoadManifest(manifestPath)
			if err == nil {
				manifests = append(manifests, m)
			}
		}
	}
	return manifests, nil
}

// Registry stores loaded modules and provides feature flag resolution
type Registry struct {
	modules []*Manifest
}

func NewRegistry(modules []*Manifest) *Registry {
	return &Registry{modules: modules}
}

func (r *Registry) GetFlagsForTenant(tenantID string, complianceLevel string) map[string]bool {
	flags := make(map[string]bool)
	for _, m := range r.modules {
		if m.Enabled {
			flags[m.Name] = true
		}
	}
	return flags
}
