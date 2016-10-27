package myJson

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func PrintStructJson(value interface{}) (string, error) {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Struct {
		err := errors.New("The argument of printJson should be a struct")
		return "", err
	}
	resultBytes, err := printJSON(reflectValue, false)
	if err != nil {
		fmt.Println("err is: ", err.Error())
		return "", err
	}
	return string(resultBytes), err
}

func printBuiltin(value reflect.Value) ([]byte, error) {
	var b bytes.Buffer
L:
	switch value.Kind() {
	case reflect.Ptr:
		value = value.Elem()
		goto L
	case reflect.Bool:
		b.WriteString(strconv.FormatBool(value.Bool()))
	case reflect.String:
		mapVal := value.String()
		b.WriteString("\"")
		b.WriteString(mapVal)
		b.WriteString("\"")
	case reflect.Int32, reflect.Int64, reflect.Int, reflect.Int16, reflect.Int8, reflect.Uintptr:
		mapVal := value.Int()
		b.WriteString(strconv.FormatInt(int64(mapVal), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		mapVal := value.Uint()
		b.WriteString(strconv.FormatUint(uint64(mapVal), 10))
	case reflect.Float32, reflect.Float64:
		mapVal := value.Float()
		b.WriteString(strconv.FormatFloat(mapVal, 'g', 10, 64))
	default:
		return make([]byte, 0), errors.New("not builtin type")
	}
	return b.Bytes(), nil
}

func printJSON(value reflect.Value, anonymous bool) ([]byte, error) {
	var b bytes.Buffer

L:
	switch value.Kind() {
	case reflect.Ptr:
		value = value.Elem()
		goto L
	case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func:
		err := fmt.Errorf("unsurported kind: %s", value.Kind())
		return make([]byte, 0), err
	case reflect.Bool, reflect.String, reflect.Int32, reflect.Int64, reflect.Int, reflect.Int16, reflect.Int8, reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		reBytes, err := printBuiltin(value)
		if err != nil {
			return make([]byte, 0), err
		}
		b.Write(reBytes)
	case reflect.Struct:
		if !anonymous {
			b.WriteString("{")
		}
		for i := 0; i < value.NumField(); i++ {
			if i > 0 {
				b.WriteString(",")
			}
			f := getField(value, i)
			fAnonymous := value.Type().Field(i).Anonymous
			if !(f.Kind() == reflect.Struct && fAnonymous) {
				b.WriteString("\"")
				b.WriteString(value.Type().Field(i).Name)
				b.WriteString("\":")
			}
			val, err := printJSON(f, fAnonymous)
			if err != nil {
				return make([]byte, 0), err
			}
			b.Write(val)
		}
		if !anonymous {
			b.WriteString("}")
		}
	case reflect.Array, reflect.Slice:
		b.WriteString("[")
		for i := 0; i < value.Len(); i++ {
			if i > 0 {
				b.WriteString(",")
			}
			reBytes, err := printJSON(value.Index(i), false)
			if err != nil {
				return make([]byte, 0), err
			}
			b.Write(reBytes)
		}
		b.WriteString("]")
	case reflect.Map:
		b.WriteString("{")
		mapKeys := value.MapKeys()
		for _, key := range mapKeys {
			kBytes, err := printJSON(key, false)
			if err != nil {
				return make([]byte, 0), err
			}
			b.Write(kBytes)
			b.WriteString(":")

			reBytes, err := printJSON(value.MapIndex(key), false)
			if err != nil {
				return make([]byte, 0), err
			}
			b.Write(reBytes)
		}
		b.WriteString("}")
	}
	return b.Bytes(), nil
}

func getField(v reflect.Value, i int) reflect.Value {
	val := v.Field(i)
	if val.Kind() == reflect.Interface && !val.IsNil() {
		val = val.Elem()
	}
	return val
}
