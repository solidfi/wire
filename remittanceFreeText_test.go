package wire

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// RemittanceFreeText creates a RemittanceFreeText
func mockRemittanceFreeText() *RemittanceFreeText {
	rft := NewRemittanceFreeText()
	rft.LineOne = "Remittance Free Text Line One"
	rft.LineTwo = "Remittance Free Text Line Two"
	rft.LineThree = "Remittance Free Text Line Three"
	return rft
}

// TestMockRemittanceFreeText validates mockRemittanceFreeText
func TestMockRemittanceFreeText(t *testing.T) {
	rft := mockRemittanceFreeText()

	require.NoError(t, rft.Validate(), "mockRemittanceFreeText does not validate and will break other tests")
}

// TestRemittanceFreeTextLineOneAlphaNumeric validates RemittanceFreeText LineOne is alphanumeric
func TestRemittanceFreeTextLineOneAlphaNumeric(t *testing.T) {
	rft := mockRemittanceFreeText()
	rft.LineOne = "®"

	err := rft.Validate()

	require.EqualError(t, err, fieldError("LineOne", ErrNonAlphanumeric, rft.LineOne).Error())
}

// TestRemittanceFreeTextLineTwoAlphaNumeric validates RemittanceFreeText LineTwo is alphanumeric
func TestRemittanceFreeTextLineTwoAlphaNumeric(t *testing.T) {
	rft := mockRemittanceFreeText()
	rft.LineTwo = "®"

	err := rft.Validate()

	require.EqualError(t, err, fieldError("LineTwo", ErrNonAlphanumeric, rft.LineTwo).Error())
}

// TestRemittanceFreeTextLineThreeAlphaNumeric validates RemittanceFreeText LineThree is alphanumeric
func TestRemittanceFreeTextLineThreeAlphaNumeric(t *testing.T) {
	rft := mockRemittanceFreeText()
	rft.LineThree = "®"

	err := rft.Validate()

	require.EqualError(t, err, fieldError("LineThree", ErrNonAlphanumeric, rft.LineThree).Error())
}

// TestParseRemittanceFreeTextWrongLength parses a wrong RemittanceFreeText record length
func TestParseRemittanceFreeTextWrongLength(t *testing.T) {
	var line = "{8750}Remittance Free Text Line One                                                                                                              Remittance Free Text Line Two                                                                                                               Remittance Free Text Line Three                                                                                                             "
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseRemittanceFreeText()

	require.EqualError(t, err, r.parseError(fieldError("LineThree", ErrValidLength)).Error())
}

// TestParseRemittanceFreeTextReaderParseError parses a wrong RemittanceFreeText reader parse error
func TestParseRemittanceFreeTextReaderParseError(t *testing.T) {
	var line = "{8750}Re®ittance Free Text Line One                                                                                                               Remittance Free Text Line Two                                                                                                               Remittance Free Text Line Three                                                                                                            "
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseRemittanceFreeText()

	require.EqualError(t, err, r.parseError(fieldError("LineOne", ErrNonAlphanumeric, "Re®ittance Free Text Line One")).Error())

	_, err = r.Read()

	require.EqualError(t, err, r.parseError(fieldError("LineOne", ErrNonAlphanumeric, "Re®ittance Free Text Line One")).Error())
}

// TestRemittanceFreeTextTagError validates a RemittanceFreeText tag
func TestRemittanceFreeTextTagError(t *testing.T) {
	rft := mockRemittanceFreeText()
	rft.tag = "{9999}"

	require.EqualError(t, rft.Validate(), fieldError("tag", ErrValidTagForType, rft.tag).Error())
}

// TestStringRemittanceFreeTextVariableLength parses using variable length
func TestStringRemittanceFreeTextVariableLength(t *testing.T) {
	var line = "{8750}"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseRemittanceFreeText()
	require.Nil(t, err)

	line = "{8750}                                                                                                                                                                                                                                                                                                                                                                                                                                    NNN"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseRemittanceFreeText()
	require.EqualError(t, err, r.parseError(NewTagMaxLengthErr()).Error())

	line = "{8750}****************************"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseRemittanceFreeText()
	require.EqualError(t, err, r.parseError(NewTagMaxLengthErr()).Error())

	line = "{8750}*"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseRemittanceFreeText()
	require.Equal(t, err, nil)
}

// TestStringRemittanceFreeTextOptions validates string() with options
func TestStringRemittanceFreeTextOptions(t *testing.T) {
	var line = "{8750}"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseRemittanceFreeText()
	require.Equal(t, err, nil)

	str := r.currentFEDWireMessage.RemittanceFreeText.String()
	require.Equal(t, str, "{8750}                                                                                                                                                                                                                                                                                                                                                                                                                                    ")

	str = r.currentFEDWireMessage.RemittanceFreeText.String(true)
	require.Equal(t, str, "{8750}*")
}
