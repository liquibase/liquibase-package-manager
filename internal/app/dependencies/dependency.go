package dependencies

type Dependency map[string]string

func (d Dependency) GetName() string {
	var r string
	for k := range d {
		r = k
	}
	return r
}

func (d Dependency) GetVersion() string {
	var r string
	for _, v := range d {
		r = v
	}
	return r
}