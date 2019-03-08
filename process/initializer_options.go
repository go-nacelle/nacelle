package process

import "time"

// InitializerConfigFunc is a function used to append additional
// metadata to an initializer during registration.
type InitializerConfigFunc func(*InitializerMeta)

// WithInitializerName assigns a name to an initializer, visible in logs.
func WithInitializerName(name string) InitializerConfigFunc {
	return func(meta *InitializerMeta) { meta.name = name }
}

// WithInitializerTimeout sets the time limit for the initializer.
func WithInitializerTimeout(timeout time.Duration) InitializerConfigFunc {
	return func(meta *InitializerMeta) { meta.initTimeout = timeout }
}

// WithFinalizerTimeout sets the time limit for the finalizer.
func WithFinalizerTimeout(timeout time.Duration) InitializerConfigFunc {
	return func(meta *InitializerMeta) { meta.finalizeTimeout = timeout }
}
