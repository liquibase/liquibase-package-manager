package lpm

//Dependency main "package":"tag" object
type Dependency map[string]string

//NewDependency instantiates a Dependency from a pkgname and a tag
func NewDependency(pkgname, tag string) Dependency {
	return Dependency{pkgname: tag}
}

//GetName get key from Dependency map
func (d Dependency) GetName() string {
	var r string
	for k := range d {
		r = k
		break
	}
	return r
}

//GetVersion get value from Dependency map
func (d Dependency) GetVersion() string {
	var r string
	for _, v := range d {
		r = v
		break
	}
	return r
}
