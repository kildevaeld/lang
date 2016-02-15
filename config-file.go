package lang

import "os"

const LANG_ROOT_ENV = "LANG_ROOT"

// ConfigDir returns the configuration directory for Packer.
func ConfigDir() (string, error) {

	if langRoot := os.Getenv(LANG_ROOT_ENV); langRoot != "" {
		return langRoot, nil
	}
	return configDir()

}
