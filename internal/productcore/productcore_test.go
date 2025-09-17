package productcore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v12/fastly"
	"github.com/fastly/go-fastly/v12/fastly/products"
	"github.com/fastly/go-fastly/v12/internal/productcore"
)

func TestDeleteMissingServiceID(t *testing.T) {
	t.Parallel()

	err := productcore.Delete(context.TODO(), &productcore.DeleteInput{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestGetMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := productcore.Get[*products.NullOutput](context.TODO(), &productcore.GetInput{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestPatchMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := productcore.Patch[*products.NullOutput](context.TODO(), &productcore.PatchInput[*products.NullInput]{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestPutMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := productcore.Put[*products.NullOutput](context.TODO(), &productcore.PutInput[*products.NullInput]{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}
