//golangcitest:args -Egokeyword
//golangcitest:config_path testdata/configs/gokeyword.yml
package testdata

func GoKeyword() {
	go func() {}() // want `detected direct use of `go` keyword: derp`
}
