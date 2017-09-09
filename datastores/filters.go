package datastores

import (
	"log"
	"regexp"
)

var filters = make([]*regexp.Regexp, 0, 1)

//SetFilters ...Compile the configured filters and add them to the array
func SetFilters(newFilters []string) {
	for _, filter := range newFilters {
		log.Println("Filter: ", filter)
		regex, err := regexp.Compile(filter)
		if err != nil {
			log.Println("Compile error: ", err)
		}
		filters = append(filters, regex)
	}
}

//GetFilters ...Retrieve the list of filters
func GetFilters() []*regexp.Regexp {
	return filters
}
