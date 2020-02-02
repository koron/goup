package godlremote

// File represents a file on the golang.org downloads page.
// It should be kept in sync with the upload code in x/build/cmd/release.
//
// See https://github.com/golang/build/blob/5bb938ef020fb4b7f22d366b1e0dc8f9b425cc2f/cmd/release/upload.go#L46-L57
type File struct {
	Filename       string `json:"filename"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	Version        string `json:"version"`
	ChecksumSHA256 string `json:"sha256"`
	Size           int64  `json:"size"`
	Kind           string `json:"kind"` // "archive", "installer", "source"
}
