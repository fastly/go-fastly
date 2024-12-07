//go:generate rm -f api.go
//go:generate service_linked_product -api

package fanout

const (
	ProductName = "Fanout"
	ProductID   = "fanout"
)
