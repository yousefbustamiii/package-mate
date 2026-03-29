package components

import "strings"

// Resolve finds an InstallItem by any of its CLI alias names or display name.
// Returns the item, its parent section, and whether it was found.
func Resolve(alias string) (*InstallItem, *Section, bool) {
	needle := strings.ToLower(strings.TrimSpace(alias))

	// Priority 1: Check display name (case-insensitive)
	for si := range AllSections {
		for ii := range AllSections[si].Items {
			item := &AllSections[si].Items[ii]
			if strings.EqualFold(item.Name, needle) {
				return item, &AllSections[si], true
			}
		}
	}

	// Priority 2: Check canonical binary (case-insensitive)
	for si := range AllSections {
		for ii := range AllSections[si].Items {
			item := &AllSections[si].Items[ii]
			if item.Binary != "" && strings.EqualFold(item.Binary, needle) {
				return item, &AllSections[si], true
			}
		}
	}

	return nil, nil, false
}
