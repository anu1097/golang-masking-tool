package mask

import (
	"fmt"
	"testing"
	"time"

	"github.com/anu1097/golang-masking-tool/customMasker"
	"github.com/anu1097/golang-masking-tool/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilValuePassed(t *testing.T) {
	//for this one test the filter does not matter
	maskTool := NewMaskTool(filter.AllFieldFilter())
	masked := maskTool.MaskDetails(nil)
	require.Nil(t, masked)
}

func TestValueFilter(t *testing.T) {
	t.Run("DefaultValueFilter", func(t *testing.T) {
		const issuedToken = "abcd1234"
		maskTool := NewMaskTool(filter.ValueFilter(issuedToken))
		t.Run("string", func(t *testing.T) {
			record := "Authorization: Bearer " + issuedToken
			filteredData := maskTool.MaskDetails(record)
			require.NotNil(t, filteredData)
			assert.Equal(t, "Authorization: Bearer [filtered]", filteredData)

			// fmt.Println(filteredData)
			// "Authorization: Bearer [filtered]"
		})

		t.Run("struct", func(t *testing.T) {
			type myRecord struct {
				ID   string
				Data string
			}
			record := myRecord{
				ID:   "userId",
				Data: issuedToken,
			}

			filteredData := maskTool.MaskDetails(record)
			require.NotNil(t, filteredData)
			copied, ok := filteredData.(myRecord)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, "userId", copied.ID)
			assert.Equal(t, filter.GetFilteredLabel(), copied.Data)
			// fmt.Println(copied)
			// "{userId [filtered]}"
		})

		t.Run("array", func(t *testing.T) {
			record := []string{
				"userId",
				"data",
				issuedToken,
			}

			filteredData := maskTool.MaskDetails(record)
			require.NotNil(t, filteredData)
			assert.Equal(t, []string([]string{"userId", "data", filter.GetFilteredLabel()}), filteredData)
			// fmt.Println(copied)
			// "{userId [filtered]}"
		})

		t.Run("map", func(*testing.T) {
			mapRecord := map[string]interface{}{
				"data": issuedToken,
			}
			filteredData := maskTool.MaskDetails(mapRecord)
			require.NotNil(t, filteredData)
			assert.Equal(t, map[string]interface{}{"data": "[filtered]"}, filteredData)
		})

	})
	t.Run("CustomValueFilter", func(t *testing.T) {
		const issuedToken = "abcd1234"
		maskTool := NewMaskTool(filter.CustomValueFilter(issuedToken, customMasker.MPassword))
		t.Run("string", func(t *testing.T) {
			record := "Authorization: Bearer " + issuedToken
			filteredData := maskTool.MaskDetails(record)
			require.NotNil(t, filteredData)
			assert.Equal(t, "Authorization: Bearer ************", filteredData)

			// fmt.Println(filteredData)
			// "Authorization: Bearer [filtered]"
		})

		t.Run("struct", func(t *testing.T) {
			type myRecord struct {
				ID   string
				Data string
			}
			record := myRecord{
				ID:   "userId",
				Data: issuedToken,
			}

			filteredData := maskTool.MaskDetails(record)
			require.NotNil(t, filteredData)
			copied, ok := filteredData.(myRecord)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, "userId", copied.ID)
			assert.Equal(t, "************", copied.Data)
			// fmt.Println(copied)
			// "{userId [filtered]}"
		})

		t.Run("array", func(t *testing.T) {
			record := []string{
				"userId",
				"data",
				issuedToken,
			}

			filteredData := maskTool.MaskDetails(record)
			require.NotNil(t, filteredData)
			assert.Equal(t, []string([]string{"userId", "data", "************"}), filteredData)
			// fmt.Println(copied)
			// "{userId [filtered]}"
		})

		t.Run("map", func(*testing.T) {
			mapRecord := map[string]interface{}{
				"data": issuedToken,
			}
			filteredData := maskTool.MaskDetails(mapRecord)
			require.NotNil(t, filteredData)
			assert.Equal(t, map[string]interface{}{"data": "************"}, filteredData)
		})

	})
}

