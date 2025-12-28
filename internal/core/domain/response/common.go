package modelResponse

type Common struct {
	HttpCode int `json:"-"`

	Data     any             `json:"data"`
	Metadata *CommonMetadata `json:"metadata,omitempty"`
}

type CommonMetadata struct {
	TotalResult int  `json:"total_result"`
	Sort        Sort `json:"sort"`
}

type Sort struct {
	Key   string `json:"key"`
	Order string `json:"order"`
}
