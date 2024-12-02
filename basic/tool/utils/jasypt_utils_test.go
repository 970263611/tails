package utils

import (
	"fmt"
	"testing"
)

func TestJasyptDec(t *testing.T) {

	dMsg, err := Decrypt("LD+PW0BgUXZ1RhWFVgDioQ==", "zhh.")

	fmt.Printf("%v\n%v\n ", dMsg, err)

	dec, err := JasyptDec("LD+PW0BgUXZ1RhWFVgDioQ==", "zhh.")

	fmt.Printf("%v\n%v\n ", dec, err)

}
