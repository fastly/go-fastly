package productcore

import "github.com/fastly/go-fastly/v19/fastly"

func makeURL(productID string, subComponents []string) string {
	path := []string{"enabled-products", "v1", productID}

	if len(subComponents) > 0 {
		path = append(path, subComponents...)
	}

	return fastly.ToSafeURL(path...)
}
