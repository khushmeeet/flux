package fluxgen

import "sort"

type By func(p1, p2 *page) bool

type pageSort struct {
	pages pages
	by    By
}

func (by By) Sort(pages pages) {
	ps := &pageSort{
		pages: pages,
		by:    by,
	}
	sort.Sort(ps)
}

func (p pageSort) Len() int {
	return len(p.pages)
}

func (p pageSort) Swap(i, j int) {
	p.pages[i], p.pages[j] = p.pages[j], p.pages[i]
}

func (p pageSort) Less(i, j int) bool {
	return p.by(&p.pages[i], &p.pages[j])
}

func descendingOrderByDate(p1, p2 *page) bool {
	return p1.Date.After(p2.Date)
}
