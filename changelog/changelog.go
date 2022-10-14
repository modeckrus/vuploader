package changelog

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

//Parse CHANGELOG.md

type ChangeLog []Change

// 1 -> 100
// 12 -> 120
// 123 -> 123
func littleVersion(v string) int64 {
	i, err := strconv.ParseInt(v, 10, 64)
	if i == 300 {
		fmt.Println(i)
	}
	if err != nil {
		return 0
	}
	if len(v) == 1 {
		return i * 100
	}
	if len(v) == 2 {
		return i * 100
	}
	if len(v) == 3 {
		return i * 100
	}
	return i
}

// 1.2.3+456 -> 001002003456
// 1.2.30+456 -> 001002030456
// 1.2.300+456 -> 001002300456
// 1.20.3+456 -> 001020003456
// 1.20.30+456 -> 001020030456
// 1.20.300+456 -> 001020300456
// 1.200.3+456 -> 001200003456
// 1.200.30+456 -> 001200030456
func Parse(version string) int64 {
	result := int64(0)
	buildSplitted := strings.Split(string(version), "+")
	version = buildSplitted[0]
	buildVersionStr := buildSplitted[1]

	v := strings.Split(string(version), ".")
	if len(v) != 3 {
		return 0
	}
	buildVersionInt, _ := strconv.ParseInt(buildVersionStr, 10, 64)
	result = littleVersion(v[0])*10000000 + littleVersion(v[1])*10000 + littleVersion(v[2])*10 + buildVersionInt
	return result
}

type Change struct {
	Version  string   `json:"version"`
	Futures  []string `json:"futures"`
	Bugfixes []string `json:"bugfixes"`
	Message  string   `json:"message"`
}

func (i Change) String() string {
	result := fmt.Sprintf("<i><b>%s</b></i>", i.Version)
	if len(i.Futures) > 0 {
		result += fmt.Sprintf("\n<b>Новое:</b>\n<i>%s</i>", strings.Join(i.Futures, "\n"))
	}
	if len(i.Bugfixes) > 0 {
		result += fmt.Sprintf("\n<b>Багфиксы:</b>\n<i>%s</i>", strings.Join(i.Bugfixes, "\n"))
	}
	if i.Message != "" {
		result += fmt.Sprintf("\n<b>Описание:</b>\n<i>%s</i>", i.Message)
	}
	return result
}

type SortByVersion []Change

func (a SortByVersion) Len() int      { return len(a) }
func (a SortByVersion) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByVersion) Less(i, j int) bool {
	return ParseVersion(a[i].Version) < ParseVersion(a[j].Version)
}
func ParseVersion(v string) int64 {
	result := int64(0)
	buildSplitted := strings.Split(string(v), "+")
	v = buildSplitted[0]
	buildVersionStr := buildSplitted[1]

	splitted := strings.Split(string(v), ".")
	for i, item := range splitted {
		var index = (2 - i)
		multiplier := int64(math.Pow(1000, float64(index)))
		// multiplier *= 1000
		key := littleVersion(item)
		// key, _ := strconv.ParseInt(item, 10, 64)
		result += key * multiplier
	}
	buildVersion, err := strconv.ParseInt(buildVersionStr, 10, 64)
	if err != nil {
		return result
	}
	result += buildVersion
	return result
}

func NewChangelog(path string) (ChangeLog, error) {
	result := ChangeLog{}
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}
	err = json.NewDecoder(file).Decode(&result)
	sort.Sort(SortByVersion(result))
	return result, err
}

func (i ChangeLog) last() Change {
	return i[len(i)-1]
}
func (i ChangeLog) NumberVersion() string {
	lastVersion := i.last()

	return string(lastVersion.Version)
}
func (i ChangeLog) LastVersion() string {

	lastVersion := i.last()

	return lastVersion.String()
}
