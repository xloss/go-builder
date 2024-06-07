package builder

import "errors"

var (
	// UpdateNoSets means that when creating update query no set values
	UpdateNoSets = errors.New("no sets")
)
