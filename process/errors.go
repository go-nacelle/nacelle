package process

type (
	errMeta struct {
		err        error
		source     namedInitializer
		silentExit bool
	}

	errMetaSet []errMeta
)

func (set errMetaSet) Error() string {
	return "<multiple errors from process group>"
}

func coerceToSet(err error, source namedInitializer) errMetaSet {
	if set, ok := err.(errMetaSet); ok {
		return set
	}

	return errMetaSet{errMeta{err: err, source: source}}
}
