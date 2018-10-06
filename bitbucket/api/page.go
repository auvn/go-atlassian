package api

type Page struct {
	Size          int64 `json:"size"`
	IsLastPage    bool  `json:"isLastPage"`
	Limit         int64 `json:"limit"`
	NextPageStart int64 `json:"nextPageStart"`
}

func (p Page) Next(path Path) Path {
	return Path{
		Path: path.Path,
		Params: path.Params.
			WithInt64("limit", p.Limit).
			WithInt64("start", p.NextPageStart),
	}
}
