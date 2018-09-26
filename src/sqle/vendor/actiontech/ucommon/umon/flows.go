package umon

type UmonFlows []UmonFlow
type UmonFlow struct {
	Tags         map[string]string
	FlowToCompId string
	Children     []UmonFlow
}

/*
	return
		1. first child if match
		2. this if no child match
		3. null if no match
*/
func (u *UmonFlow) Match(tags map[string]string) *UmonFlow {
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
