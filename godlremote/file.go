package godlremote

// File represents a file on the go.dev downloads page.
// It should be kept in sync with the upload code in x/build/cmd/release.
//
// See https://github.com/golang/build/blob/aa7fa4b2107cceba70eb6901c91e3be4b8bc419b/cmd/release/upload.go#L38-L47
type File struct {
	Filename       string `json:"filename"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	Version        string `json:"version"`
	ChecksumSHA256 string `json:"sha256"`
	Size           int64  `json:"size"`
	Kind           string `json:"kind"` // "archive", "installer", "source"
}
