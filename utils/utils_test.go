package utils

import (
	"fmt"
	"testing"
)

func TestValidateStrongPassword(t *testing.T) {
	b1 := ValidateStrongPassword("123")
	b2 := ValidateStrongPassword("123456789")
	b3 := ValidateStrongPassword("A1234567")
	b4 := ValidateStrongPassword("AB123456")
	b5 := ValidateStrongPassword("A@123")
	fmt.Println(b1, b2, b3, b4, b5)
}
