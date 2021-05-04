package util

type StringKeyMap map[string]interface{}

func (r *StringKeyMap) Keys() []string {
	var ks []string
	for k := range *r {
		ks = append(ks, k)
	}

	return ks
}
