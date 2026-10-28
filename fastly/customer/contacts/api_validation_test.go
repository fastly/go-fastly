package contacts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v15/fastly"
)

func TestClient_List_validation(t *testing.T) {
	_, err := List(context.TODO(), fastly.TestClient, &ListInput{
		CustomerID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingCustomerID)
}

func TestClient_Create_validation(t *testing.T) {
	_, err := Create(context.TODO(), fastly.TestClient, &CreateInput{
		CustomerID: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingCustomerID)
}

func TestClient_Delete_validation(t *testing.T) {
	err := Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		CustomerID: nil,
		ContactID:  fastly.ToPointer("abc"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingCustomerID)

	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		CustomerID: fastly.ToPointer("abc"),
		ContactID:  nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingContactID)
}
