package response

type Comment struct {
	ID         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_dat,"`
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}
