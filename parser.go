package gphoto2cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	autoDetechRegExp = regexp.MustCompile("([a-zA-Z0-9_ ]+)(usb:(\\d+),(\\d+))")
	errorMsgRegExp   = regexp.MustCompile(`\*\*\* Error \*\*\* *\n(.+)\n+\*\*\* Error( \(([-0-9 ]+): '([^']+)'\))* \*\*\*`)
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
	v := &ConfigValue{
		Choices: make(map[int]string),
	}

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
				choiceParts := strings.SplitN(value, " ", 2)
				if len(choiceParts) == 2 {
					index := choiceParts[0]
					choice := choiceParts[1]
					i, _ := strconv.Atoi(index)
					v.Choices[i] = choice
				}
			}
		}
	}

	return v
}

func (p *parser) ParseErrorMessage(input string) error {
	matches := errorMsgRegExp.FindStringSubmatch(input)
	if len(matches) != 5 {
		return nil
	}

	errorMessage := matches[1]
	errorCode := matches[3]
	errorKind := matches[4]

	if errorCode == "" && errorKind == "" {
		return fmt.Errorf(errorMessage)
	}

	return fmt.Errorf("%s (%s): %s", errorKind, errorCode, errorMessage)
}

func (p *parser) ParseCapture(input string) *CaptureResult {
	lines := strings.Split(input, "\n")
	r := &CaptureResult{}
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "New file is in location"):
			r.OnCamera = strings.TrimSpace(strings.TrimPrefix(line, "New file is in location"))
		case strings.HasPrefix(line, "Saving file as"):
			r.LocalFile = strings.TrimSpace(strings.TrimPrefix(line, "Saving file as"))
		case strings.HasPrefix(line, "Deleting file"):
			r.Deleted = true
		}
	}

	return r
}
