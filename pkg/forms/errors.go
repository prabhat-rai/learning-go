package forms

type errors map[string][]string

// Implement an Add() method to add error messages for a given field to the map.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Implement a Get() method to retrieve the first error message for a given field from the map.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

// Implement a GetAllErrors() method to retrieve all error messages
func (e errors) GetAllErrors() [] string {
	var errs []string

	for i :=range e {
		errs = append(errs, i + ":" + e[i][0])
	}

	return errs
}