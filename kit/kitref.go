package kit

type KitRefList struct {
	References []KitRef `yaml:"references"`
}

type KitRef struct {
	Alias  string `yaml:"alias"`
	Global bool   `yaml:"global"`
	Path   string `yaml:"path"`
	URL    string `yaml:"url"`
}

func newKitRefList() KitRefList {
	return KitRefList{References: nil}
}
