package lang

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type paths struct {
	root string
}

func (self paths) Root(args ...string) string {
	list := []string{self.root}
	return filepath.Join(append(list, args...)...)
}

func (self paths) Cache(args ...string) string {
	list := []string{self.root, "cache"}
	return filepath.Join(append(list, args...)...)
}

func (self paths) Source(args ...string) string {
	list := []string{self.root, "sources"}
	return filepath.Join(append(list, args...)...)
}

func (self paths) Current(args ...string) string {
	list := []string{self.root, "current"}
	return filepath.Join(append(list, args...)...)
}

func (self paths) Temp(args ...string) string {
	list := []string{self.root, "temp"}
	return filepath.Join(append(list, args...)...)
}

func (self paths) Ensure() error {

	if err := ensureDir(self.Temp()); err != nil {
		return err
	}

	if err := ensureDir(self.Cache()); err != nil {
		return err
	}

	if err := ensureDir(self.Source()); err != nil {
		return err
	}

	if err := ensureDir(self.Current()); err != nil {
		return err
	}

	if err := ensureDir(self.Root()); err != nil {
		return err
	}

	return nil
}

type Language struct {
	paths      paths
	definition Definition
}

func (self *Language) List() []Version {
	return self.definition.Stable
}

func (self *Language) Definition() Definition {
	return self.definition
}

func (self *Language) GetName() string {
	return self.definition.Name
}

func (self *Language) GetVersion(version string, oss OS, arch Arch, binary bool) *Version {

	if "unstable" == strings.ToLower(version) {
		return &self.definition.Unstable[0]
	}
	ver := version
	for _, v := range self.definition.Stable {
		if version == "latest" && v.Os == oss && v.Arch == arch && v.Binary == binary {
			if v.Latest {
				return &v
			}
			continue
		}

		if v.Version[0] == 'v' && version[0] != 'v' {
			ver = "v" + version
		} else if v.Version[0] != 'v' && version[0] == 'v' {
			ver = version[1:]
		}

		if v.Version == ver && v.Os == oss && v.Arch == arch && v.Binary == binary {
			return &v
		}
	}

	return nil

}

func (self *Language) Use(version Version) error {
	sourceDir := self.paths.Source(fmt.Sprintf("%s-%s-%s", version.Version, version.Os, version.Arch))
	fmt.Printf("Sourcedir %s", sourceDir)
	if !dirExists(sourceDir) {
		return fmt.Errorf("Not installed")
	}

	target := self.paths.Current()

	if dirExists(target) {
		os.RemoveAll(target)
	}

	err := os.Symlink(sourceDir, target)

	if err != nil {
		return err
	}

	versionFile := self.paths.Root("version")
	ioutil.WriteFile(versionFile, []byte(version.Version), 0755)

	return nil

}

func (self *Language) Remove(version Version) error {
	//sourceDir := self.paths.Source(fmt.Sprintf("%s-%s-%s", version.Version, version.Os, version.Arch))
	return nil
}

func (self *Language) Install(version Version, progressCB func(step Step, progress, total int64)) error {

	if progressCB == nil {
		progressCB = func(step Step, progress, total int64) {}
	}

	target, err := self.download(version, progressCB)

	if err != nil {
		return err
	}

	defer os.RemoveAll(target)

	sourceDir := self.paths.Source(fmt.Sprintf("%s-%s-%s", version.Version, version.Os, version.Arch))

	shouldCopy := false
	if len(version.Build) > 0 {
		//progressCB(Compile, 0, 0)
		if err := self.compile(target, sourceDir, version); err != nil {
			return err
		}
		if !dirExists(sourceDir) {
			shouldCopy = true
		}

	} else {
		shouldCopy = true
	}

	if shouldCopy {
		total := analyzeDir(target)
		var index int64 = 0
		if err := copyDir(target, sourceDir, func(s, t string) {
			index++
			progressCB(Install, index, total)
		}); err != nil {
			return err
		}
	}

	return err
}

func (self *Language) download(version Version, progressCB func(step Step, progress, total int64)) (string, error) {

	if version.Source.Type == URL {
		fileName := filepath.Base(version.Source.Link)
		localFile := self.paths.Cache(fileName)

		if !fileExists(localFile) {
			err := download(version.Source.Link, localFile, progressCB)

			if err != nil {
				return "", err
			}

		}

		tmpDir := self.paths.Temp(fmt.Sprintf("%s-%s-%s", version.Version, version.Os, version.Arch))
		if version.Source.Target != "" {
			tmpDir = self.paths.Temp(version.Source.Target)
		}

		if !dirExists(tmpDir) {

			if version.Source.Hash.Type != "" {
				if err := ValidateFile(version.Source.Hash.Type, version.Source.Hash.Value, localFile); err != nil {
					os.Remove(localFile)
					return "", err
				}
			}

			err := unpack(localFile, self.paths.Temp(), func(progress, total int64) {
				progressCB(Unpack, progress, total)
			})

			if err != nil {
				os.Remove(localFile)
				return "", err
			}
		}

		return tmpDir, nil

	} else {
		return "", errors.New("not implemented yet")
	}

}

func (self *Language) compile(working, prefix string, version Version) error {
	return compile(working, prefix, self, version)
}

func NewLanguage(path string, def Definition) (*Language, error) {

	lang := &Language{paths{path}, def}

	err := lang.paths.Ensure()

	return lang, err
}
