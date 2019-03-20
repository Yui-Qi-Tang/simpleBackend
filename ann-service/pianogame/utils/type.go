package utils

import "reflect"

// GetTypeName return the type of value as string
func GetTypeName(value interface{}) string {
	valueOf := reflect.ValueOf(value)

	if valueOf.Type().Kind() == reflect.Ptr {
		return reflect.Indirect(valueOf).Type().Name()
	}

	return (valueOf.Type().Name())
}
