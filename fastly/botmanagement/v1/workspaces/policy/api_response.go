package policy

// Action is the action applied to matching requests.
type Action struct {
	// Type is the type of action to apply.
	Type string `json:"type"`
}

// Category is a group of bots with similar characteristics or behaviors.
// Categories are defined by Fastly and cannot be created or deleted.
type Category struct {
	// CategoryID is the unique identifier of the category.
	CategoryID string `json:"category_id"`
	// Name is the display name of the category.
	Name string `json:"name"`
	// Description is a description of the bots the category contains.
	Description string `json:"description,omitempty"`
	// Action is the action applied to every bot in the category that is
	// configured with the `inherit` action. Must be one of `allow`,
	// `block`, or `challenge`.
	Action Action `json:"action"`
}

// Categories is the API response structure for the list categories
// operation.
type Categories struct {
	// Data is the list of returned categories.
	Data []Category `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// Bot is a single bot within a category. Bots are defined by Fastly and
// cannot be created or deleted.
type Bot struct {
	// BotID is the unique identifier of the bot.
	BotID string `json:"bot_id"`
	// Name is the display name of the bot.
	Name string `json:"name"`
	// CategoryID is the ID of the category the bot belongs to. It is only
	// populated when listing all bots in a workspace.
	CategoryID string `json:"category_id,omitempty"`
	// Action is the action applied to requests from the bot. Must be one of
	// `allow`, `block`, `challenge`, or `inherit`, which applies the action of
	// the category the bot belongs to.
	Action Action `json:"action"`
}

// Bots is the API response structure for the list bots operations.
type Bots struct {
	// Data is the list of returned bots.
	Data []Bot `json:"data"`
	// Meta contains pagination metadata.
	Meta Meta `json:"meta"`
}

// Meta is the pagination metadata returned by the list operations.
type Meta struct {
	// Limit is the `limit` value used when making the request.
	Limit *int `json:"limit"`
	// Total is the count of results matching the request.
	Total *int `json:"total"`
	// NextCursor is the cursor used to retrieve the next page of results. It
	// is nil when there are no further pages.
	NextCursor *string `json:"next_cursor"`
}
