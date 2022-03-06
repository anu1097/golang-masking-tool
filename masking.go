package mask

import (
	"reflect"

	"github.com/anu1097/golang-mask-utility/customMasker"
	"github.com/anu1097/golang-mask-utility/filter"
)

type Masking interface {
	UpdateCustomMaskingChar(maskingChar customMasker.MaskingCharacter)
	UpdateFilterLabel(filterlabel string)
	GetFilterLabel() string
	AppendFilters(filters ...filter.Filter)
	GetFilters() filter.Filters
	MaskDetails(v interface{}) interface{}
	clone(fieldName string, value reflect.Value, tag string) reflect.Value
}

type masking struct {
	filterList filter.Filters
}

// Get a new masking instance. Pass your required filters
func NewMaskingInstance(filters ...filter.Filter) *masking {
	filter.SetCustomMaskerInstance(customMasker.NewMasker())
	var filterList = filter.Filters{}
	filterList = append(filterList, filters...)
	return &masking{
		filterList: filterList,
	}
}

// Call to update masking character for custom masker
func (x *masking) UpdateCustomMaskingChar(maskingChar customMasker.MaskingCharacter) {
	filter.UpdateCustomMaskingChar(maskingChar)
}

// Call to update filter label
func (x *masking) UpdateFilterLabel(filterlabel string) {
	filter.SetFilteredLabel(filterlabel)
}

// Call to get filter label
func (x *masking) GetFilterLabel() string {
	return filter.GetFilteredLabel()
}

// Append to existing list of filters in masking instance
func (x *masking) AppendFilters(filters ...filter.Filter) {
	x.filterList = append(x.filterList, filters...)
}

// Get complete list of existing filters used by masking instance
func (x *masking) GetFilters() filter.Filters {
	return x.filterList
}

// Call to Mask Details from a given instance
func (x *masking) MaskDetails(v interface{}) interface{} {
	return x.clone("", reflect.ValueOf(v), "").Interface()
}

// Internal function which masks based on filters and returns a clone of the data passed
func (x *masking) clone(fieldName string, value reflect.Value, tag string) reflect.Value {
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
			dst.Elem().Field(i).Set(x.clone(f.Name, fv, tagValue))
		}

	case reflect.Map:
		dst = reflect.MakeMap(src.Type())
		keys := src.MapKeys()
		for i := 0; i < src.Len(); i++ {
			mValue := src.MapIndex(keys[i])
			dst.SetMapIndex(keys[i], x.clone(keys[i].String(), mValue, ""))
		}

	case reflect.Array, reflect.Slice:
		dst = reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		for i := 0; i < src.Len(); i++ {
			dst.Index(i).Set(x.clone(fieldName, src.Index(i), ""))
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
