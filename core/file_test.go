package core

import (
	"fmt"
	"testing"
)

func TestFile(t *testing.T) {
	ok := DoesFileExist("../core/core.go")
	fmt.Println(ok)
}
