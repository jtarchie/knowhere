package runtime

import "github.com/jtarchie/knowhere/address"

type Address struct{}

func (a *Address) Parse(fullAddress string) (map[string]string, bool) {
	return address.Parse(fullAddress, true)
}
