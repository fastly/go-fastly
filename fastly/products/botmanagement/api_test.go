package botmanagement_test

import (
	"context"
	"testing"

	"github.com/fastly/go-fastly/v14/fastly"
	"github.com/fastly/go-fastly/v14/fastly/products"
	"github.com/fastly/go-fastly/v14/fastly/products/botmanagement"
	"github.com/fastly/go-fastly/v14/internal/productcore"
	"github.com/fastly/go-fastly/v14/internal/test_utils"

	"github.com/stretchr/testify/require"
)

var serviceID = fastly.TestDeliveryServiceID

func TestUpdateConfigurationMissingContentGuard(t *testing.T) {
	t.Parallel()

	_, err := botmanagement.UpdateConfiguration(context.TODO(), nil, serviceID, botmanagement.ConfigureInput{ContentGuard: ""})

	require.ErrorIs(t, err, botmanagement.ErrMissingContentGuard)
}

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          botmanagement.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          botmanagement.Get,
		ProductID:     botmanagement.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[botmanagement.EnableOutput, products.NullInput]{
		OpNoInputFn: botmanagement.Enable,
		ProductID:   botmanagement.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      botmanagement.Get,
		ProductID: botmanagement.ProductID,
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[botmanagement.ConfigureOutput]{
		Phase:     "default",
		OpFn:      botmanagement.GetConfiguration,
		ProductID: botmanagement.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o botmanagement.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.ContentGuard, "test '%s'", tc.Name)
			require.Equalf(t, "off", *o.Configuration.ContentGuard, "test '%s'", tc.Name)
		},
	}),
	productcore.NewUpdateConfigurationTest(&productcore.UpdateConfigurationTestInput[botmanagement.ConfigureOutput, botmanagement.ConfigureInput]{
		OpFn:      botmanagement.UpdateConfiguration,
		Input:     botmanagement.ConfigureInput{ContentGuard: "on"},
		ProductID: botmanagement.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o botmanagement.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.ContentGuard, "test '%s'", tc.Name)
			require.Equalf(t, "on", *o.Configuration.ContentGuard, "test '%s'", tc.Name)
		},
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[botmanagement.ConfigureOutput]{
		Phase:     "after update",
		OpFn:      botmanagement.GetConfiguration,
		ProductID: botmanagement.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o botmanagement.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.ContentGuard, "test '%s'", tc.Name)
			require.Equalf(t, "on", *o.Configuration.ContentGuard, "test '%s'", tc.Name)
		},
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: botmanagement.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          botmanagement.Get,
		ProductID:     botmanagement.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementAndConfigurationDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}
