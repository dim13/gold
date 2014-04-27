package gold

import (
	"errors"
	"strings"
	"time"
)

var (
	TimeFormat = "January 2, 2006"
	readMore   = "<!--readmore-->"
)

type Articles []*Article
type Article struct {
	Date     time.Time
	Title    string
	Slug     string
	Body     string
	Tags     Tags
	Enabled  bool
	Author   string
	Comments Comments
}

type YearMap map[int]Articles
type MonthMap map[int]Articles

func (a *Article) makeSlug() {
	r := strings.NewReplacer(" ", "-")
	a.Slug = r.Replace(strings.TrimSpace(a.Title))
}

func (a *Article) Publish() {
	a.Date = time.Now()
	a.Enabled = true
}

func (a *Article) Suppress() {
	a.Enabled = false
}

func (a *Article) AddComment(c *Comment) {
	a.Comments.Add(c)
}

func (a Articles) Len() int           { return len(a) }
func (a Articles) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Articles) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func (a *Articles) Add(article *Article) error {
	article.Date = time.Now()
	if article.Slug == "" {
		article.makeSlug()
	}
	_, err := a.Find(article.Slug)
	if err == nil {
		return errors.New("duplicate slug " + article.Slug)
	}
	*a = append(*a, article)
	return nil
}

func (a Articles) Find(slug string) (*Article, error) {
	for i, _ := range a {
		if a[i].Slug == slug {
			return a[i], nil
		}
	}
	return &Article{}, errors.New("not found")
}

func (a Articles) Page(page, app int) (Articles, int, int) {
	var next, prev int

	lastpage := len(a)/app + 1

	if page <= 1 {
		page = 1
	} else {
		prev = page - 1
	}

	if page >= lastpage {
		page = lastpage
	} else {
		next = page + 1
	}

	from := (page - 1) * app
	to := from + app - 1
	if to > len(a) {
		to = len(a)
	}

	return a[from:to], next, prev
}

func (a Article) PostDate() string {
	return a.Date.Local().Format(TimeFormat)
}

func (a Article) RssDate() string {
	return a.Date.Local().Format(time.RFC1123Z)
}

func (a Article) ReadMore() string {
	if i := strings.Index(a.Body, readMore); i > 0 {
		return a.Body[:i]
	}
	return a.Body
}

func (a Article) HasMore() bool {
	return strings.Contains(a.Body, readMore)
}

func (a Articles) Year(year int) (A Articles) {
	for _, v := range a {
		if v.Date.Year() == year {
			A = append(A, v)
		}
	}
	return A
}

func (a Articles) Month(month int) (A Articles) {
	for _, v := range a {
		if v.Date.Month() == time.Month(month) {
			A = append(A, v)
		}
	}
	return A
}

func (a Articles) YearMonth(year, month int) (A Articles) {
	for _, v := range a {
		if v.Date.Year() == year &&
			v.Date.Month() == time.Month(month) &&
			v.Enabled {
			A = append(A, v)
		}
	}
	return A
}

func (a Articles) Enabled() (A Articles) {
	for _, v := range a {
		if v.Enabled {
			A = append(A, v)
		}
	}
	return A
}

func (a Articles) YearMap() YearMap {
	ym := make(YearMap)
	for _, v := range a.Enabled() {
		y := v.Date.Year()
		ym[y] = append(ym[y], v)
	}
	return ym
}

func (a Articles) MonthMap() MonthMap {
	mm := make(MonthMap)
	for _, v := range a.Enabled() {
		m := int(v.Date.Month())
		mm[m] = append(mm[m], v)
	}
	return mm
}