func TestVariousDatastructuresForVariousScenarios(t *testing.T) {
	customMasker := NewMaskingInstance(
		filter.ValueFilter("blue"),
	)

	t.Run("struct", func(t *testing.T) {
		type testData struct {
			ID    int
			Name  string
			Label string
		}

		t.Run("original data is not modified when filtered", func(t *testing.T) {
			data := &testData{
				ID:    100,
				Name:  "blue",
				Label: "five",
			}
			v := customMasker.MaskDetails(data)
			require.NotNil(t, v)
			copied, ok := v.(*testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, filter.GetFilteredLabel(), copied.Name)
			assert.Equal(t, "blue", data.Name)
			assert.Equal(t, "five", data.Label)
			assert.Equal(t, "five", copied.Label)
			assert.Equal(t, 100, copied.ID)
		})

		t.Run("non-ptr struct can be modified", func(t *testing.T) {
			data := testData{
				Name:  "blue",
				Label: "five",
			}
			v := customMasker.MaskDetails(data)
			require.NotNil(t, v)
			copied, ok := v.(testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, filter.GetFilteredLabel(), copied.Name)
			assert.Equal(t, "five", copied.Label)
		})

		t.Run("nested structure can be modified", func(t *testing.T) {
			type testDataParent struct {
				Child testData
			}

			data := &testDataParent{
				Child: testData{
					Name:  "blue",
					Label: "five",
				},
			}
			v := customMasker.MaskDetails(data)
			require.NotNil(t, v)
			copied, ok := v.(*testDataParent)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, filter.GetFilteredLabel(), copied.Child.Name)
			assert.Equal(t, "five", copied.Child.Label)
		})

		t.Run("unexported field is also copied", func(t *testing.T) {
			type myStruct struct {
				unexported string
				Exported   string
			}

			data := &myStruct{
				unexported: "red",
				Exported:   "orange",
			}
			v := customMasker.MaskDetails(data)
			require.NotNil(t, v)
			copied, ok := v.(*myStruct)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, "red", data.unexported)
			assert.Equal(t, "orange", data.Exported)
		})

		t.Run("original type", func(t *testing.T) {
			type myType string
			type myData struct {
				Name myType
			}
			data := &myData{
				Name: "miss blue",
			}
			v := customMasker.MaskDetails(data)
			require.NotNil(t, v)
			copied, ok := v.(*myData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, myType("miss "+filter.GetFilteredLabel()), copied.Name)
		})

		t.Run("various field", func(t *testing.T) {
			type child struct{}
			type myStruct struct {
				Func      func() time.Time
				Chan      chan int
				Bool      bool
				Bytes     []byte
				Interface interface{}
				Child     *child
			}
			data := &myStruct{
				Func:  time.Now,
				Chan:  make(chan int),
				Bool:  true,
				Bytes: []byte("timeless"),
				Child: nil,
			}
			v := customMasker.MaskDetails(data)
			require.NotNil(t, v)
			copied, ok := v.(*myStruct)
			require.True(t, ok)
			require.NotNil(t, copied)

			// function type is not compareable, but it's ok if not nil
			assert.NotNil(t, copied.Func)
			assert.Equal(t, data.Chan, copied.Chan)
			assert.Equal(t, data.Bool, copied.Bool)
			assert.Equal(t, data.Bytes, copied.Bytes)
		})

		t.Run("map data", func(t *testing.T) {
			data := map[string]*testData{
				"xyz": {
					Name:  "blue",
					Label: "five",
				},
			}
			v := customMasker.MaskDetails(data)
			require.NotNil(t, v)
			copied, ok := v.(map[string]*testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, filter.GetFilteredLabel(), copied["xyz"].Name)
			assert.Equal(t, "five", copied["xyz"].Label)
		})

		t.Run("array data", func(t *testing.T) {
			data := []testData{
				{
					Name:  "orange",
					Label: "five",
				},
				{
					Name:  "blue",
					Label: "five",
				},
			}
			v := customMasker.MaskDetails(data)
			require.NotNil(t, v)
			copied, ok := v.([]testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, "orange", copied[0].Name)
			assert.Equal(t, filter.GetFilteredLabel(), copied[1].Name)
			assert.Equal(t, "five", copied[1].Label)
		})

		t.Run("array data with ptr", func(t *testing.T) {
			data := []*testData{
				{
					Name:  "orange",
					Label: "five",
				},
				{
					Name:  "blue",
					Label: "five",
				},
			}
			v := customMasker.MaskDetails(data)
			require.NotNil(t, v)
			copied, ok := v.([]*testData)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, "orange", copied[0].Name)
			assert.Equal(t, filter.GetFilteredLabel(), copied[1].Name)
			assert.Equal(t, "five", copied[1].Label)
		})

	})

}

