package filter

var loggerFilters Filters = Filters{}

func AddFilters(sampleFilter ...Filter) {
	loggerFilters = append(loggerFilters, sampleFilter...)
}

func GetFilters() Filters {
	return loggerFilters
}
