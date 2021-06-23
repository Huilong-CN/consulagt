package consulagt

import "math/rand"

//Random 随机一个可用服务器
func Random(svrs []string) (string, bool) {
	count := len(svrs)
	if count <= 0 {
		return "", false
	} else if count == 1 {
		return svrs[0], true
	}
	return svrs[rand.Intn(len(svrs))], true
}
