package utilities

import "log"

type generic interface {
	[]byte | string | interface{}
}

func Must[T generic](data T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return data
}
func ErrorHandler(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
