package mask

import (
	"reflect"

	"github.com/anu1097/golang-masking-tool/customMasker"
	"github.com/anu1097/golang-masking-tool/filter"
)

// Get a new masking instance. Pass your required filters
//
// Example:
//
//	var maskingInstance = NewMaskTool(filter.FieldFilter("Phone"))
func NewMaskTool(filters ...filter.Filter) masking {
	return *NewMaskingInstance(filters...)
}

type Masking interface {
	// Call to update masking character for custom masker
	UpdateCustomMaskingChar(maskingChar customMasker.MaskingCharacter)

	// Call to update filter label
	UpdateFilterLabel(filterlabel string)

	// Call to get filter label
	GetFilterLabel() string

	// Append to existing list of filters in masking instance
	AppendFilters(filters ...filter.Filter)

	// Get complete list of existing filters used by masking instance
	GetFilters() filter.Filters

	// Call to Mask Details from a given instance
	MaskDetails(v interface{}) interface{}

	// Internal function which masks based on filters and returns a clone of the data passed
	clone(fieldName string, value reflect.Value, tag string) reflect.Value
}

type masking struct {
	filterList filter.Filters
}

// Get a pointer to new masking instance. Pass your required filters
//
// Example:
//
//	var maskingInstance = NewMaskingInstance(filter.FieldFilter("Phone"))
func NewMaskingInstance(filters ...filter.Filter) *masking {
	filter.SetCustomMaskerInstance(customMasker.NewMasker())
	var filterList = filter.Filters{}
	filterList = append(filterList, filters...)
	return &masking{
		filterList: filterList,
	}
}

func (x *masking) UpdateCustomMaskingChar(maskingChar customMasker.MaskingCharacter) {
	filter.UpdateCustomMaskingChar(maskingChar)
}

func (x *masking) UpdateFilterLabel(filterlabel string) {
	filter.SetFilteredLabel(filterlabel)
}

func (x *masking) GetFilterLabel() string {
	return filter.GetFilteredLabel()
}

func (x *masking) AppendFilters(filters ...filter.Filter) {
	x.filterList = append(x.filterList, filters...)
}

func (x *masking) GetFilters() filter.Filters {
	return x.filterList
}

func (x *masking) MaskDetails(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	return x.clone("", reflect.ValueOf(v), "").Interface()
}

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
			if fv.Type().Kind() == reflect.Ptr && fv.Elem().Kind() == reflect.String {
				a := x.clone(f.Name, fv.Elem(), tagValue).Interface().(string)
				dst.Elem().Field(i).Set(reflect.ValueOf(&a))
			} else {
				dst.Elem().Field(i).Set(x.clone(f.Name, fv, tagValue))
			}
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
