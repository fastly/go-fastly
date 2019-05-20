package fastly

import (
	"testing"
	"time"
)

func TestClient_GetDictionaryInfo(t *testing.T) {
	var (
		nd  *DictionaryInfo
		err error
	)
	record(t, "dictionary_info/get", func(c *Client) {
		nd, err = c.GetDictionaryInfo(&GetDictionaryInfoInput{
			Service: testServiceID,
			Version: 682,
			ID:      "6mJXuCi2Mf19uHRefZJZBg",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if !nd.LastUpdated.Equal(time.Date(2016, time.May, 3, 16, 11, 41, 0, time.UTC)) {
		t.Errorf("bad last_updated: %v", nd.LastUpdated)
	}
	if nd.ItemCount != 4 {
		t.Errorf("bad item_count: %d", nd.ItemCount)
	}
	if nd.Digest != "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a" {
		t.Errorf("bad digest: %q", nd.Digest)
	}
}

func TestClient_GetDictionaryInfo_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDictionaryInfo(&GetDictionaryInfoInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryInfo(&GetDictionaryInfoInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDictionaryInfo(&GetDictionaryInfoInput{
		Service: "foo",
		Version: 1,
		ID:      "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
