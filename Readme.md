# Golang Masking Tool
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/anu1097/golang-masking-tool/blob/main/LICENSE)
![GolangVersion](https://img.shields.io/github/go-mod/go-version/anu1097/golang-masking-tool)
[![CircleCI](https://circleci.com/gh/anu1097/golang-masking-tool/tree/main.svg?style=svg)](https://circleci.com/gh/anu1097/golang-masking-tool/tree/main)
[![Release](https://img.shields.io/github/v/release/anu1097/golang-masking-tool)](https://github.com/anu1097/golang-masking-tool/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/anu1097/golang-masking-tool.svg)](https://pkg.go.dev/github.com/anu1097/golang-masking-tool)

Golang Masking Tool is a simple utility of creating a masker tool which you can use to mask sensitive information.
You can use a variety of filters with custom masking types to assist you.

Inspired by two repositories - 
1. [zlog](https://github.com/m-mizutani/zlog)
2. [Golang Masker](https://github.com/ggwhite/go-masker)

Both libraries were solving similar usecase but didn't cover all use cases I was looking for. Zlog is a libray that focuses more on logging and its filterating features are not exposed to be used separately. While Golang Masker doesn't cover all data types. Hence it didn't solve my use-cases.

So I combined both of them to create this library and sharing the best of both them. Added some more uses cases too.

* [Getting Started](#Getting-Started)

# Getting Started

```
$ go get -u github.com/anu1097/golang-masking-tool
```


## Usage

- [Basic example](#basic-example)
	- [Creating a Masking Instance](#create-masking-instance)
- [Filter sensitive data](#filter-sensitive-data)
    - [By specified field](#by-specified-field)
    - [By specified field-prefix](#by-specified-field-prefix)
	- [By specified value](#by-specified-value)
	- [By custom type](#by-custom-type)
	- [By struct tag](#by-struct-tag)
	- [By data pattern (e.g. personal information)](#by-regex-pattern)
    - [All Fields Filter](#by-allfields-filter)
- [Customise Masking Tool](#customise-masking-tool)
	- [Update Custom Masker Character](#update-custom-masker-character)
	- [Update Default Filter](#update-default-filter)
	- [Append More Filters](#append-more-filter)

## Basic Example

### Create Masking Instance

```
    var maskingInstance = NewMaskTool()
```

## Filter Sensitive Data
### By Specified Field

Default 

```
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

	// fmt.Println(filteredData)
	// {userId [filtered]}
```

Custom Mask Type

```
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
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	// {userId**** 090-***0-0000 http://admin:xxxxx@localhost:1234/uri dum****@dummy.com J**n D**e 1 AB R****** 4444-4******44-4444}
```


### By Specified Field-Prefix

Default
```
	type myRecord struct {
		ID          string
		SecurePhone string
	}
	record := myRecord{
		ID:          "userId",
		SecurePhone: "090-0000-0000",
	}

	maskTool := NewMaskTool(filter.FieldPrefixFilter("Secure"))
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	// {userId [filtered]}
```

Custom 
```
	maskTool := NewMaskTool(filter.CustomFieldPrefixFilter("Secure", customMasker.MMobile))
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	// {userId 090-***0-0000}

```
### By Specified Value

Default

```
	const issuedToken = "abcd1234"
	maskTool := NewMaskTool(filter.ValueFilter(issuedToken))
	record := "Authorization: Bearer " + issuedToken
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	// Authorization: Bearer [filtered]
```

Custom Mask Type

```
	const issuedToken = "abcd1234"
	maskTool := NewMaskTool(filter.CustomValueFilter(issuedToken, customMasker.MPassword))
	record := "Authorization: Bearer " + issuedToken
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	// Authorization: Bearer ************
```
### By custom type

Default

```
	type password string
	type myRecord struct {
		ID       string
		Password password
	}
	record := myRecord{
		ID:       "userId",
		Password: "abcd1234",
	}
	maskTool := NewMaskTool(filter.CustomTypeFilter(password(""), customMasker.MPassword))
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	// {userId [filtered]}
```

Custom Mask Type

```
	type password string
	type myRecord struct {
		ID       string
		Password password
	}
	record := myRecord{
		ID:       "userId",
		Password: "abcd1234",
	}

	maskTool := NewMaskTool(filter.CustomTypeFilter(password(""), customMasker.MPassword))
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	// {userId ************}
```
### By struct tag

Default

```
	type myRecord struct {
		ID    string
		EMail string `mask:"secret"` //Use secret for default filter
	}
	record := myRecord{
		ID:    "userId",
		EMail: "dummy@dummy.com",
	}

	maskTool := NewMaskTool(filter.TagFilter())
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	// {userId [filtered]}
```

Custom Mask Type

```
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

	// fmt.Println(filteredData)
	// {userId dum****@dummy.com 9191***191}
```
### By Regex Pattern

Default Phone Filter

```
	type myRecord struct {
		ID    string
		Phone string
	}
	record := myRecord{
		ID:    "userId",
		Phone: "090-0000-0000",
	}
	maskTool := NewMaskTool(filter.PhoneFilter())
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	// {userId [filtered]}
```

Custom Phone Filter

```
	type myRecord struct {
		ID    string
		Phone string
	}
	record := myRecord{
		ID:    "userId",
		Phone: "090-0000-0000",
	}
	maskTool := NewMaskTool((filter.CustomPhoneFilter(customMasker.MMobile)))
	filteredData := maskTool.MaskDetails(record)

	// fmt.Println(filteredData)
	//{userId 090-***0-0000}
```

Default Email Filter

```
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

	// fmt.Println(filteredData)
	// {userId [filtered]}
```

Custom Email Filter
```

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

	// fmt.Println(filteredData)
	// {userId dum****@dummy.com}
```

Custom Regex Filter
```
	customRegex := "^https:\\/\\/(dummy-backend.)[0-9a-z]*.com\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)$"
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

	// fmt.Println(filteredData)
	// {userId [filtered]}
```

Custom Regex Filter With Mask Type
```
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

	// fmt.Println(filteredData)
	// {userId ************}
```

### By AllFields Filter

Default
```
	type child struct {
		Data string
	}
	s := "test"
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
	}
	mask := NewMaskingInstance(
		filter.AllFieldFilter(),
	)

	filteredData := mask.MaskDetails(data)

	// fmt.Println(filteredData)
	// {<nil> <nil> false [] [] [] <nil> {} 0x1400009b180 [filtered]}

```

With Custom Mask Type
```
	mask := NewMaskingInstance(
		filter.CustomAllFieldFilter(customMasker.MPassword),
	)
	filteredData := mask.MaskDetails(data)

	// fmt.Println(filteredData)
	// {<nil> <nil> false [] [] [] <nil> {} 0x1400009b180 ************}

```
## Custom Mask Types

|Type        |Const        |Tag        |Description                                                                                            |
|:----------:|:-----------:|:---------:|:------------------------------------------------------------------------------------------------------|
|Name        |MName        |name       |mask the second letter and the third letter                                                            |
|Password    |MPassword    |password   |always return `************`                                                                           |
|Address     |MAddress     |addr       |keep first 6 letters, mask the rest                                                                    |
|Email       |MEmail       |email      |keep domain and the first 3 letters                                                                    |
|Mobile      |MMobile      |mobile     |mask 3 digits from the 4'th digit                                                                      |
|Telephone   |MTelephone   |tel        |remove `(`, `)`, ` `, `-` chart, and mask last 4 digits of telephone number, format to `(??)????-????` |
|ID          |MID          |id         |mask last 4 digits of ID number                                                                        |
|CreditCard  |MCreditCard  |credit     |mask 6 digits from the 7'th digit                                                                      |
|Secret      |MStruct      |secret     |Uses default filtered string. Use only with Struct Tag filter                                                                                        |


## Customise Masking Tool

### Update Default Filter
```
	maskTool := NewMaskTool(filter.FieldFilter("Phone"))
	maskTool.SetFilteredLabel("CustomFilterString")
	// maskTool.GetFilteredLabel()
    // CustomFilterString
```

### Update Custom Masker Character
```
	maskTool := NewMaskTool(filter.FieldFilter("Phone"))
	maskTool.UpdateCustomMaskingChar(customMasker.PCross)
```
### Append More Filter
```
	maskTool := NewMaskTool(filter.FieldFilter("Phone"))
	maskTool.AppendFilters(filter.EmailFilter())
```
## License

- MIT License
- Author: Anuraag Gupta <anuraagg.kval3@gmail.com>
