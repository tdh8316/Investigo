package scan

type Result struct {
	Username string
	Site     string

	URLTemplate   string
	ProbeTemplate string
	Link          string

	Exists  bool
	Proxied bool
	Err     error
}

type Config struct {
	UserAgent    string
	WithTor      bool
	Download     bool
	Concurrency  int
	MaxBodyBytes int64
}

type ValidationFailure struct {
	Site          string
	UsedUsername  string
	UnusedUsername string

	Used   Result
	Unused Result
}
