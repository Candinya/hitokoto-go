package types

type Version struct {
	ProtocolVersion string `json:"protocol_version"` // only supports 1.0.0 now
	BundleVersion   string `json:"bundle_version"`
	UpdatedAt       uint64 `json:"updated_at"`
	Categories      struct {
		Path      string `json:"path"`
		TimeStamp uint64 `json:"timestamp"`
	} `json:"categories"`
	Sentences []struct {
		Name      string `json:"name"`
		Key       string `json:"key"`
		Path      string `json:"path"`
		TimeStamp uint64 `json:"timestamp"`
	} `json:"sentences"`
}
