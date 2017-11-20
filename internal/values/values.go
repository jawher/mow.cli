package values

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/jawher/mow.cli/internal/container"
)

// BoolValued is an interface values can implement to indicate that they are a bool option, i.e. can be set without providing a value with just -f for example
type BoolValued interface {
	flag.Value
	// IsBoolFlag should return true to indicate that this value is a bool value
	IsBoolFlag() bool
}

// EnumValued is an interface ti indicate that a value can hold an enum
type EnumValued interface {
	flag.Value
	// Validate makes sure the value passed is valid
	Validate(string, []container.Validator) (string, error)
}

// MultiValued is an interface ti indicate that a value can hold multiple values
type MultiValued interface {
	flag.Value
	// Clear should clear the list of values
	Clear()
}

// DefaultValued in an interface to determine if the value stored is the default value, and thus does not need be shown in the help message
type DefaultValued interface {
	// IsDefault should return true if the value stored is the default value, and thus does not need be shown in the help message
	IsDefault() bool
}

/******************************************************************************/
/* BOOL                                                                        */
/******************************************************************************/

// BoolValue is a flag.Value type holding boolean values
type BoolValue bool

var (
	_ flag.Value    = NewBool(new(bool), false)
	_ BoolValued    = NewBool(new(bool), false)
	_ DefaultValued = NewBool(new(bool), false)
)

// NewBool creates a new bool value
func NewBool(into *bool, v bool) *BoolValue {
	*into = v
	return (*BoolValue)(into)
}

// Set sets the value from a provided string
func (bo *BoolValue) Set(s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*bo = BoolValue(b)
	return nil
}

// IsBoolFlag returns true
func (bo *BoolValue) IsBoolFlag() bool {
	return true
}

func (bo *BoolValue) String() string {
	return fmt.Sprintf("%v", *bo)
}

// IsDefault return true if the bool value is false
func (bo *BoolValue) IsDefault() bool {
	return !bool(*bo)
}

/******************************************************************************/
/* ENUM                                                                        */
/******************************************************************************/

// EnumValue is a flag.Value type holding string values
type EnumValue string

var (
	_ flag.Value    = NewEnum(new(string), "")
	_ EnumValued    = NewEnum(new(string), "")
	_ DefaultValued = NewEnum(new(string), "")
)

// NewEnum creates a new string value
func NewEnum(into *string, v string) *EnumValue {
	*into = v
	return (*EnumValue)(into)
}

// Set sets the value from a provided string
func (sa *EnumValue) Set(s string) error {
	*sa = EnumValue(s)
	return nil
}

func (sa *EnumValue) String() string {
	return fmt.Sprintf("%#v", *sa)
}

// IsDefault return true if the string value is empty
func (sa *EnumValue) IsDefault() bool {
	return string(*sa) == ""
}

// Validate validates the specified enum value, returns the wanted value and
// calls the relevant callback if defined.
func (sa *EnumValue) Validate(v string, vv []container.Validator) (string, error) {
	for _, x := range vv {
		if x.User == v {
			if x.Callback != nil {
				err := x.Callback()
				if err != nil {
					return "", err
				}
			}
			return x.Value, nil
		}
	}

	// If we are here the value is invalid, let's give the user the list of
	// valid values in our validation list.
	help := ""
	for i := 0; i < len(vv); i++ {
		help += vv[i].User
		if i == (len(vv) - 1) {
			help += "."
		} else {
			help += ", "
		}
	}
	return "", fmt.Errorf("Invalid value %s, valid values are %s", v, help)
}

/******************************************************************************/
/* STRING                                                                        */
/******************************************************************************/

// StringValue is a flag.Value type holding string values
type StringValue string

var (
	_ flag.Value    = NewString(new(string), "")
	_ DefaultValued = NewString(new(string), "")
)

// NewString creates a new string value
func NewString(into *string, v string) *StringValue {
	*into = v
	return (*StringValue)(into)
}

// Set sets the value from a provided string
func (sa *StringValue) Set(s string) error {
	*sa = StringValue(s)
	return nil
}

func (sa *StringValue) String() string {
	return fmt.Sprintf("%#v", *sa)
}

// IsDefault return true if the string value is empty
func (sa *StringValue) IsDefault() bool {
	return string(*sa) == ""
}

/******************************************************************************/
/* INT                                                                        */
/******************************************************************************/

// IntValue is a flag.Value type holding int values
type IntValue int

var (
	_ flag.Value = NewInt(new(int), 0)
)

// NewInt creates a new int value
func NewInt(into *int, v int) *IntValue {
	*into = v
	return (*IntValue)(into)
}

// Set sets the value from a provided string
func (ia *IntValue) Set(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*ia = IntValue(int(i))
	return nil
}

func (ia *IntValue) String() string {
	return fmt.Sprintf("%v", *ia)
}

/******************************************************************************/
/* STRINGS                                                                    */
/******************************************************************************/

// StringsValue is a flag.Value type holding string slices values
type StringsValue []string

var (
	_ flag.Value    = NewStrings(new([]string), nil)
	_ MultiValued   = NewStrings(new([]string), nil)
	_ DefaultValued = NewStrings(new([]string), nil)
)

// NewStrings creates a new multi-string value
func NewStrings(into *[]string, v []string) *StringsValue {
	*into = v
	return (*StringsValue)(into)
}

// Set sets the value from a provided string
func (sa *StringsValue) Set(s string) error {
	*sa = append(*sa, s)
	return nil
}

func (sa *StringsValue) String() string {
	res := "["
	for idx, s := range *sa {
		if idx > 0 {
			res += ", "
		}
		res += fmt.Sprintf("%#v", s)
	}
	return res + "]"
}

// Clear clears the slice
func (sa *StringsValue) Clear() {
	*sa = nil
}

// IsDefault return true if the string slice is empty
func (sa *StringsValue) IsDefault() bool {
	return len(*sa) == 0
}

/******************************************************************************/
/* INTS                                                                       */
/******************************************************************************/

// IntsValue is a flag.Value type holding int values
type IntsValue []int

var (
	_ flag.Value    = NewInts(new([]int), nil)
	_ MultiValued   = NewInts(new([]int), nil)
	_ DefaultValued = NewInts(new([]int), nil)
)

// NewInts creates a new multi-int value
func NewInts(into *[]int, v []int) *IntsValue {
	*into = v
	return (*IntsValue)(into)
}

// Set sets the value from a provided string
func (ia *IntsValue) Set(s string) error {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*ia = append(*ia, int(i))
	return nil
}

func (ia *IntsValue) String() string {
	res := "["
	for idx, s := range *ia {
		if idx > 0 {
			res += ", "
		}
		res += fmt.Sprintf("%v", s)
	}
	return res + "]"
}

// Clear clears the slice
func (ia *IntsValue) Clear() {
	*ia = nil
}

// IsDefault return true if the int slice is empty
func (ia *IntsValue) IsDefault() bool {
	return len(*ia) == 0
}