func TestAllFieldFilter(t *testing.T) {
	type ID string
	type child struct {
		Data string
	}
	s := "test"
	var id ID = "id"
	type myStruct struct {
		Func      func() time.Time
		Chan      chan int
		Bool      bool
		Bytes     []byte
		Strs      []string
		StrsPtr   []*string
		Interface interface{}
		Child     child
		ChildPtr  *child
		Data      string
		Str       string
		Pstr      *string
		ID        ID
		PID       *ID
	}
	data := &myStruct{
		Func:      time.Now,
		Chan:      make(chan int),
		Bool:      true,
		Bytes:     []byte("timeless"),
		Strs:      []string{"aa"},
		StrsPtr:   []*string{&s},
		Interface: &s,
		Child:     child{Data: "x"},
		ChildPtr:  &child{Data: "y"},
		Data:      "data",
		Str:       s,
		Pstr:      &s,
		ID:        id,
		PID:       &id,
	}

	t.Run("default allfield filter", func(t *testing.T) {
		mask := NewMaskingInstance(
			filter.AllFieldFilter(),
		)

		v := mask.MaskDetails(data)
		require.NotNil(t, v)
		copied, ok := v.(*myStruct)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Nil(t, copied.Func)
		assert.Nil(t, copied.Chan)
		assert.Nil(t, copied.Bytes)
		assert.Nil(t, copied.Strs)
		assert.Nil(t, copied.StrsPtr)
		assert.Nil(t, copied.Interface)
		assert.Empty(t, copied.Child.Data)
		assert.Empty(t, copied.ChildPtr.Data)
		assert.Equal(t, filter.GetFilteredLabel(), copied.Data)
		assert.Equal(t, ("test"), data.Str)
		assert.Equal(t, filter.GetFilteredLabel(), copied.Str)
		assert.Equal(t, ("test"), *data.Pstr)
		assert.Equal(t, filter.GetFilteredLabel(), *copied.Pstr)
		assert.Equal(t, ID("id"), data.ID)
		assert.Equal(t, ID(filter.GetFilteredLabel()), copied.ID)
		assert.Equal(t, ID("id"), *data.PID)
		assert.Equal(t, ID(filter.GetFilteredLabel()), *copied.PID)
	})

	t.Run("custom allfield filter", func(t *testing.T) {
		mask := NewMaskingInstance(
			filter.CustomAllFieldFilter(customMasker.MPassword),
		)

		v := mask.MaskDetails(data)
		require.NotNil(t, v)
		copied, ok := v.(*myStruct)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Nil(t, copied.Func)
		assert.Nil(t, copied.Chan)
		assert.Nil(t, copied.Bytes)
		assert.Nil(t, copied.Strs)
		assert.Nil(t, copied.StrsPtr)
		assert.Nil(t, copied.Interface)
		assert.Empty(t, copied.Child.Data)
		assert.Empty(t, copied.ChildPtr.Data)
		assert.Equal(t, "************", copied.Data)
		assert.Equal(t, ("test"), data.Str)
		assert.Equal(t, "************", copied.Str)
		assert.Equal(t, ("test"), *data.Pstr)
		assert.Equal(t, "************", *copied.Pstr)
		assert.Equal(t, ID("id"), data.ID)
		assert.Equal(t, ID("************"), copied.ID)
		assert.Equal(t, ID("id"), *data.PID)
		assert.Equal(t, ID("************"), *copied.PID)
	})

}

