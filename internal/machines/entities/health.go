package entities

// Health is an internal value type encapsulating operations with machine health state,
// and it is responsible for recalculation machine health depending of updates, missing
// for specific machine.
type health struct {
	level HealthLevel
}

// Recalculates health level depending on missing updates. Returns another example of
// enitity, so it should be used like immutable type.
func (h *health) Recalculate(mu []MissingUpdate) health {
	var newLevel HealthLevel = Healthy

	for _, missing := range mu {
		if missing.Severity == Critical || missing.Severity == Important {
			newLevel = Danger
			break
		}
		if missing.Severity == Low {
			newLevel = Warning
		}
	}

	return health{
		level: HealthLevel(newLevel),
	}
}
