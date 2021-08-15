package version

// Info is information about the version of the application.
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"build_date"`
	Go      string `json:"go"`
}