func TestTypeFilter(t *testing.T) {
	t.Run("CustomTypeFilter", func(t *testing.T) {

		type password string
		type myRecord struct {
			ID       string
			Password password
		}
		record := myRecord{
			ID:       "userId",
			Password: "abcd1234",
		}

		t.Run("Type Filter with Mask Type", func(t *testing.T) {
			maskTool := NewMaskTool(filter.CustomTypeFilter(password(""), customMasker.MPassword))
			filteredData := maskTool.MaskDetails(record)
			require.NotNil(t, filteredData)
			copied, ok := filteredData.(myRecord)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, password("************"), copied.Password)
			assert.Equal(t, "userId", copied.ID)

		})

		// fmt.Println(copied)
		// {userId [filtered]}
	})

	t.Run("TypeFilter", func(t *testing.T) {
		type password string
		type myRecord struct {
			ID       string
			Password password
		}
		record := myRecord{
			ID:       "userId",
			Password: "abcd1234",
		}

		t.Run("Default Type Filter", func(t *testing.T) {
			maskTool := NewMaskTool(filter.TypeFilter(password("")))
			filteredData := maskTool.MaskDetails(record)
			require.NotNil(t, filteredData)
			copied, ok := filteredData.(myRecord)
			require.True(t, ok)
			require.NotNil(t, copied)
			assert.Equal(t, password(filter.GetFilteredLabel()), copied.Password)
			assert.Equal(t, "userId", copied.ID)
		})

	})
}

func TestTagFilter(t *testing.T) {
	t.Run("default ", func(t *testing.T) {
		type myRecord struct {
			ID    string
			EMail string `mask:"secret"`
		}
		record := myRecord{
			ID:    "userId",
			EMail: "dummy@dummy.com",
		}

		maskTool := NewMaskTool(filter.TagFilter())
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, filter.GetFilteredLabel(), copied.EMail)
		assert.Equal(t, "userId", copied.ID)

		// fmt.Println(copied)
		// {userId [filtered]}
	})

	t.Run("custom ", func(t *testing.T) {
		type myRecord struct {
			ID    string
			EMail string `mask:"email"`
			Phone string `mask:"mobile"`
		}
		record := myRecord{
			ID:    "userId",
			EMail: "dummy@dummy.com",
			Phone: "9191919191",
		}

		maskTool := NewMaskTool(filter.TagFilter(customMasker.MEmail, customMasker.MMobile))
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, "dum****@dummy.com", copied.EMail)
		assert.Equal(t, "9191***191", copied.Phone)
		assert.Equal(t, "userId", copied.ID)

		// fmt.Println(copied)
		// {userId [filtered]}
	})
}

func TestPiiPhoneNumber(t *testing.T) {
	maskTool := NewMaskTool(filter.PhoneFilter())

	t.Run("string", func(t *testing.T) {
		stringRecord := "090-0000-0000"
		filteredData := maskTool.MaskDetails(stringRecord)
		require.NotNil(t, filteredData)
		assert.Equal(t, filter.GetFilteredLabel(), filteredData)

		// fmt.Println(filteredData)
		// [filtered]
	})

	t.Run("struct", func(t *testing.T) {
		type myRecord struct {
			ID    string
			Phone string
		}
		record := myRecord{
			ID:    "userId",
			Phone: "090-0000-0000",
		}
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, filter.GetFilteredLabel(), copied.Phone)
		assert.Equal(t, "userId", copied.ID)

		// fmt.Println(copied)
		// {userId [filtered]}
	})

	t.Run("map", func(*testing.T) {
		mapRecord := map[string]interface{}{
			"phone": "090-0000-0000",
		}
		filteredData := maskTool.MaskDetails(mapRecord)
		require.NotNil(t, filteredData)
		assert.Equal(t, map[string]interface{}{"phone": "[filtered]"}, filteredData)
	})

}

