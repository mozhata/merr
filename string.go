package merr

import "fmt"

// BuildErrMsg used to format error message
func BuildErrMsg(msgs ...interface{}) string {
	if len(msgs) == 0 || msgs == nil {
		return ""
	}
	if len(msgs) == 1 {
		if v, ok := msgs[0].(string); ok {
			return v
		}
		if v, ok := msgs[0].(error); ok {
			return v.Error()
		}
	}
	if len(msgs) > 1 {
		return fmt.Sprintf(msgs[0].(string), msgs[1:]...)
	}
	return ""
}
