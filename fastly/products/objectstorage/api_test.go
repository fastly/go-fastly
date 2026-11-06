package objectstorage_test

import (
	"testing"

	"github.com/fastly/go-fastly/v19/fastly/products"
	"github.com/fastly/go-fastly/v19/fastly/products/objectstorage"
	"github.com/fastly/go-fastly/v19/internal/productcore"
	"github.com/fastly/go-fastly/v19/internal/test_utils"
)

var functionalTests = []*test_utils.FunctionalTest{
	productcore.NewDisableTest(&productcore.DisableTestInput{
		Phase:         "ensure disabled before testing",
		OpFn:          objectstorage.Disable,
		IgnoreFailure: true,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[objectstorage.EnableOutput]{
		Phase:         "before enablement",
		OpFn:          objectstorage.Get,
		ProductID:     objectstorage.ProductID,
		ExpectFailure: true,
	}),
	productcore.NewEnableTest(&productcore.EnableTestInput[objectstorage.EnableOutput, products.NullInput]{
		OpNoInputFn: objectstorage.Enable,
		ProductID:   objectstorage.ProductID,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[objectstorage.EnableOutput]{
		Phase:     "after enablement",
		OpFn:      objectstorage.Get,
		ProductID: objectstorage.ProductID,
	}),
	productcore.NewDisableTest(&productcore.DisableTestInput{
		OpFn: objectstorage.Disable,
	}),
	productcore.NewGetTest(&productcore.GetTestInput[objectstorage.EnableOutput]{
		Phase:         "after disablement",
		OpFn:          objectstorage.Get,
		ProductID:     objectstorage.ProductID,
		ExpectFailure: true,
	}),
}

func TestEnablement(t *testing.T) {
	test_utils.ExecuteFunctionalTests(t, functionalTests)
}
