package lang

/*func compile(sourceDir string, config Config, version Version, stepCb func(step Step)) error {
	errOut, _ := os.Create(filepath.Join(config.Temp, version.Name()+"-build.error"))
	stdOut, _ := os.Create(filepath.Join(config.Temp, version.Name()+"-build.error"))

	defer errOut.Close()
	defer stdOut.Close()

	if stepCb == nil {
		stepCb = func(Step) {}
	}

	if runtime.GOOS == "windows" {

	} else {

		mkCmd := func(c string, args ...string) *exec.Cmd {
			cmd := exec.Command(c, args...)
			cmd.Dir = sourceDir
			cmd.Env = os.Environ()
			cmd.Stdout = stdOut
			cmd.Stderr = errOut
			return cmd
		}

		stepCb(Configure)
		target := filepath.Join(config.Source, version.Name())
		cmd := mkCmd("python", "./configure", "--prefix="+target)

		err := cmd.Run()
		if err != nil {
			return err
		}

		stepCb(Build)
		cmd = mkCmd("make")

		err = cmd.Run()

		if err != nil {
			return err
		}

		stepCb(Install)
		cmd = mkCmd("make", "install")

		err = cmd.Run()

		if err != nil {
			return err
		}

	}

	return nil

}*/
