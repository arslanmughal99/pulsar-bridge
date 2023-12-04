package main

import (
	"github.com/gookit/validate"
)

func ValidateDto(dto any) map[string]string {
	v := validate.Struct(dto)
	v.StopOnError = false

	e := v.ValidateE()

	if len(e) <= 0 {
		return nil
	}

	errs := map[string]string{}

	for k, v := range e {
		for _, r := range v {
			errs[k] = r
		}
	}

	return errs
}
