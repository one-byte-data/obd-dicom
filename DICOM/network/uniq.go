package network

var (
	uniqid = 1
)

func Resetuniq() {
	uniqid = 1
}

func uniq8() byte {
	uniqid++
	return byte(uniqid & 0xff);
}

func uniq8odd() byte {

	if uniqid&0x01 == 1 {
		uniqid++
		return byte(uniqid & 0xff)
	}
	uniqid++
	uniqid++
	return byte(uniqid & 0xff)
}
