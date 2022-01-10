package fastly_test

import (
	"fmt"
	"log"

	"github.com/fastly/go-fastly/v5/fastly"
)

func ExampleClient_NewListServicesPaginator() {
	client, err := fastly.NewClient("your_api_token")
	if err != nil {
		log.Fatal(err)
	}

	paginator := client.NewListServicesPaginator(
		&fastly.ListServicesInput{
			PerPage: 50,
		},
	)

	var ss []*fastly.Service
	for paginator.HasNext() {
		data, err := paginator.GetNext()
		if err != nil {
			break
		}
		ss = append(ss, data...)
	}

	fmt.Printf("retrieved %d services\n", len(ss))
}
