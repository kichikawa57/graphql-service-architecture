package model

import (
	"fmt"
	"io"
	"strconv"
)

// This file is just to create the package

type UserID struct {
	value int
}

func NewUserID(id int) UserID {
	return UserID{value: id}
}

func (u UserID) String() string {
	return fmt.Sprintf(`"%d"`, u.value)
}

func (u UserID) Int() int {
	return u.value
}

// MarshalGQL implements the graphql.Marshaler interface
func (u UserID) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Itoa(u.value))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (u *UserID) UnmarshalGQL(v any) error {
	switch v := v.(type) {
	case string:
		if result, err := strconv.Atoi(v); err != nil {
			return err
		} else {
			*u = NewUserID(result)
		}
		return nil
	case int:
		*u = NewUserID(v)
		return nil
	default:
		return fmt.Errorf("%T is not a url.URL", v)
	}
}