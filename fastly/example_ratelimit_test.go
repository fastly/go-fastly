package fastly_test

import (
	"context"
	"fmt"
	"log"

	"github.com/fastly/go-fastly/v13/fastly"
)

func ExampleClient_RateLimitRemaining() {
	token := "your_api_token"
	sid := "your_service_id"
	dictName := "your_dict_name"

	c, err := fastly.NewClient(token)
	if err != nil {
		log.Fatal(err)
	}

	v, err := c.LatestVersion(context.TODO(), &fastly.LatestVersionInput{ServiceID: sid})
	if err != nil {
		log.Fatal(err)
	}

	dict, err := c.GetDictionary(context.TODO(), &fastly.GetDictionaryInput{
		ServiceID:      sid,
		ServiceVersion: *v.Number,
		Name:           dictName,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.CreateDictionaryItem(context.TODO(), &fastly.CreateDictionaryItemInput{
		ServiceID:    sid,
		DictionaryID: *dict.DictionaryID,
		ItemKey:      fastly.ToPointer("test-dictionary-item"),
		ItemValue:    fastly.ToPointer("value"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Writes remaining before '429 Too Many Requests': %d\n", c.RateLimitRemaining())
	fmt.Printf("Next rate limit reset expected at %v\n", c.RateLimitReset())

	for i := 1; i < 10; i++ {
		_, err := c.UpdateDictionaryItem(context.TODO(), &fastly.UpdateDictionaryItemInput{
			ServiceID:    sid,
			DictionaryID: *dict.DictionaryID,
			ItemKey:      "test-dictionary-item",
			ItemValue:    fmt.Sprintf("value%d", i),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Writes remaining before '429 Too Many Requests': %d\n", c.RateLimitRemaining())
		fmt.Printf("Next rate limit reset expected at %v\n", c.RateLimitReset())
	}
}
