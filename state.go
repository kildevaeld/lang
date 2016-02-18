package lang

type State struct {
	Installed map[string]Versions
}

func (self *State) Add(name string, v Version) {
	if _, ok := self.Installed[name]; !ok {
		self.Installed[name] = Versions{}
	}
	if self.List(name).Contains(v) {
		return
	}

	self.Installed[name] = append(self.Installed[name], v)
}

func (self *State) Remove(name string, v Version) {
	/*if _, ok := self.Installed; !ok {
			return
		}

	    for i, vv := range self.Installed {
	        if
	    }*/
}

func (self *State) List(name string) Versions {
	if v, ok := self.Installed[name]; ok {
		return v
	}
	return Versions{}
}
