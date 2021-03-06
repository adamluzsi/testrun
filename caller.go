package testcase

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

var testcasePkgDirPath string

func init() {
	_, specFilePath, _, _ := runtime.Caller(0)
	testcasePkgDirPath = path.Dir(specFilePath)
}

func (spec *Spec) callerLocationName(skip int) string {
	locationName := func(file string, line int) string {
		return fmt.Sprintf(`%s:%d`, path.Base(file), line)
	}
	for i := 0; i < 1024; i++ {
		_, file, line, ok := runtime.Caller(1 + skip + i) // 1 means skip this file
		if !ok {
			return ""
		}
		// fast path when caller located in a *_test.go file
		if strings.HasSuffix(file, `_test.go`) {
			return locationName(file, line)
		}
		// skip testcase packages
		if strings.HasPrefix(file, testcasePkgDirPath) {
			continue
		}
		// skip stdlib testing
		if strings.Contains(file, `go/src/testing/`) {
			continue
		}
		// skip stdlib runtime
		if strings.Contains(file, `go/src/runtime/`) {
			continue
		}
		return locationName(file, line)
	}
	return ""
}
