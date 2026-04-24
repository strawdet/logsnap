package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SchemaInfo is a lightweight summary of a schema for listing.
type SchemaInfo struct {
	Name           string   `json:"name"`
	Description    string   `json:"description,omitempty"`
	RequiredFields []string `json:"required_fields"`
	AllowedLevels  []string `json:"allowed_levels,omitempty"`
}

// ListSchemas returns summaries of all saved schemas, sorted by name.
func ListSchemas(dir string) ([]SchemaInfo, error) {
	schemasDir := filepath.Join(dir, "schemas")
	entries, err := os.ReadDir(schemasDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []SchemaInfo{}, nil
		}
		return nil, fmt.Errorf("read schemas dir: %w", err)
	}

	var infos []SchemaInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(schemasDir, e.Name()))
		if err != nil {
			continue
		}
		var s Schema
		if err := json.Unmarshal(data, &s); err != nil {
			continue
		}
		infos = append(infos, SchemaInfo{
			Name:           s.Name,
			Description:    s.Description,
			RequiredFields: s.RequiredFields,
			AllowedLevels:  s.AllowedLevels,
		})
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})
	return infos, nil
}
