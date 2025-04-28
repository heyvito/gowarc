package cdxj

import "strings"

func Parse(data []byte) (ItemCollection, error) {
	lines := strings.Split(string(data), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	var col ItemCollection
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		it, err := ItemFromString(line)
		if err != nil {
			return nil, err
		}
		col = append(col, it)
	}
	return col, nil
}
