package productcore

import "github.com/fastly/go-fastly/v9/fastly"

func makeURL(productID, serviceID string, subComponents []string) string {
	path := []string{"enabled-products", "v1", productID}

	if serviceID != "" {
		path = append(path, "services", serviceID)
	}

	if len(subComponents) > 0 {
		path = append(path, subComponents...)
	}

	return fastly.ToSafeURL(path...)
}
