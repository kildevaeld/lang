package lang

type Build struct {
	Interpreter string
	Exec        string
}

type Hash struct {
	Type  string
	Value string
}

type Source struct {
	Type   SourceType
	Link   string
	Target string
	Hash   Hash
}

type Version struct {
	Version string
	Arch    Arch
	Os      OS
	Binary  bool
	Latest  bool
	Source  Source
}

type Export struct {
	Binary  string
	Library string
}

type Definition struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stable      []Version `json:"stable"`
	Unstable    Version   `json:"unstable"`
	Export      Export
}
