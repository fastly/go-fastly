package ddosprotection_test

import (
	"context"
	"testing"

	"github.com/fastly/go-fastly/v12/fastly"
	"github.com/fastly/go-fastly/v12/fastly/products"
	"github.com/fastly/go-fastly/v12/fastly/products/ddosprotection"
	"github.com/fastly/go-fastly/v12/internal/productcore"
	"github.com/fastly/go-fastly/v12/internal/test_utils"

	"github.com/stretchr/testify/require"
)

var serviceID = fastly.TestDeliveryServiceID

func TestUpdateConfigurationMissingMode(t *testing.T) {
	t.Parallel()

	_, err := ddosprotection.UpdateConfiguration(context.TODO(), nil, serviceID, ddosprotection.ConfigureInput{Mode: ""})

	require.ErrorIs(t, err, ddosprotection.ErrMissingMode)
}

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          ddosprotection.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[ddosprotection.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          ddosprotection.Get,
		ProductID:     ddosprotection.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[ddosprotection.EnableOutput, products.NullInput]{
		OpNoInputFn: ddosprotection.Enable,
		ProductID:   ddosprotection.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[ddosprotection.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      ddosprotection.Get,
		ProductID: ddosprotection.ProductID,
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[ddosprotection.ConfigureOutput]{
		Phase:     "default",
		OpFn:      ddosprotection.GetConfiguration,
		ProductID: ddosprotection.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o ddosprotection.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "log", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	productcore.NewUpdateConfigurationTest(&productcore.UpdateConfigurationTestInput[ddosprotection.ConfigureOutput, ddosprotection.ConfigureInput]{
		OpFn:      ddosprotection.UpdateConfiguration,
		Input:     ddosprotection.ConfigureInput{Mode: "block"},
		ProductID: ddosprotection.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o ddosprotection.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	productcore.NewGetConfigurationTest(&productcore.GetConfigurationTestInput[ddosprotection.ConfigureOutput]{
		Phase:     "after update",
		OpFn:      ddosprotection.GetConfiguration,
		ProductID: ddosprotection.ProductID,
		CheckOutputFn: func(t *testing.T, tc *test_utils.FunctionalTest, o ddosprotection.ConfigureOutput) {
			require.NotNilf(t, o.Configuration.Mode, "test '%s'", tc.Name)
			require.Equalf(t, "block", *o.Configuration.Mode, "test '%s'", tc.Name)
		},
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: ddosprotection.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[ddosprotection.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          ddosprotection.Get,
		ProductID:     ddosprotection.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablementAndConfigurationDelivery(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestDeliveryServiceID)
}

func TestEnablementAndConfigurationCompute(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests, fastly.TestComputeServiceID)
}
