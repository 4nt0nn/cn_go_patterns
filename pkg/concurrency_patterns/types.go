package pkg_cnp_types

// Future represents the type that the InnerFuture will satisfy.
type Future interface {
	Result() (string, error)
}
