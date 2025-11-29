package model

import "time"

// Worldbook stores lightweight world-building snippets associated with a role.
type Worldbook struct {
	ID        string         `json:"id"`
	RoleID    string         `json:"role_id"`
	Data      map[string]any `json:"data"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type WorldSummary struct {
	Summary  string   `json:"summary"`
	Scene    string   `json:"scene"`
	Timeline string   `json:"timeline"`
	NPCs     []string `json:"npcs"`
	Entries  map[string][]string `json:"entries,omitempty"`
}

func (w *Worldbook) Summary() *WorldSummary {
	if w == nil {
		return nil
	}
	summary := &WorldSummary{Entries: map[string][]string{}}
	if w.Data != nil {
		if v, ok := w.Data["summary"].(string); ok {
			summary.Summary = v
		}
		if v, ok := w.Data["scene"].(string); ok {
			summary.Scene = v
		}
		if v, ok := w.Data["timeline"].(string); ok {
			summary.Timeline = v
		}
		if raw, ok := w.Data["npcs"]; ok {
			switch typed := raw.(type) {
			case []any:
				for _, item := range typed {
					if s, ok := item.(string); ok {
						summary.NPCs = append(summary.NPCs, s)
					}
				}
			case []string:
				summary.NPCs = append(summary.NPCs, typed...)
			}
		}
		if raw, ok := w.Data["entries"]; ok {
			switch typed := raw.(type) {
			case map[string]any:
				for k, v := range typed {
					var vals []string
					switch vv := v.(type) {
					case []any:
						for _, item := range vv {
							if s, ok := item.(string); ok {
								vals = append(vals, s)
							}
						}
					case []string:
						vals = append(vals, vv...)
					case string:
						vals = append(vals, vv)
					}
					if len(vals) > 0 {
						summary.Entries[k] = vals
					}
				}
			}
		}
	}
	return summary
}
