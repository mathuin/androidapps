package main

import (
	"fmt"
	"reflect"
	"testing"
)

var NewSetting_tests = []struct {
	desc   string
	envvar string
	result *Setting
	err    error
}{
	{"Port", "ANDROIDAPPS_PORT", &Setting{description: "Port", envvar: "ANDROIDAPPS_PORT"}, nil},
	{"Port", "", nil, fmt.Errorf("missing envvar")},
	{"", "ANDROIDAPPS_PORT", nil, fmt.Errorf("missing description")},
	{"", "", nil, fmt.Errorf("missing description")},
}

func Test_NewSetting_1(t *testing.T) {
	for _, tt := range NewSetting_tests {
		output, err := NewSetting(tt.desc, tt.envvar)
		if reflect.DeepEqual(output, tt.result) == false {
			t.Errorf("Given desc=%+#v and envvar=%+#v, wanted (%+#v, %+#v), got (%+#v, %+#v) instead", tt.desc, tt.envvar, tt.result, tt.err, output, err)
		}
	}
}

// test set_value
// - what happens with no envvar and no flag value
//   (error)
// - what happens with no envvar and flag value
//   (value set to flag_value)
// - what happens with envvar and no flag value
//   (value set to envvar)
// - what happens with envvar and flag value
//   (value set to envvar)

// test apply_settings
// - create a settings variable of two with good settings
//   silent success
// - create a settings variable of two with one bad setting
//   log fatal
//   (return error?)
