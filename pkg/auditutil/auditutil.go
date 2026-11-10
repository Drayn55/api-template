package auditutil

import "encoding/json"

// ToJSON safely converts any struct to JSON string pointer
func ToJSON(v any) *string {
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	s := string(b)
	return &s
}

