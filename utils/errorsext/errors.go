package errorsext

import "strings"

type MultiError struct {
	errs []error
}

func (me MultiError) Error() string {
	ss := make([]string, len(me.errs))
	for i := range me.errs {
		ss[i] = me.errs[i].Error()
	}
	return strings.Join(ss, "\n")
}

type MultiErrorBuilder struct {
	me MultiError
}

func (meb *MultiErrorBuilder) PushIfErr(err error) {
	if err != nil {
		meb.Push(err)
	}
}

func (meb *MultiErrorBuilder) Push(err error) {
	meb.me.errs = append(meb.me.errs, err)
}

func (meb *MultiErrorBuilder) Error() error {
	if len(meb.me.errs) == 0 {
		return nil
	}
	return meb.me
}
