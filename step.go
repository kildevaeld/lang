//go:generate stringer -type=Step,Arch,OS,SourceType
package lang

import (
	"fmt"
	"strings"
)

type Step int
type Arch int
type OS int
type SourceType int

const (
	Unpack Step = iota
	Download
	Compile
	Configure
	//Build
	Install
)

const (
	Linux OS = iota + 100
	Darwin
	Windows
	Android
)

const (
	X86 Arch = iota + 200
	X64
	Arm
	Armbe
	Arm64
)

const (
	Git SourceType = iota + 300
	URL
)

func (self SourceType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self)), nil
}

func (self *SourceType) UnmarshalJSON(b []byte) error {

	switch strings.Replace(strings.ToLower(string(b)), "\"", "", 2) {
	case "url":
		*self = URL
		break
	case "git":
		*self = Git
		break

	}

	return nil
}

func (self Arch) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self)), nil
}

func (self *Arch) UnmarshalJSON(b []byte) error {
	switch strings.Replace(strings.ToLower(string(b)), "\"", "", 2) {
	case "x64", "amd64":
		*self = X64
		break
	case "386", "x86":
		*self = X86
		break
	case "arm":
		*self = Arm
		break
	}

	return nil
}

func (self Step) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self)), nil
}

func (self *Step) UnmarshalJSON(b []byte) error {
	switch strings.Replace(strings.ToLower(string(b)), "\"", "", 2) {
	case "Unpack", "unpack":
		*self = Unpack
	}

	return nil
}

func (self OS) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", self)), nil
}

func (self *OS) UnmarshalJSON(b []byte) error {
	switch strings.Replace(strings.ToLower(string(b)), "\"", "", 2) {
	case "linux":
		*self = Linux
		break
	case "windows", "win":
		*self = Windows
		break
	case "darwin", "ios", "osx":
		*self = Darwin
		break
	}

	return nil
}
