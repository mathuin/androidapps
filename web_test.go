package main

import (
	"html/template"
	"testing"
)

var obfuscate_tests = []struct {
	input  string
	output template.HTML
}{
	{"pie", "&#x70;&#x69;&#x65;"},
}

func Test_obfuscate(t *testing.T) {
	for _, tt := range obfuscate_tests {
		output := obfuscate(tt.input)
		if output != tt.output {
			t.Errorf("Given %+#v, wanted %+#v, got %+#v", tt.input, tt.output, output)
		}
	}
}

var mailto_tests = []struct {
	input  string
	output template.HTMLAttr
}{
	{"pie", "href=\"&#x6d;&#x61;&#x69;&#x6c;&#x74;&#x6f;&#x3a;&#x70;&#x69;&#x65;\""},
}

func Test_mailto(t *testing.T) {
	for _, tt := range mailto_tests {
		output := mailto(tt.input)
		if output != tt.output {
			t.Errorf("Given %+#v, wanted %+#v, got %+#v", tt.input, tt.output, output)
		}
	}
}
