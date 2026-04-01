package operations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v14/fastly"
)

func TestClient_Tags(t *testing.T) {
	ctx := context.TODO()

	serviceID := fastly.TestDeliveryServiceID

	var err error

	// Create tag #1.
	const tagName1 = "test-tag-tags"
	var tag1 *OperationTag
	fastly.Record(t, "create_tag", func(c *fastly.Client) {
		tag1, err = CreateTag(ctx, c, &CreateTagInput{
			ServiceID:   fastly.ToPointer(serviceID),
			Name:        fastly.ToPointer(tagName1),
			Description: fastly.ToPointer("go-fastly test tag"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, tag1)
	require.NotEmpty(t, tag1.ID)
	require.Equal(t, tagName1, tag1.Name)
	require.Equal(t, "go-fastly test tag", tag1.Description)

	defer func() {
		fastly.Record(t, "delete_tag", func(c *fastly.Client) {
			_ = DeleteTag(ctx, c, &DeleteTagInput{
				ServiceID: fastly.ToPointer(serviceID),
				TagID:     fastly.ToPointer(tag1.ID),
			})
		})
	}()

	// Describe tag.
	var described *OperationTag
	fastly.Record(t, "describe_tag", func(c *fastly.Client) {
		described, err = DescribeTag(ctx, c, &DescribeTagInput{
			ServiceID: fastly.ToPointer(serviceID),
			TagID:     fastly.ToPointer(tag1.ID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, described)
	require.Equal(t, tag1.ID, described.ID)
	require.Equal(t, tagName1, described.Name)

	// Update tag.
	// NOTE: The API currently requires the tag "name" to be provided on PATCH,
	// even when only updating other fields (e.g. description).
	var updated *OperationTag
	fastly.Record(t, "update_tag", func(c *fastly.Client) {
		updated, err = UpdateTag(ctx, c, &UpdateTagInput{
			ServiceID:   fastly.ToPointer(serviceID),
			TagID:       fastly.ToPointer(tag1.ID),
			Name:        fastly.ToPointer(tagName1),
			Description: fastly.ToPointer("updated"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, tag1.ID, updated.ID)
	require.Equal(t, tagName1, updated.Name)
	require.Equal(t, "updated", updated.Description)

	// Create tag #2 (for exercising tag create/delete recordings and pagination fixture stability).
	const tagName2 = "test-tag-tags-pagination"
	var tag2 *OperationTag
	fastly.Record(t, "create_tag_2", func(c *fastly.Client) {
		tag2, err = CreateTag(ctx, c, &CreateTagInput{
			ServiceID:   fastly.ToPointer(serviceID),
			Name:        fastly.ToPointer(tagName2),
			Description: fastly.ToPointer("go-fastly test tag 2"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, tag2)
	require.NotEmpty(t, tag2.ID)
	require.Equal(t, tagName2, tag2.Name)

	defer func() {
		fastly.Record(t, "delete_tag_2", func(c *fastly.Client) {
			_ = DeleteTag(ctx, c, &DeleteTagInput{
				ServiceID: fastly.ToPointer(serviceID),
				TagID:     fastly.ToPointer(tag2.ID),
			})
		})
	}()

	// List tags (existing behavior).
	var tags *OperationTags
	fastly.Record(t, "list_tags", func(c *fastly.Client) {
		tags, err = ListTags(ctx, c, &ListTagsInput{
			ServiceID: fastly.ToPointer(serviceID),
			Limit:     fastly.ToPointer(100),
			Page:      fastly.ToPointer(0),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, tags)

	// Avoid brittle assertions: just ensure our created tag exists in the returned data.
	found := false
	for _, item := range tags.Data {
		if item.ID == tag1.ID {
			found = true
			break
		}
	}
	require.True(t, found, "expected created tag to appear in list")

	// ---- Pagination test for tags ----
	//
	// IMPORTANT: Tag listing order is not guaranteed, so don't assert that our newly created tags
	// must appear in the first two pages. Instead, exercise the paginator logic (page increments,
	// limit respected) and assert that page 0 and page 1 return different items.
	limit := 1
	p := NewTagPaginator(ctx, fastly.TestClient, &ListTagsInput{
		ServiceID: fastly.ToPointer(serviceID),
		Limit:     &limit,
		Page:      fastly.ToPointer(0),
	})

	var (
		page0 []OperationTag
		page1 []OperationTag
	)

	fastly.Record(t, "list_tags_page_0", func(c *fastly.Client) {
		p.SetClient(c)
		page0, err = p.GetNext()
	})
	require.NoError(t, err)
	require.Len(t, page0, 1)

	fastly.Record(t, "list_tags_page_1", func(c *fastly.Client) {
		p.SetClient(c)
		page1, err = p.GetNext()
	})
	require.NoError(t, err)
	require.Len(t, page1, 1)

	// Pagination sanity: different pages should not return the same item.
	require.NotEqual(t, page0[0].ID, page1[0].ID)

	// Optional sanity: if the API reports total > fetched, paginator should still indicate more pages.
	require.True(t, p.HasNext())
}
