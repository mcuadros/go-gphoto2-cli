package gphoto2cli

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	autoDetechRegExp = regexp.MustCompile("([a-zA-Z0-9_ ]+)(usb:(\\d+),(\\d+))")
	errorMsgRegExp   = regexp.MustCompile(`\*\*\* Error \*\*\* *\n(.+)\n\*\*\* Error \(([-0-9]+): '([^']+)'\) \*\*\*`)
)

type parser struct{}

func (g *parser) ParseAutoDetect(output string) []*Port {
	m := autoDetechRegExp.FindAllStringSubmatch(output, -1)
	if m == nil {
		return nil
	}

	var result []*Port
	for _, match := range m {
		result = append(result, &Port{
			Model: strings.TrimSpace(string(match[1])),
			Port:  string(match[2]),
		})
	}

	return result
}

func (p *parser) ParseCameraSummary(text string) *CameraSummary {
	c := &CameraSummary{}
	c.Manufacturer = p.extractValue(text, "Manufacturer")
	c.Model = p.extractValue(text, "Model")
	c.Version = p.extractValue(text, "Version")
	c.SerialNumber = p.extractValue(text, "Serial Number")
	c.CaptureFormats = p.extractList(text, "Capture Formats")
	c.DisplayFormats = p.extractList(text, "Display Formats")

	return c
}

func (p *parser) extractValue(text, key string) string {
	re := regexp.MustCompile(key + `:\s*(.*?)\n`)
	matches := re.FindStringSubmatch(text)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

func (p *parser) extractList(text, key string) []string {
	re := regexp.MustCompile(key + `:\s*(.*?)\n`)
	matches := re.FindStringSubmatch(text)
	if len(matches) == 2 {
		return strings.Split(matches[1], ", ")
	}
	return nil
}

func (p *parser) ParseGetConfig(output string) *ConfigValue {
	lines := strings.Split(output, "\n")
	v := &ConfigValue{}

	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]
		switch key {
		case "Label":
			v.Label = value
		case "Readonly":
			v.Readonly = value == "1"
		case "Type":
			v.Type = value
		case "Current":
			v.Current = value
		case "Choice":
			if v.Type == "RADIO" {
				v.Choices = strings.Split(value, " ")
			}
		}
	}

	return v
}

func (p *parser) ParseErrorMessage(input string) error {
	matches := errorMsgRegExp.FindStringSubmatch(input)
	if len(matches) != 4 {
		return fmt.Errorf("invalid error message format: %q", input)
	}

	errorMessage := matches[1]
	errorCode := matches[2]
	errorKind := matches[3]

	return fmt.Errorf("%s (%s): %s", errorKind, errorCode, errorMessage)
}
