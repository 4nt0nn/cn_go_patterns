package pkg_types

import "context"

type Circuit func(context.Context) (string, error)
