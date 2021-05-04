package models

import "awesomeProject/pkg/util"

type Tags struct {
	Tags []string `json:"tags"`
}

type TagSummary struct {
	Tag         string   `json:"tag"`
	Count       int      `json:"count"`
	Articles    []int    `json:"articles"`
	RelatedTags []string `json:"related_tags"`
}

func (r *Articles) GetTagSummary(tag string) TagSummary {
	var as []int
	rt := make(util.StringKeyMap)

	for k, a := range *r {
		as = append(as, k)
		for _, t := range a.Tags {
			if t != tag {
				rt[t] = 1
			}
		}
	}

	ts := TagSummary{
		Tag:         tag,
		Count:       len(*r),
		Articles:    as,
		RelatedTags: rt.Keys(),
	}

	return ts
}
