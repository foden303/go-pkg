package err

import "fmt"

func (i *impl) Error() string {
	if i.err != nil {
		return fmt.Sprintf("%s: %s. Cause: %v", i.ErrorCode, i.ErrorMessage, i.err)
	}
	return fmt.Sprintf("%s: %s", i.ErrorCode, i.ErrorMessage)
}

func (i *impl) Err() error {
	return i.err
}
