package common

var Foo = 5

func Must(e error) {
	if e != nil {
		panic(e)
	}
}
