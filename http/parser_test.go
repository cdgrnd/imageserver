package http

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/pierrre/imageserver"
)

var _ Parser = ListParser{}

func TestListParserParse(t *testing.T) {
	parser := ListParser{
		&SourceParser{},
	}
	req, err := http.NewRequest("GET", "http://localhost?source=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	if !params.Has(imageserver.SourceParam) {
		t.Fatal("not set")
	}
}

func TestListParserParseError(t *testing.T) {
	parser := ListParser{
		&QualityParser{},
	}
	req, err := http.NewRequest("GET", "http://localhost?quality=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestListParserResolve(t *testing.T) {
	parser := ListParser{
		&SourceParser{},
	}

	httpParam := parser.Resolve(imageserver.SourceParam)
	if httpParam != imageserver.SourceParam {
		t.Fatal("not equals")
	}

	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &SourceParser{}

func TestSourceParserParse(t *testing.T) {
	parser := &SourceParser{}
	req, err := http.NewRequest("GET", "http://localhost?source=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	source, err := params.GetString(imageserver.SourceParam)
	if err != nil {
		t.Fatal(err)
	}
	if source != "foo" {
		t.Fatal("not equals")
	}
}

func TestSourceParserParseUndefined(t *testing.T) {
	parser := &SourceParser{}
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}

	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has(imageserver.SourceParam) {
		t.Fatal("should not be set")
	}
}

func TestSourceParserResolve(t *testing.T) {
	parser := &SourceParser{}
	httpParam := parser.Resolve(imageserver.SourceParam)
	if httpParam != imageserver.SourceParam {
		t.Fatal("not equals")
	}
	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &SourcePathParser{}

func TestSourcePathParserParse(t *testing.T) {
	parser := &SourcePathParser{}
	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	source, err := params.GetString(imageserver.SourceParam)
	if err != nil {
		t.Fatal(err)
	}
	if source != "/foobar" {
		t.Fatal("not equals")
	}
}

func TestSourcePathParserResolve(t *testing.T) {
	parser := &SourcePathParser{}
	httpParam := parser.Resolve(imageserver.SourceParam)
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &SourceURLParser{}

func TestSourceURLParserParse(t *testing.T) {
	storeURL, err := url.Parse("http://store/image")
	if err != nil {
		t.Fatal(err)
	}
	parser := &SourceURLParser{
		Parser: &SourcePathParser{},
		Base:   storeURL,
	}
	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	source, err := params.Get(imageserver.SourceParam)
	if err != nil {
		t.Fatal(err)
	}
	if source, ok := source.(*url.URL); !ok {
		t.Fatal("source is not an URL")
	} else {
		if source.String() != "http://store/image/foobar" {
			t.Fatal("not equals")
		}
	}
}

func TestSourceURLParserParseUndefined(t *testing.T) {
	storeURL, err := url.Parse("http://store/image")
	if err != nil {
		t.Fatal(err)
	}
	parser := &SourceURLParser{
		Parser: &SourceParser{},
		Base:   storeURL,
	}
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has(imageserver.SourceParam) {
		t.Fatal("should not be set")
	}
}

func TestSourceURLParserParseError(t *testing.T) {
	storeURL, err := url.Parse("http://store/image")
	if err != nil {
		t.Fatal(err)
	}
	parser := &SourceURLParser{
		Parser: &QualityParser{},
		Base:   storeURL,
	}
	req, err := http.NewRequest("GET", "http://localhost?quality=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSourceURLParserResolve(t *testing.T) {
	storeURL, err := url.Parse("http://store/image")
	if err != nil {
		t.Fatal(err)
	}
	parser := &SourceURLParser{
		Parser: &SourceParser{},
		Base:   storeURL,
	}
	httpParam := parser.Resolve(imageserver.SourceParam)
	if httpParam != imageserver.SourceParam {
		t.Fatal("not equals")
	}
}

var _ Parser = &FormatParser{}

func TestFormatParserParse(t *testing.T) {
	parser := &FormatParser{}
	req, err := http.NewRequest("GET", "http://localhost?format=jpg", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	format, err := params.GetString("format")
	if err != nil {
		t.Fatal(err)
	}
	if format != "jpeg" {
		t.Fatal("not equals")
	}
}

func TestFormatParserParseUndefined(t *testing.T) {
	parser := &FormatParser{}
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has("format") {
		t.Fatal("should not be set")
	}
}

func TestFormatParserParseError(t *testing.T) {
	parser := &FormatParser{}
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{"format": 666}
	err = parser.Parse(req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestFormatParserResolve(t *testing.T) {
	parser := &FormatParser{}

	httpParam := parser.Resolve("format")
	if httpParam != "format" {
		t.Fatal("not equals")
	}

	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &QualityParser{}

func TestQualityParserParse(t *testing.T) {
	parser := &QualityParser{}
	req, err := http.NewRequest("GET", "http://localhost?quality=50", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	quality, err := params.GetInt("quality")
	if err != nil {
		t.Fatal(err)
	}
	if quality != 50 {
		t.Fatal("not equals")
	}
}

func TestQualityParserParseError(t *testing.T) {
	parser := &QualityParser{}
	req, err := http.NewRequest("GET", "http://localhost?quality=foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err == nil {
		t.Fatal("no error")
	}
	if err, ok := err.(*imageserver.ParamError); !ok {
		t.Fatal("wrong error type")
	} else {
		param := err.Param
		if param != "quality" {
			t.Fatal("wrong param")
		}
	}
}

func TestQualityParserResolve(t *testing.T) {
	parser := &QualityParser{}
	httpParam := parser.Resolve("quality")
	if httpParam != "quality" {
		t.Fatal("not equals")
	}
	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &GammaCorrectionParser{}

func TestGammaCorrectionParserParse(t *testing.T) {
	parser := &GammaCorrectionParser{}
	req, err := http.NewRequest("GET", "http://localhost?gamma_correction=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	res, err := params.GetBool("gamma_correction")
	if err != nil {
		t.Fatal(err)
	}
	if res != true {
		t.Fatalf("unexpected result: got %t, want %t", res, true)
	}
}

func TestGammaCorrectionParserResolve(t *testing.T) {
	parser := &GammaCorrectionParser{}

	res := parser.Resolve("gamma_correction")
	expected := "gamma_correction"
	if res != expected {
		t.Fatalf("got %s, want %s", res, expected)
	}

	res = parser.Resolve("foobar")
	expected = ""
	if res != expected {
		t.Fatalf("got %s, want %s", res, expected)
	}
}

func TestParseQueryString(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?string=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	ParseQueryString("string", req, params)
	s, err := params.GetString("string")
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatal("not equals")
	}
}

func TestParseQueryStringUndefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	ParseQueryString("string", req, params)
	if params.Has("string") {
		t.Fatal("should not be set")
	}
}

func TestParseQueryInt(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?int=42", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryInt("int", req, params)
	if err != nil {
		t.Fatal(err)
	}
	i, err := params.GetInt("int")
	if err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Fatal("not equals")
	}
}

func TestParseQueryIntUndefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryInt("int", req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has("int") {
		t.Fatal("should not be set")
	}
}

func TestParseQueryIntError(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?int=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryInt("int", req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParseQueryFloat(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?float=12.34", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryFloat("float", req, params)
	if err != nil {
		t.Fatal(err)
	}
	f, err := params.GetFloat("float")
	if err != nil {
		t.Fatal(err)
	}
	if f != 12.34 {
		t.Fatal("not equals")
	}
}

func TestParseQueryFloatUndefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryFloat("float", req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has("float") {
		t.Fatal("should not be set")
	}
}

func TestParseQueryFloatError(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?float=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryFloat("float", req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParseQueryBool(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?bool=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryBool("bool", req, params)
	if err != nil {
		t.Fatal(err)
	}
	b, err := params.GetBool("bool")
	if err != nil {
		t.Fatal(err)
	}
	if b != true {
		t.Fatal("not equals")
	}
}

func TestParseQueryBoolUndefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryBool("bool", req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has("bool") {
		t.Fatal("should not be set")
	}
}

func TestParseQueryBoolError(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?bool=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryBool("bool", req, params)
	if err == nil {
		t.Fatal("no error")
	}
}
