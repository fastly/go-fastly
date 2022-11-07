package fastly_test

import (
	"fmt"
	"log"

	"github.com/fastly/go-fastly/v7/fastly"
)

func ExampleClient_NewListACLEntriesPaginator() {
	client, err := fastly.NewClient("your_api_token")
	if err != nil {
		log.Fatal(err)
	}

	paginator := client.NewListACLEntriesPaginator(
		&fastly.ListACLEntriesInput{
			ServiceID: "your_service_id",
			ACLID:     "your_acl_id",
			PerPage:   50,
		},
	)

	var es []*fastly.ACLEntry
	for paginator.HasNext() {
		data, err := paginator.GetNext()
		if err != nil {
			break
		}
		es = append(es, data...)
	}

	fmt.Printf("retrieved %d ACL entries\n", len(es))
}
