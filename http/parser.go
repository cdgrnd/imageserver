package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pierrre/imageserver"
)

// Parser represents a HTTP Request parser.
type Parser interface {
	// Parse parses a Request and fill Params.
	Parse(*http.Request, imageserver.Params) error

	// Resolve resolves a param to a HTTP param.
	// It returns the resolved HTTP param, or an empty string.
	Resolve(param string) (httpParam string)
}

// ListParser is a list of HTTP Parser.
type ListParser []Parser

// Parse implements Parser.
func (lp ListParser) Parse(req *http.Request, params imageserver.Params) error {
	for _, subParser := range lp {
		err := subParser.Parse(req, params)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve implements Parser.
func (lp ListParser) Resolve(param string) string {
	for _, subParser := range lp {
		httpParam := subParser.Resolve(param)
		if httpParam != "" {
			return httpParam
		}
	}
	return ""
}

// SourceParser is a HTTP Parser that takes the source param from the query.
type SourceParser struct{}

// Parse implements Parser.
func (parser *SourceParser) Parse(req *http.Request, params imageserver.Params) error {
	ParseQueryString(imageserver.SourceParam, req, params)
	return nil
}

// Resolve implements Parser.
func (parser *SourceParser) Resolve(param string) string {
	if param != imageserver.SourceParam {
		return ""
	}
	return imageserver.SourceParam
}

// SourcePathParser is a HTTP Parser that takes the source param from the path.
type SourcePathParser struct {
}

// Parse implements Parser.
func (parser *SourcePathParser) Parse(req *http.Request, params imageserver.Params) error {
	params.Set(imageserver.SourceParam, req.URL.Path)
	return nil
}

// Resolve implements Parser.
func (parser *SourcePathParser) Resolve(param string) string {
	return ""
}

// SourceURLParser is a Parser that takes the source from the sub Parser and adds it to the Base URL.
type SourceURLParser struct {
	Parser
	Base *url.URL
}

// Parse implements Parser.
func (parser *SourceURLParser) Parse(req *http.Request, params imageserver.Params) error {
	err := parser.Parser.Parse(req, params)
	if err != nil {
		return err
	}
	if !params.Has(imageserver.SourceParam) {
		return nil
	}
	source, _ := params.Get(imageserver.SourceParam)
	u := copyURL(parser.Base)
	u.Path += fmt.Sprint(source)
	params.Set(imageserver.SourceParam, u)
	return nil
}

// Resolve implements Parser.
func (parser *SourceURLParser) Resolve(param string) string {
	return parser.Parser.Resolve(param)
}

// FormatParser is an http Parser that takes the "format" param from query.
type FormatParser struct {
}

// Parse implements Parser.
func (parser *FormatParser) Parse(req *http.Request, params imageserver.Params) error {
	ParseQueryString("format", req, params)
	if !params.Has("format") {
		return nil
	}
	format, err := params.GetString("format")
	if err != nil {
		return err
	}
	if format == "jpg" {
		format = "jpeg"
	}
	params.Set("format", format)
	return nil
}

// Resolve implements Parser.
func (parser *FormatParser) Resolve(param string) string {
	if param != "format" {
		return ""
	}
	return "format"
}

// QualityParser is an http Parser that takes the "quality" param from query.
type QualityParser struct {
}

// Parse implements Parser.
func (parser *QualityParser) Parse(req *http.Request, params imageserver.Params) error {
	return ParseQueryInt("quality", req, params)
}

// Resolve implements Parser.
func (parser *QualityParser) Resolve(param string) string {
	if param != "quality" {
		return ""
	}
	return "quality"
}

// GammaCorrectionParser is a HTTP Parser for gamma correct.
type GammaCorrectionParser struct{}

// Parse implements Parser.
func (parser *GammaCorrectionParser) Parse(req *http.Request, params imageserver.Params) error {
	return ParseQueryBool("gamma_correction", req, params)
}

// Resolve implements Parser.
func (parser *GammaCorrectionParser) Resolve(param string) string {
	if param == "gamma_correction" {
		return "gamma_correction"
	}
	return ""
}

// ParseQueryString takes the param from the query string and add it to params.
func ParseQueryString(param string, req *http.Request, params imageserver.Params) {
	s := req.URL.Query().Get(param)
	if s != "" {
		params.Set(param, s)
	}
}

// ParseQueryInt takes the param from the query string, parse it as an int and add it to params.
func ParseQueryInt(param string, req *http.Request, params imageserver.Params) error {
	s := req.URL.Query().Get(param)
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return newParseTypeParamError(param, "int", err)
	}
	params.Set(param, i)
	return nil
}

// ParseQueryFloat takes the param from the query string, parse it as a float64 and add it to params.
func ParseQueryFloat(param string, req *http.Request, params imageserver.Params) error {
	s := req.URL.Query().Get(param)
	if s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return newParseTypeParamError(param, "float", err)
	}
	params.Set(param, f)
	return nil
}

// ParseQueryBool takes the param from the query string, parse it as an bool and add it to params.
func ParseQueryBool(param string, req *http.Request, params imageserver.Params) error {
	s := req.URL.Query().Get(param)
	if s == "" {
		return nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return newParseTypeParamError(param, "bool", err)
	}
	params.Set(param, b)
	return nil
}

func newParseTypeParamError(param string, parseType string, parseErr error) *imageserver.ParamError {
	return &imageserver.ParamError{
		Param:   param,
		Message: fmt.Sprintf("parse %s: %s", parseType, parseErr.Error()),
	}
}