func TestCustomPiiPhoneNumber(t *testing.T) {
	maskTool := NewMaskTool(filter.CustomPhoneFilter(customMasker.MMobile))

	t.Run("string", func(t *testing.T) {
		stringRecord := "090-0000-0000"
		filteredData := maskTool.MaskDetails(stringRecord)
		require.NotNil(t, filteredData)
		assert.Equal(t, "090-***0-0000", filteredData)

		// fmt.Println(filteredData)
		// [filtered]
	})

	t.Run("struct", func(t *testing.T) {
		type myRecord struct {
			ID    string
			Phone string
		}
		record := myRecord{
			ID:    "userId",
			Phone: "090-0000-0000",
		}
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, "090-***0-0000", copied.Phone)
		assert.Equal(t, "userId", copied.ID)

		// fmt.Println(copied)
		// {userId [filtered]}
	})

	t.Run("map", func(*testing.T) {
		mapRecord := map[string]interface{}{
			"phone": "090-0000-0000",
		}
		filteredData := maskTool.MaskDetails(mapRecord)
		require.NotNil(t, filteredData)
		assert.Equal(t, map[string]interface{}{"phone": "090-***0-0000"}, filteredData)
	})

}
func TestPiiEmail(t *testing.T) {
	type myRecord struct {
		ID    string
		Email string
	}
	record := myRecord{
		ID:    "userId",
		Email: "dummy@dummy.com",
	}
	maskTool := NewMaskTool(filter.EmailFilter())
	filteredData := maskTool.MaskDetails(record)
	require.NotNil(t, filteredData)
	copied, ok := filteredData.(myRecord)
	require.True(t, ok)
	require.NotNil(t, copied)
	assert.Equal(t, filter.GetFilteredLabel(), copied.Email)
	assert.Equal(t, "userId", copied.ID)

	// fmt.Println(copied)
	// {userId [filtered]}
}

func TestCustomPiiEmail(t *testing.T) {
	type myRecord struct {
		ID    string
		Email string
	}
	record := myRecord{
		ID:    "userId",
		Email: "dummy@dummy.com",
	}
	maskTool := NewMaskTool(filter.CustomEmailFilter(customMasker.MEmail))
	filteredData := maskTool.MaskDetails(record)
	require.NotNil(t, filteredData)
	copied, ok := filteredData.(myRecord)
	require.True(t, ok)
	require.NotNil(t, copied)
	assert.Equal(t, "dum****@dummy.com", copied.Email)
	assert.Equal(t, "userId", copied.ID)

	// fmt.Println(copied)
	// {userId [filtered]}
}

func TestPiiCustomRegexNumber(t *testing.T) {

	customRegex := "^https:\\/\\/(dummy-backend.)[0-9a-z]*.com\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)$"
	t.Run("Default Filter", func(t *testing.T) {
		type myRecord struct {
			ID   string
			Link string
		}
		record := myRecord{
			ID:   "userId",
			Link: "https://dummy-backend.dummy.com/v2/random",
		}
		maskTool := NewMaskTool(filter.CustomRegexFilter(customRegex))
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, filter.GetFilteredLabel(), copied.Link)
		assert.Equal(t, "userId", copied.ID)
		// fmt.Println(copied)
		// {userId [filtered]}
	})

	t.Run("Custom Filter", func(t *testing.T) {
		type myRecord struct {
			ID   string
			Link string
		}
		record := myRecord{
			ID:   "userId",
			Link: "https://dummy-backend.dummy.com/v2/random",
		}
		maskTool := NewMaskTool(filter.CustomRegexFilterWithMType(customRegex, customMasker.MPassword))
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, "************", copied.Link)
		assert.Equal(t, "userId", copied.ID)
	})

}

