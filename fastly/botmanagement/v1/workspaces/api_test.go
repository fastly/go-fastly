package workspaces

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_Workspaces(t *testing.T) {
	ctx := context.TODO()

	var err error

	// List workspaces. Workspaces cannot be created through the API, so the
	// test account must already have at least one.
	var workspaces []Workspace
	fastly.Record(t, "list_workspaces", func(c *fastly.Client) {
		workspaces, err = List(ctx, c, &ListInput{})
	})
	require.NoError(t, err)
	require.NotEmpty(t, workspaces, "expected at least one workspace")

	// List workspaces 10 per page to exercise auto-pagination. The result
	// should match the single-page listing.
	var paginated []Workspace
	fastly.Record(t, "list_workspaces_paginated", func(c *fastly.Client) {
		paginated, err = List(ctx, c, &ListInput{
			Limit: new(10),
		})
	})
	require.NoError(t, err)
	require.Len(t, paginated, len(workspaces))

	workspace := workspaces[0]

	// Get the workspace.
	var got *Workspace
	fastly.Record(t, "get_workspace", func(c *fastly.Client) {
		got, err = Get(ctx, c, &GetInput{
			WorkspaceID: new(workspace.ID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, workspace.ID, got.ID)
	require.Equal(t, workspace.Name, got.Name)
	require.Equal(t, workspace.ProtectionMode, got.ProtectionMode)

	// Restore the original description at the end.
	defer func() {
		fastly.Record(t, "restore_workspace", func(c *fastly.Client) {
			_, err = Update(ctx, c, &UpdateInput{
				WorkspaceID: new(workspace.ID),
				Description: new(workspace.Description),
			})
		})
		if err != nil {
			t.Errorf("error restoring workspace: %v", err)
		}
	}()

	// Update the workspace description.
	const updatedDescription = "go-fastly test description"

	var updated *Workspace
	fastly.Record(t, "update_workspace", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			WorkspaceID: new(workspace.ID),
			Description: new(updatedDescription),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, workspace.ID, updated.ID)
	require.Equal(t, updatedDescription, updated.Description)
	require.Equal(t, workspace.ProtectionMode, updated.ProtectionMode)
}

func TestClient_Workspaces_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Get(ctx, fastly.TestClient, &GetInput{WorkspaceID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingWorkspaceID)

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{WorkspaceID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingWorkspaceID)
}
