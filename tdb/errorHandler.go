package tdb

import "log"

type generic interface {
	[]byte | string | interface{}
}

// must is a generic helper function that handles errors and returns the data if no error occurred.
// It takes a generic type T that can be []byte, string, or interface{}, along with an error.
// If an error is present, it logs the error and terminates the program.
// Returns the original data if no error occurred.
func must[T generic](data T, err error) T {
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

// errorHandler checks if an error is present and logs it fatally if it exists.
// It takes an error parameter and terminates the program if the error is not nil.
func errorHandler(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
