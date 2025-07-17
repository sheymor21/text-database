package tdb

import "log"

type generic interface {
	[]byte | string | interface{}
}

func must[T generic](data T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return data
}
func errorHandler(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
