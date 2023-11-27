package fastly_test

import (
	"fmt"
	"log"
	"os"

	"github.com/fastly/go-fastly/v8/fastly"
)

func ExampleNewPaginator() {
	client, err := fastly.NewClient(os.Getenv("FASTLY_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// Supported type parameter values:
	// fastly.ACLEntry
	// fastly.DictionaryItem
	// fastly.Service
	p := fastly.NewPaginator[fastly.Service](
		// Fastly API client.
		client,
		// Fastly API inputs.
		&fastly.ListInput{
			Direction: "ascend",
			Sort:      "type",
			PerPage:   50,
		},
		// Fastly API path.
		//
		// Supported values:
		// fastly.ACLEntriesPath
		// fastly.DictionaryItemsPath
		// fastly.ServicePath
		fastly.ServicePath,
	)

	var results []*fastly.Service
	for p.HasNext() {
		data, err := p.GetNext()
		if err != nil {
			fmt.Printf("failed to get next page (remaining: %d): %s", p.Remaining(), err)
			return
		}
		results = append(results, data...)
	}

	fmt.Printf("%#v\n", len(results))
	for _, service := range results {
		fmt.Printf("%#v (%s)\n", service.Name, service.Type)
	}
}
