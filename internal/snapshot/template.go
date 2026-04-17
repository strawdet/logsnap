package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Template struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels,omitempty"`
	Tags      []string          `json:"tags,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

func templatePath(dir, name string) string {
	return filepath.Join(dir, "template_"+name+".json")
}

func SaveTemplate(dir string, t *Template) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	t.CreatedAt = time.Now()
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal template: %w", err)
	}
	return os.WriteFile(templatePath(dir, t.Name), data, 0644)
}

func LoadTemplate(dir, name string) (*Template, error) {
	data, err := os.ReadFile(templatePath(dir, name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template %q not found", name)
		}
		return nil, fmt.Errorf("read template: %w", err)
	}
	var t Template
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("unmarshal template: %w", err)
	}
	return &t, nil
}

func DeleteTemplate(dir, name string) error {
	p := templatePath(dir, name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("template %q not found", name)
	}
	return os.Remove(p)
}

func ApplyTemplate(dir, name string, s *Snapshot) error {
	t, err := LoadTemplate(dir, name)
	if err != nil {
		return err
	}
	if s.Labels == nil {
		s.Labels = make(map[string]string)
	}
	for k, v := range t.Labels {
		if _, exists := s.Labels[k]; !exists {
			s.Labels[k] = v
		}
	}
	for _, tag := range t.Tags {
		s.Tags = appendIfMissing(s.Tags, tag)
	}
	return nil
}

func appendIfMissing(tags []string, tag string) []string {
	for _, t := range tags {
		if t == tag {
			return tags
		}
	}
	return append(tags, tag)
}
