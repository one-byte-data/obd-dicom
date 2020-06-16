package network

var (
	uniqid = 1
)

func Resetuniq() {
	uniqid = 1
}

func Uniq8() byte {
	uniqid++
	return byte(uniqid & 0xff)
}

func Uniq16() uint16 {
	uniqid++
	return uint16(uniqid & 0xffff)
}

func Uniq8odd() byte {

	if uniqid&0x01 == 1 {
		uniqid++
		return byte(uniqid & 0xff)
	}
	uniqid++
	uniqid++
	return byte(uniqid & 0xff)
}

func Uniq16odd() uint16 {

	if uniqid&0x01 == 1 {
		uniqid++
		return uint16(uniqid & 0xffff)
	}
	uniqid++
	uniqid++
	return uint16(uniqid & 0xffff)
}
