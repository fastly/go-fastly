package contacts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v15/fastly"
)

func TestClient_Contacts(t *testing.T) {
	ctx := context.TODO()

	var (
		err        error
		customerID string
	)
	fastly.Record(t, "get_current_user", func(c *fastly.Client) {
		var u *fastly.User
		u, err = c.GetCurrentUser(ctx)
		if err == nil && u != nil && u.CustomerID != nil {
			customerID = *u.CustomerID
		}
	})
	require.NoError(t, err)
	require.NotEmpty(t, customerID)

	// Fastly refuses to delete the last contact of a given type, so we
	// create a guard contact first that we leave in place for the duration
	// of the test (and best-effort delete during cleanup).
	//
	// NOTE: When recreating the fixtures, update both emails.
	guardEmail := "go-fastly-test+contact-guard+20260522@example.com"
	email := "go-fastly-test+contact+20260522@example.com"

	var guard *Contact
	fastly.Record(t, "create_guard", func(c *fastly.Client) {
		guard, err = Create(ctx, c, &CreateInput{
			CustomerID:  fastly.ToPointer(customerID),
			ContactType: fastly.ToPointer("emergency"),
			Name:        fastly.ToPointer("guard contact"),
			Email:       fastly.ToPointer(guardEmail),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, guard)

	var co *Contact
	fastly.Record(t, "create", func(c *fastly.Client) {
		co, err = Create(ctx, c, &CreateInput{
			CustomerID:  fastly.ToPointer(customerID),
			ContactType: fastly.ToPointer("emergency"),
			Name:        fastly.ToPointer("test contact"),
			Email:       fastly.ToPointer(email),
		})
	})
	require.NoError(t, err)
	require.NotEmpty(t, co.ContactID)
	require.Equal(t, email, co.Email)

	// Best-effort cleanup: guard deletion may fail if it is the last
	// emergency contact on the account, which is fine.
	defer func() {
		fastly.Record(t, "cleanup", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				CustomerID: fastly.ToPointer(customerID),
				ContactID:  fastly.ToPointer(co.ContactID),
			})
			_ = Delete(ctx, c, &DeleteInput{
				CustomerID: fastly.ToPointer(customerID),
				ContactID:  fastly.ToPointer(guard.ContactID),
			})
		})
	}()

	var cs []*Contact
	fastly.Record(t, "list", func(c *fastly.Client) {
		cs, err = List(ctx, c, &ListInput{
			CustomerID: fastly.ToPointer(customerID),
		})
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(cs), 1)

	fastly.Record(t, "delete", func(c *fastly.Client) {
		err = Delete(ctx, c, &DeleteInput{
			CustomerID: fastly.ToPointer(customerID),
			ContactID:  fastly.ToPointer(co.ContactID),
		})
	})
	require.NoError(t, err)
}
