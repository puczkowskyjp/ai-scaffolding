package planning

import "errors"

// ErrInvalidSelection reports that UI agent-selection state does not match the plan.
var ErrInvalidSelection = errors.New("invalid agent selection")
