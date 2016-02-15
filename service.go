package lang

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

var HostArch Arch
var HostOs OS

func init() {
	switch runtime.GOOS {
	case "darwin":
		HostOs = Darwin
		break
	case "windows":
		HostOs = Windows
		break
	case "linux":
		HostOs = Linux
		break
	case "android":
		HostOs = Android
	}

	switch runtime.GOARCH {
	case "amd64":
		HostArch = X64
	case "386":
		HostArch = X86
	}
}

type Config struct {
	Root string
}

type Service struct {
	config Config
	langs  map[string]*Language
}

func (self *Service) Languages() []string {
	var out []string
	for k, _ := range self.langs {
		out = append(out, k)
	}
	return out
}

func (self *Service) HasLanguage(name string) bool {
	for ln, _ := range self.langs {
		if ln == strings.ToLower(name) {
			return true
		}
	}
	return false
}

func (self *Service) GetLanguage(name string) *Language {
	for ln, ll := range self.langs {
		if ln == strings.ToLower(name) {
			return ll
		}
	}
	return nil
}

func (self *Service) AddDefinition(d Definition) (*Language, error) {
	lcName := strings.ToLower(d.Name)
	path := filepath.Join(self.config.Root, lcName)
	lang, err := NewLanguage(path, d)
	if err != nil {
		return nil, err
	}
	self.langs[lcName] = lang
	return lang, nil
}

func (self *Service) Use(lang, version string, binary bool) error {
	var language *Language
	var ok bool
	if language, ok = self.langs[strings.ToLower(lang)]; !ok {
		return fmt.Errorf("No language: %s", lang)
	}

	v := language.GetVersion(version, HostOs, HostArch, binary)

	if v == nil {
		v := language.GetVersion(version, HostOs, HostArch, !binary)
		if v == nil {
			return fmt.Errorf("Version %s not found for language %s", v, lang)
		}
	}

	return language.Use(*v)
}

func (self *Service) Environ() []string {
	var path, library []string
	for _, ln := range self.langs {
		export := ln.Definition().Export
		if export.Binary != "" {
			path = append(path, ln.paths.Current(export.Binary))
		}
		if export.Library != "" {
			library = append(library, ln.paths.Current(export.Library))
		}
	}
	path = append(path, "$PATH")
	library = append(library, "$LD_LIBRARY_PATH")

	return []string{"PATH=" + strings.Join(path, ":"), "LD_LIBRARY_PATH=" + strings.Join(library, ":")}

}

func (self *Service) Install(lang, v string, binary bool, progressCB func(step Step, progres, total int64)) error {

	var language *Language
	var ok bool
	if language, ok = self.langs[strings.ToLower(lang)]; !ok {
		return fmt.Errorf("No language: %s", lang)
	}

	version := language.GetVersion(v, HostOs, HostArch, binary)

	if version == nil {
		version := language.GetVersion(v, HostOs, HostArch, !binary)
		if version == nil {
			return fmt.Errorf("Version %s not found for language %s", v, lang)
		}
	}

	return language.Install(*version, progressCB)

}

func New(config Config) *Service {
	return &Service{
		config: config,
		langs:  make(map[string]*Language),
	}
}
