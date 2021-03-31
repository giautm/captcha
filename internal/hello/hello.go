package hello

import (
	"fmt"
)

// Hello return hello message
func Hello(name string) string {
	return fmt.Sprintf("hello %s", name)
}
