//golangcitest:args -Egokeyword
package testdata

func GoKeyword() {
	go func() {}() // want `detected direct use of `go` keyword`
}
