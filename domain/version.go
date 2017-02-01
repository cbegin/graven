package domain

import (
	"strings"
	"fmt"
	"strconv"
)

type Version struct {
	Major     int
	Minor     int
	Patch     int
	Qualifier string
}

func (v *Version) Parse(s string) (err error) {
	parts := strings.Split(s, "-")
	var full string
	if len(parts) == 2 {
		full = parts[0]
		v.Qualifier = parts[1]
	} else if len(parts) == 1 {
		full = parts[0]
	} else {
		return fmt.Errorf("Too many qualifier delimiters (-) in version string: %s", s)
	}
	parts = strings.Split(full, ".")
	if len(parts) == 3 {
		if v.Major, err = strconv.Atoi(parts[0]); err != nil {
			return err
		}
		if v.Minor, err = strconv.Atoi(parts[1]); err != nil {
			return err
		}
		if v.Patch, err = strconv.Atoi(parts[2]); err != nil {
			return err
		}
	} else if len(parts) == 2 {
		if v.Major, err = strconv.Atoi(parts[0]); err != nil {
			return err
		}
		if v.Minor, err = strconv.Atoi(parts[1]); err != nil {
			return err
		}
	} else if len(parts) == 1 {
		if v.Major, err = strconv.Atoi(parts[0]); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Too many qualifier delimiters (-) in version string: %s", s)
	}
	return err

}

func (v *Version) ToString() string {
	if len(v.Qualifier) > 0 {
		return fmt.Sprintf("%v.%v.%v-%v", v.Major, v.Minor, v.Patch, v.Qualifier)
	} else {
		return fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, v.Patch)
	}
}

