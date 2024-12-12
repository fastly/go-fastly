package products_test

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	p "github.com/fastly/go-fastly/v9/internal/products"
	"github.com/stretchr/testify/require"
)

func TestDeleteMissingServiceID(t *testing.T) {
	t.Parallel()

	err := p.Delete(&p.DeleteInput{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestGetMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := p.Get(&p.GetInput[p.NullOutput]{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestPatchMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := p.Patch(&p.PatchInput[p.NullInput, p.NullOutput]{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestPutMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := p.Put(&p.PutInput[p.NullInput, p.NullOutput]{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}
