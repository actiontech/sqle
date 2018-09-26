package log

//utility functions provide to external

func KeyRet(stage *Stage, err error) error {
	if nil == err {
		Key(stage, "succeed")
	} else {
		Key(stage, "return error (%v)", err)
	}
	return err
}

func BriefRet(stage *Stage, err error) error {
	if nil == err {
		Brief(stage, "succeed")
	} else {
		Brief(stage, "return error (%v)", err)
	}
	return err
}
