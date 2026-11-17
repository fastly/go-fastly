package botmanagement_test

import (
	"context"
	"testing"

	"github.com/fastly/go-fastly/v19/fastly"
	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/botmanagement"
	"github.com/fastly/go-fastly/v19/internal/serviceproductcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"

	"github.com/stretchr/testify/require"
)

var serviceID = fastly.TestDeliveryServiceID

func TestUpdateConfigurationMissingContentGuard(t *testing.T) {
	t.Parallel()

	_, err := botmanagement.UpdateConfiguration(context.TODO(), nil, serviceID, botmanagement.ConfigureInput{ContentGuard: ""})

	require.ErrorIs(t, err, botmanagement.ErrMissingContentGuard)
}

var functionalTests = []*test_utils.FunctionalTestService{
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          botmanagement.Disable,
		IgnoreFailure: true,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          botmanagement.Get,
		ProductID:     botmanagement.ProductID,
		ExpectFailure: true,
	}),
	serviceproductcore.NewEnableTest(&serviceproductcore.EnableTestInput[botmanagement.EnableOutput, products.ServiceNullInput]{
		OpNoInputFn: botmanagement.Enable,
		ProductID:   botmanagement.ProductID,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      botmanagement.Get,
		ProductID: botmanagement.ProductID,
	}),
	serviceproductcore.NewGetConfigurationTest(&serviceproductcore.GetConfigurationTestInput[botmanagement.ConfigureOutput]{
		Phase:     "default",
		OpFn:      botmanagement.GetConfiguration,
		ProductID: botmanagement.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, o botmanagement.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.ContentGuard, "test '%s'", tc.Name)
			require.Equalf(t, "off", *o.Configuration.ContentGuard, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewUpdateConfigurationTest(&serviceproductcore.UpdateConfigurationTestInput[botmanagement.ConfigureOutput, botmanagement.ConfigureInput]{
		OpFn:      botmanagement.UpdateConfiguration,
		Input:     botmanagement.ConfigureInput{ContentGuard: "on"},
		ProductID: botmanagement.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, o botmanagement.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.ContentGuard, "test '%s'", tc.Name)
			require.Equalf(t, "on", *o.Configuration.ContentGuard, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewGetConfigurationTest(&serviceproductcore.GetConfigurationTestInput[botmanagement.ConfigureOutput]{
		Phase:     "after update",
		OpFn:      botmanagement.GetConfiguration,
		ProductID: botmanagement.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTestService, o botmanagement.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.ContentGuard, "test '%s'", tc.Name)
			require.Equalf(t, "on", *o.Configuration.ContentGuard, "test '%s'", tc.Name)
		},
	}),
	serviceproductcore.NewDisableTest(&serviceproductcore.DisableTestInput{
		OpFn: botmanagement.Disable,
	}),
	serviceproductcore.NewGetTest(&serviceproductcore.GetTestInput[botmanagement.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          botmanagement.Get,
		ProductID:     botmanagement.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementAndConfigurationDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTestsService(t, functionalTests, fastly.TestDeliveryServiceID)
}
