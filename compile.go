package lang

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Template struct {
	Source string
}

func Interpolate(name, str string, ctx Template) (string, error) {
	t, err := template.New(name).Parse(str)

	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, ctx)
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

func run_command(build Command, workingDir string, template Template, stdout io.Writer, stderr io.Writer) (err error) {

	if HostOs == Windows {
		return errors.New("windows not supported yet")
	} else {
		cmdStr := build.Exec

		if cmdStr, err = Interpolate(build.Exec, cmdStr, template); err != nil {
			return err
		}

		args := strings.Split(cmdStr, " ")
		if build.Interpreter != "" {
			cmdStr = build.Interpreter
		} else {
			cmdStr = args[0]
			if len(args) == 1 {
				args = []string{}
			} else {
				args = args[1:]
			}
		}

		cmd := exec.Command(cmdStr, args...)

		cmd.Dir = workingDir
		cmd.Env = os.Environ()
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		return cmd.Run()
	}

}

func compile(workingDir, prefix string, lang *Language, version Version) (err error) {

	errOut, _ := os.Create(lang.paths.Temp(fmt.Sprintf("%s-build.error", version.Version)))
	stdOut, _ := os.Create(lang.paths.Temp(fmt.Sprintf("%s-build.log", version.Version)))

	defer errOut.Close()
	defer stdOut.Close()
	template := Template{
		Source: prefix,
	}
	for _, cmd := range version.Build {
		if err = run_command(cmd, workingDir, template, stdOut, errOut); err != nil {
			return err
		}
	}

	return nil
}
