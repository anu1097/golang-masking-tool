package mask

import (
	"reflect"

	"github.com/anu1097/golang-mask-utility/filter"
	"github.com/anu1097/golang-mask-utility/masker"
)

var filterList filter.Filters = filter.Filters{}

type Masking interface {
	AddFilters(filters ...filter.Filter)
	GetFilters() filter.Filters
	MaskDetails(v interface{}) interface{}
	Clone(fieldName string, value reflect.Value, tag string) reflect.Value
	UpdateMaskingCharacter()
}

type masking struct {
	filterList filter.Filters
}

func NewMasking(filters ...filter.Filter) *masking {
	filterList = append(filterList, filters...)
	return &masking{
		filterList: filterList,
	}
}

func (x *masking) UpdateCustomMaskingChar(maskingChar masker.MaskinCharacter) {
	filter.UpdateCustomMaskingChar(maskingChar)
}
func (x *masking) AddFilters(filters ...filter.Filter) {
	x.filterList = append(x.filterList, filters...)
}

func (x *masking) GetFilters() filter.Filters {
	return x.filterList
}

func (x *masking) MaskDetails(v interface{}) interface{} {
	return x.Clone("", reflect.ValueOf(v), "").Interface()
}

func (x *masking) Clone(fieldName string, value reflect.Value, tag string) reflect.Value {
	adjustValue := func(ret reflect.Value) reflect.Value {
		switch value.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Array:
			return ret
		default:
			return ret.Elem()
		}
	}

	src := value
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return reflect.New(value.Type()).Elem()
		}
		src = value.Elem()
	}

	var dst reflect.Value
	maskingFilter, shouldMask := filter.CheckShouldMask(x.filterList, fieldName, src.Interface(), tag)
	if shouldMask {
		dst = reflect.New(src.Type())
		switch src.Kind() {
		case reflect.String:
			filteredData := maskingFilter.MaskString(value.String())
			dst.Elem().SetString(filteredData)
		case reflect.Array, reflect.Slice:
			dst = dst.Elem()
		}
		return adjustValue(dst)
	}

	switch src.Kind() {
	case reflect.String:
		dst = reflect.New(src.Type())
		filtered := x.filterList.ReplaceString(value.String())
		dst.Elem().SetString(filtered)

	case reflect.Struct:
		dst = reflect.New(src.Type())
		t := src.Type()

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fv := src.Field(i)
			if !fv.CanInterface() {
				continue
			}
			tagValue := f.Tag.Get(filter.GetTagKey())
			dst.Elem().Field(i).Set(x.Clone(f.Name, fv, tagValue))
		}

	case reflect.Map:
		dst = reflect.MakeMap(src.Type())
		keys := src.MapKeys()
		for i := 0; i < src.Len(); i++ {
			mValue := src.MapIndex(keys[i])
			dst.SetMapIndex(keys[i], x.Clone(keys[i].String(), mValue, ""))
		}

	case reflect.Array, reflect.Slice:
		dst = reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		for i := 0; i < src.Len(); i++ {
			dst.Index(i).Set(x.Clone(fieldName, src.Index(i), ""))
		}

	case reflect.Interface:
		dst = reflect.New(src.Type())
		data := value.Interface()
		stringData, ok := data.(string)
		if !ok {
			dst.Elem().Set(src)
		} else {
			filtered := x.filterList.ReplaceString(stringData)
			dst.Elem().Set(reflect.ValueOf(filtered))
		}

	default:
		dst = reflect.New(src.Type())
		dst.Elem().Set(src)
	}
	return adjustValue(dst)
}
