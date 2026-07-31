package input

var (
	Actions map[string]*Binding = map[string]*Binding{
		"MoveForward":  &Binding{},
		"MoveBackward": &Binding{},
		"MoveLeft":     &Binding{},
		"MoveRight":    &Binding{},
		"MoveUp":       &Binding{},
		"MoveDown":     &Binding{},
		"MoveFast":     &Binding{},
	}
)

// Retrieve default keys for each key scope
func DefaultBindings(scope string) map[string]string {
	switch scope {
	case "Character":
		return map[string]string{
			"V":           "ToggleViewMode",
			"LeftAlt":     "ToggleMouseCapture",
			"D":           "MoveRight",
			"A":           "MoveLeft",
			"W":           "MoveForward",
			"S":           "MoveBackward",
			"Space":       "MoveUp",
			"LeftControl": "MoveDown",
			"LeftShift":   "MoveFast",
		}
	default:
		return map[string]string{}
	}
}
