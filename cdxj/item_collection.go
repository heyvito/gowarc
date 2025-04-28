package cdxj

type ItemCollection []*Item

func (i ItemCollection) Find(predicate func(*Item) bool) *Item {
	for _, item := range i {
		if predicate(item) {
			return item
		}
	}
	return nil
}

func (i ItemCollection) ByURL(url string) *Item {
	return i.Find(func(item *Item) bool { return item.URL == url })
}

func (i ItemCollection) AllByURL(url string) ItemCollection {
	var col ItemCollection
	for _, item := range i {
		if item.URL == url {
			col = append(col, item)
		}
	}
	return col
}

func (i ItemCollection) AllURLs() []string {
	var urls []string
	var set map[string]struct{}
	for _, item := range i {
		if _, ok := set[item.URL]; !ok {
			set[item.URL] = struct{}{}
			urls = append(urls, item.URL)
		}
	}
	return urls
}
