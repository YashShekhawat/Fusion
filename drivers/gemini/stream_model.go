package gemini

type streamResponse struct {
	Candidates []streamCandidate `json:"candidates"`
}

type streamCandidate struct {
	Content streamContent `json:"content"`
}

type streamContent struct {
	Parts []streamPart `json:"parts"`
}

type streamPart struct {
	Text string `json:"text"`
}
