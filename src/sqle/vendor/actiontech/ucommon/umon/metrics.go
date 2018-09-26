package umon

type UmonMetrics map[string]UmonMetric
type UmonMetric struct {
	Tags     map[string]string
	Enable   bool
	Children []UmonMetric
}

/*
	return
		1. first child if match
		2. this if no child match
		3. null if no match
*/
func (u *UmonMetric) Match(tags map[string]string) *UmonMetric {
	if nil != u.Tags {
		for tag, val := range u.Tags {
			if "" != tags[tag] && val != tags[tag] {
				return nil
			}
		}
	}
	if nil != u.Children {
		for _, child := range u.Children {
			if m := child.Match(tags); nil != m {
				return m
			}
		}
	}
	return u
}
