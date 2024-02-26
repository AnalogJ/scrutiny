package errors_test

import (
	"testing"

	"github.com/analogj/scrutiny/webapp/backend/pkg/errors"
	"github.com/stretchr/testify/require"
)

//func TestCheckErr_WithoutError(t *testing.T) {
//	t.Parallel()
//
//	//assert
//	require.NotPanics(t, func() {
//		errors.CheckErr(nil)
//	})
//}

//func TestCheckErr_Error(t *testing.T) {
//	t.Parallel()
//
//	//assert
//	require.Panics(t, func() {
//		errors.CheckErr(stderrors.New("This is an error"))
//	})
//}

func TestErrors(t *testing.T) {
	t.Parallel()

	// assert
	require.Implements(t, (*error)(nil), errors.ConfigFileMissingError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.ConfigValidationError("test"), "should implement the error interface")
	require.Implements(t, (*error)(nil), errors.DependencyMissingError("test"), "should implement the error interface")
}
