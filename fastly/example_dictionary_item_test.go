package fastly_test

import (
	"fmt"
	"log"

	"github.com/fastly/go-fastly/v5/fastly"
)

func ExampleClient_NewListDictionaryItemsPaginator() {
	client, err := fastly.NewClient("your_api_token")
	if err != nil {
		log.Fatal(err)
	}

	paginator := client.NewListDictionaryItemsPaginator(
		&fastly.ListDictionaryItemsInput{
			ServiceID:    "your_service_id",
			DictionaryID: "your_dictionary_id",
			PerPage:      50,
		},
	)

	var es []*fastly.DictionaryItem
	for paginator.HasNext() {
		data, err := paginator.GetNext()
		if err != nil {
			break
		}
		es = append(es, data...)
	}

	fmt.Printf("retrieved %d Dictionary items\n", len(es))
}
