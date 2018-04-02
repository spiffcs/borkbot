package borkbot

type redditResponse struct {
	kind string
	data redditInfo
}

type redditInfo struct {
	children redditChild
}

type redditChild struct {
	data redditData
}

type redditData struct {
	url string
}
