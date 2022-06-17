package response

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

type Video struct {
	ID            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,"`
	CommentCount  int64  `json:"comment_count,"`
	IsFavorite    bool   `json:"is_favorite"`
}

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
