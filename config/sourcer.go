package config

// Sourcer pulls requested names from a variable source. This can be the
// environment, a file, a remote server, etc. This can be done on-demand
// per variable, or a cache of variables can be built on startup and then
// pulled from a cached mapping as requested.
type Sourcer interface {
	// Tags returns a list of tags which are required to get a value from
	// the source. Order matters.
	Tags() []string

	// Get will retrieve a value from the source with the given tag values.
	// The tag values passed to this method will be in the same order as
	// returned from the Tags method. The boolean flags in the output reflect
	// if the requested tag values were empty (if this field may be skipped),
	// and whether or not a value was found.
	Get(values []string) (string, bool, bool)
}
