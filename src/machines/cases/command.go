package cases

// Interface for command, implementings business use cases.
type Command interface {
	// Execute business cases.
	Execute() error
}
