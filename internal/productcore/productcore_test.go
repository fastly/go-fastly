package productcore_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/internal/productcore"
	"github.com/stretchr/testify/require"
)

func TestDeleteMissingServiceID(t *testing.T) {
	t.Parallel()

	err := productcore.Delete(&productcore.DeleteInput{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestGetMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := productcore.Get[*products.NullOutput](&productcore.GetInput{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestPatchMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := productcore.Patch[*products.NullOutput](&productcore.PatchInput[*products.NullInput]{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestPutMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := productcore.Put[*products.NullOutput](&productcore.PutInput[*products.NullInput]{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}
