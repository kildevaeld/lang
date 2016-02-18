package lang

type Command struct {
	Name        string
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
	Build   []Command
}

type Export struct {
	Binary  string
	Library string
}

type Definition struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Stable      Versions `json:"stable"`
	Unstable    Versions `json:"unstable"`
	Export      Export
	Environment map[string]string
}

func (self *Version) Equal(v Version) bool {
	return v.Version == self.Version && self.Arch == v.Arch && self.Os == self.Os && self.Binary == self.Binary
}

type Versions []Version

func (self Versions) Contains(v Version) bool {
	for _, vv := range self {
		if vv.Equal(v) {
			return true
		}
	}
	return false
}
