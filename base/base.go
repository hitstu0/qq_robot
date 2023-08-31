package base

var (
	sandboxUrl  string = "https://sandbox.api.sgroup.qq.com/"
	formalUrl   string = "https://api.sgroup.qq.com/"

	test        bool   = true
)

func GetUrl() string {
	if test {
		return sandboxUrl
	} else {
		return formalUrl
	}
}