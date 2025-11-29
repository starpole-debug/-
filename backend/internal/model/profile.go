package model

// UserHome aggregates data needed for the personal homepage.
type UserHome struct {
	User        *User             `json:"user"`
	Assets      *UserAsset        `json:"assets"`
	Stats       UserHomeStats     `json:"stats"`
	Favorites   []CommunityPost   `json:"favorites"`
	RecentViews []CommunityPost   `json:"recent_views"`
	MyPosts     []CommunityPost   `json:"my_posts"`
	CreatorHint CreatorHomePrompt `json:"creator"`
}

// UserHomeStats exposes key metrics for the dashboard.
type UserHomeStats struct {
	SessionTotal    int `json:"session_total"`
	FavoriteTotal   int `json:"favorite_total"`
	RecentViewTotal int `json:"recent_view_total"`
	PostTotal       int `json:"post_total"`
}

// CreatorHomePrompt hints whether the user has creator presence.
type CreatorHomePrompt struct {
	RolesTotal       int   `json:"roles_total"`
	PublishedRoles   int   `json:"published_roles"`
	WalletBalance    int64 `json:"wallet_balance"`
	HasCreatorAccess bool  `json:"has_creator_access"`
}