func TestFieldFilter(t *testing.T) {
	t.Run("struct with no specific mask", func(*testing.T) {
		type myRecord struct {
			ID    string
			Phone string
		}
		record := myRecord{
			ID:    "userId",
			Phone: "090-0000-0000",
		}

		maskTool := NewMaskTool(filter.FieldFilter("Phone"))
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, filter.GetFilteredLabel(), copied.Phone)
		assert.Equal(t, "userId", copied.ID)

		// fmt.Println(copied)
		// {userId [filtered]}

	})
}

func TestCustomFieldFilter(t *testing.T) {
	t.Run("struct", func(*testing.T) {
		type myRecord struct {
			ID         string
			Phone      string
			Url        string
			Email      string
			Name       string
			Address    string
			CreditCard string
		}
		record := myRecord{
			ID:         "userId",
			Phone:      "090-0000-0000",
			Url:        "http://admin:mysecretpassword@localhost:1234/uri",
			Email:      "dummy@dummy.com",
			Name:       "John Doe",
			Address:    "1 AB Road, Paradise",
			CreditCard: "4444-4444-4444-4444",
		}

		maskTool := NewMaskTool(
			filter.CustomFieldFilter("Phone", customMasker.MMobile),
			filter.CustomFieldFilter("Email", customMasker.MEmail),
			filter.CustomFieldFilter("Url", customMasker.MURL),
			filter.CustomFieldFilter("Name", customMasker.MName),
			filter.CustomFieldFilter("ID", customMasker.MID),
			filter.CustomFieldFilter("Address", customMasker.MAddress),
			filter.CustomFieldFilter("CreditCard", customMasker.MCreditCard),
		)
		maskTool.UpdateCustomMaskingChar(customMasker.PCross)
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, "090-xxx0-0000", copied.Phone)
		assert.Equal(t, "dumxxxx@dummy.com", copied.Email)
		assert.Equal(t, "userIdxxxx", copied.ID)
		assert.Equal(t, "Jxxn Dxxe", copied.Name)
		assert.Equal(t, "1 AB Rxxxxxx", copied.Address)
		assert.Equal(t, "4444-4xxxxxx44-4444", copied.CreditCard)
		assert.Equal(t, "http://admin:xxxxx@localhost:1234/uri", copied.Url)

		fmt.Println(copied)
		// {userId [filtered]}

	})

	t.Run("map", func(*testing.T) {
		mapRecord := map[string]interface{}{
			"secret": "secretData",
		}
		filter := NewMaskTool(filter.CustomFieldFilter("secret", customMasker.MEmail))
		filteredData := filter.MaskDetails(mapRecord)
		require.NotNil(t, filteredData)
		assert.Equal(t, map[string]interface{}(map[string]interface{}{"secret": "secretData"}), mapRecord)
		assert.Equal(t, map[string]interface{}(map[string]interface{}{"secret": interface{}(nil)}), filteredData)
	})

}

func TestFieldPrefixFilter(t *testing.T) {
	type myRecord struct {
		ID          string
		SecurePhone string
	}
	t.Run("default", func(*testing.T) {
		record := myRecord{
			ID:          "userId",
			SecurePhone: "090-0000-0000",
		}

		maskTool := NewMaskTool(filter.FieldPrefixFilter("Secure"))
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, filter.GetFilteredLabel(), copied.SecurePhone)
		assert.Equal(t, "userId", copied.ID)

		// fmt.Println(copied)
		// {userId [filtered]}
	})

}

func TestCustomFieldPrefixFilter(t *testing.T) {
	type myRecord struct {
		ID          string
		SecurePhone string
	}

	t.Run("custom", func(*testing.T) {
		record := myRecord{
			ID:          "userId",
			SecurePhone: "090-0000-0000",
		}

		maskTool := NewMaskTool(filter.CustomFieldPrefixFilter("Secure", customMasker.MMobile))
		filteredData := maskTool.MaskDetails(record)
		require.NotNil(t, filteredData)
		copied, ok := filteredData.(myRecord)
		require.True(t, ok)
		require.NotNil(t, copied)
		assert.Equal(t, "090-***0-0000", copied.SecurePhone)
		assert.Equal(t, "userId", copied.ID)
	})

}
