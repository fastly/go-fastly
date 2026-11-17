package serviceproductcore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
)

func TestDeleteMissingServiceID(t *testing.T) {
	t.Parallel()

	err := serviceproductcore.Delete(context.TODO(), &serviceproductcore.DeleteInput{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestGetMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := serviceproductcore.Get[*products.ServiceNullOutput](context.TODO(), &serviceproductcore.GetInput{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestPatchMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := serviceproductcore.Patch[*products.ServiceNullOutput](context.TODO(), &serviceproductcore.PatchInput[*products.ServiceNullInput]{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}

func TestPutMissingServiceID(t *testing.T) {
	t.Parallel()

	_, err := serviceproductcore.Put[*products.ServiceNullOutput](context.TODO(), &serviceproductcore.PutInput[*products.ServiceNullInput]{
		ServiceID: "",
	})

	require.ErrorIs(t, err, fastly.ErrMissingServiceID)
}
