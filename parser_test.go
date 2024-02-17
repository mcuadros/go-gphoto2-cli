package gphoto2cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var autoDetectFixture = `Model                          Port
----------------------------------------------------------
Canon EOS R10                  usb:003,017
Canon EOS R10                  usb:003,016`

func TestParseAutoDetect(t *testing.T) {
	p := &parser{}
	ports := p.ParseAutoDetect(autoDetectFixture)
	require.Len(t, ports, 2)

	for _, p := range ports {
		assert.Equal(t, "Canon EOS R10", p.Model)
		assert.True(t, strings.Contains(p.Port, "usb:"))
		assert.Len(t, p.Port, 11)
	}
}

var getConfigTextFixture = `Label: Serial Number
Readonly: 1
Type: TEXT
Current: 2c90dfcfb648160375c8af85dc53d343
END`

func TestParseGetConfigText(t *testing.T) {
	p := &parser{}
	value := p.ParseGetConfig(getConfigTextFixture)

	assert.True(t, value.Readonly)
	assert.Equal(t, "Serial Number", value.Label)
	assert.Equal(t, "TEXT", value.Type)
	assert.Equal(t, "2c90dfcfb648160375c8af85dc53d343", value.Current)
}

var getConfigRadioFixture = `Label: Picture Style
Readonly: 0
Type: RADIO
Current: Auto
Choice: 0 Auto
END`

func TestParseGetConfigRadio(t *testing.T) {
	p := &parser{}
	value := p.ParseGetConfig(getConfigRadioFixture)

	assert.False(t, value.Readonly)
	assert.Equal(t, "Picture Style", value.Label)
	assert.Equal(t, "RADIO", value.Type)
	assert.Equal(t, "Auto", value.Current)
	assert.Equal(t, []string{"0", "Auto"}, value.Choices)
}

var cameraSummaryFixture = `Camera summary:
Manufacturer: Canon.Inc
Model: Canon EOS R10
  Version: 3-1.3.0
  Serial Number: 2c90dfcfb648160375c8af85dc53d343
Vendor Extension ID: 0xb (1.0)
Vendor Extension Description:

Capture Formats: JPEG Unknown(b108) Unknown(b10b) Unknown(b982)
Display Formats: Association/Directory, Script, DPOF, MS AVI, MS Wave, JPEG, Unknown(b103), Unknown(bf02), Defined Type, Unknown(b105), Unknown(b982), Unknown(b10a), Unknown(b10b), Unknown(b109)

Device Capabilities:
        File Download, File Deletion, File Upload
        No Image Capture, No Open Capture, Canon EOS Capture, Canon EOS Capture 2
        Canon Wifi support

Storage Devices Summary:
store_00020001:
        StorageDescription: SD1
        VolumeLabel:
        Storage Type: Removable RAM (memory card)
        Filesystemtype: Digital Camera Layout (DCIM)
        Access Capability: Read-Write
        Maximum Capability: 62193139712 (59312 MB)
        Free Space (Bytes): 62159454208 (59279 MB)
        Free Space (Images): -1

Device Property Summary:
Property 0xd402:(read only) (type=0xffff) 'Canon EOS R10'
Property 0xd407:(read only) (type=0x6) 1
Property 0xd406:(readwrite) (type=0xffff) 'Unknown Initiator'
Property 0xd303:(read only) (type=0x2) 1
Battery Level(0x5001):(read only) (type=0x2) Enumeration [100,0,75,0,50] value: 100% (100)`

func TestParseCameraSummary(t *testing.T) {
	p := &parser{}
	s := p.ParseCameraSummary(cameraSummaryFixture)

	assert.Equal(t, s.Model, "Canon EOS R10")
	assert.Equal(t, s.SerialNumber, "2c90dfcfb648160375c8af85dc53d343")
	assert.Equal(t, s.Manufacturer, "Canon.Inc")
	assert.Equal(t, s.Version, "3-1.3.0")
	assert.Len(t, s.CaptureFormats, 1)
	assert.Len(t, s.DisplayFormats, 14)
}

var fixtureErrorMessage = `
*** Error ***    
serialnumberd not found in configuration tree.
*** Error (-1: 'Unspecified error') ***

`

func TestParseErrorMessage(t *testing.T) {
	p := &parser{}
	err := p.ParseErrorMessage(fixtureErrorMessage)
	assert.Equal(t, err.Error(), "Unspecified error (-1): serialnumberd not found in configuration tree.")
}

var fixtureErrorCaptureMessage = `

*** Error ***
Canon EOS Capture failed to release: Perhaps no focus?

*** Error ***
Canon EOS Capture failed to release: Perhaps no focus?
ERROR: Could not capture image.
ERROR: Could not capture.
`

func TestParseErrorMessageCapture(t *testing.T) {
	p := &parser{}
	err := p.ParseErrorMessage(fixtureErrorCaptureMessage)
	assert.Equal(t, err.Error(), "Canon EOS Capture failed to release: Perhaps no focus?")
}
