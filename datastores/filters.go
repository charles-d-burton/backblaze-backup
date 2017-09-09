package datastores

import (
	"log"
	"regexp"
)

var filters = make([]*regexp.Regexp, 0, 1)

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

func GetFilters() []*regexp.Regexp {
	return filters
}
