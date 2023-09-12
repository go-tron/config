package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	fmt.Println(c)
}
