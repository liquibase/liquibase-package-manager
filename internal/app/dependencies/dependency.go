package dependencies

//Dependency main "package":"tag" object
type Dependency map[string]string

//GetName get key from Dependency map
func (d Dependency) GetName() string {
	var r string
	for k := range d {
		r = k
	}
	return r
}

//GetVersion get value from Dependency map
func (d Dependency) GetVersion() string {
	var r string
	for _, v := range d {
		r = v
	}
	return r
}