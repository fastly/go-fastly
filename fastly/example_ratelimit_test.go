package fastly_test

import (
	"fmt"
	"log"

	"github.com/fastly/go-fastly/v10/fastly"
)

func ExampleClient_RateLimitRemaining() {
	token := "your_api_token"
	sid := "your_service_id"
	dictName := "your_dict_name"

	c, err := fastly.NewClient(token)
	if err != nil {
		log.Fatal(err)
	}

	v, err := c.LatestVersion(&fastly.LatestVersionInput{ServiceID: sid})
	if err != nil {
		log.Fatal(err)
	}

	dict, err := c.GetDictionary(&fastly.GetDictionaryInput{
		ServiceID:      sid,
		ServiceVersion: *v.Number,
		Name:           dictName,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.CreateDictionaryItem(&fastly.CreateDictionaryItemInput{
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
		_, err := c.UpdateDictionaryItem(&fastly.UpdateDictionaryItemInput{
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
