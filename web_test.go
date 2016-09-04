package main

import (
	"html/template"
	"testing"
)

var obfuscateTests = []struct {
	input  string
	output template.HTML
}{
	{"pie", "&#x70;&#x69;&#x65;"},
}

func Test_obfuscate(t *testing.T) {
	for _, tt := range obfuscateTests {
		output := obfuscate(tt.input)
		if output != tt.output {
			t.Errorf("Given %+#v, wanted %+#v, got %+#v", tt.input, tt.output, output)
		}
	}
}

var mailtoTests = []struct {
	input  string
	output template.HTMLAttr
}{
	{"pie", "href=\"&#x6d;&#x61;&#x69;&#x6c;&#x74;&#x6f;&#x3a;&#x70;&#x69;&#x65;\""},
}

func Test_mailto(t *testing.T) {
	for _, tt := range mailtoTests {
		output := mailto(tt.input)
		if output != tt.output {
			t.Errorf("Given %+#v, wanted %+#v, got %+#v", tt.input, tt.output, output)
		}
	}
}
