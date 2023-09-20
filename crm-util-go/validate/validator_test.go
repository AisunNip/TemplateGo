package validate_test

import (
	"crm-util-go/validate"
	"fmt"
	"testing"
)

func TestValidatorDigits(t *testing.T) {
	isValid := validate.IsDigits("11")
	fmt.Println("IsDigits:", isValid)
	if !isValid {
		t.Errorf("IsDigits Error")
	}
}

func TestValidatorBoolean(t *testing.T) {
	isValid := validate.IsBoolean("1")
	fmt.Println("IsBoolean:", isValid)
	if !isValid {
		t.Errorf("IsBoolean Error")
	}
}

func TestValidatorThaiIDNo(t *testing.T) {
	isValid := validate.IsThaiIDNo("1234567890123")
	fmt.Println("IsThaiIDNo:", isValid)
	if !isValid {
		t.Errorf("IsThaiIDNo Error")
	}
}

func TestValidatorMobileNo(t *testing.T) {
	isValid := validate.IsMobileNo("+66909096518")
	fmt.Println("IsMobileNo:", isValid)
	if !isValid {
		t.Errorf("IsMobileNo Error")
	}
}

func TestValidatorHomePhoneNo(t *testing.T) {
	isValid := validate.IsHomePhoneNo("027325113")
	fmt.Println("IsHomePhoneNo:", isValid)
	if !isValid {
		t.Errorf("IsHomePhoneNo Error")
	}
}

func TestValidatorEmail(t *testing.T) {
	isValid := validate.IsEmail("paravit.tun@gmail.com")
	fmt.Println("IsEmail:", isValid)
	if !isValid {
		t.Errorf("IsEmail Error")
	}
}

func TestValidatorHasStringValue(t *testing.T) {
	isValid := validate.HasStringValue("abc")
	fmt.Println("HasStringValue:", isValid)
	if !isValid {
		t.Errorf("HasStringValue Error")
	}
}
