package operations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v13/fastly"
)

func TestClient_Tags(t *testing.T) {
	ctx := context.TODO()

	serviceID := fastly.TestDeliveryServiceID

	var err error

	// Create tag.
	const tagName = "test-tag-tags"
	var tag *OperationTag
	fastly.Record(t, "create_tag", func(c *fastly.Client) {
		tag, err = CreateTag(ctx, c, &CreateTagInput{
			ServiceID:   fastly.ToPointer(serviceID),
			Name:        fastly.ToPointer(tagName),
			Description: fastly.ToPointer("go-fastly test tag"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, tag)
	require.NotEmpty(t, tag.ID)
	require.Equal(t, tagName, tag.Name)
	require.Equal(t, "go-fastly test tag", tag.Description)

	defer func() {
		fastly.Record(t, "delete_tag", func(c *fastly.Client) {
			_ = DeleteTag(ctx, c, &DeleteTagInput{
				ServiceID: fastly.ToPointer(serviceID),
				TagID:     fastly.ToPointer(tag.ID),
			})
		})
	}()

	// Describe tag.
	var described *OperationTag
	fastly.Record(t, "describe_tag", func(c *fastly.Client) {
		described, err = DescribeTag(ctx, c, &DescribeTagInput{
			ServiceID: fastly.ToPointer(serviceID),
			TagID:     fastly.ToPointer(tag.ID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, described)
	require.Equal(t, tag.ID, described.ID)
	require.Equal(t, tagName, described.Name)

	// Update tag.
	// NOTE: The API currently requires the tag "name" to be provided on PATCH,
	// even when only updating other fields (e.g. description).
	var updated *OperationTag
	fastly.Record(t, "update_tag", func(c *fastly.Client) {
		updated, err = UpdateTag(ctx, c, &UpdateTagInput{
			ServiceID:   fastly.ToPointer(serviceID),
			TagID:       fastly.ToPointer(tag.ID),
			Name:        fastly.ToPointer(tagName),
			Description: fastly.ToPointer("updated"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, tag.ID, updated.ID)
	require.Equal(t, tagName, updated.Name)
	require.Equal(t, "updated", updated.Description)

	// List tags.
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
		if item.ID == tag.ID {
			found = true
			break
		}
	}
	require.True(t, found, "expected created tag to appear in list")
}
