//go:generate rm -f api.go api_test.go
//go:generate service_linked_product

package fanout

const (
	ProductName = "Fanout"
	ProductID   = "fanout"
)
