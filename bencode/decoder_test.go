package bencode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	bencode := "d4:dictd3:1234:test3:4565:thinge4:listl11:list-item-111:list-item-2e6:numberi123456e6:string5:valuee"
	decodedValue, err := Decode(bencode)
	assert.NoError(t, err)
	expected := map[string]interface{}{
		"dict": map[string]interface{}{
			"123": "test",
			"456": "thing",
		},
		"list":   []interface{}{"list-item-1", "list-item-2"},
		"number": 123456,
		"string": "value",
	}
	assert.Equal(t, expected, decodedValue)
}

func TestDecodeNumber(t *testing.T) {
	decoder := &bencodeDecoder{data: []byte("i123456e"), position: 0}
	number, err := decoder.decodeNumber()
	assert.NoError(t, err)
	assert.Equal(t, 123456, number)
}

func TestDecodeNumberNoEndingCharacter(t *testing.T) {
	decoder := &bencodeDecoder{data: []byte("i123456"), position: 0}
	number, err := decoder.decodeNumber()
	assert.Error(t, err)
	assert.Equal(t, 0, number)
}

func TestDecodeNumberNoNumeralValue(t *testing.T) {
	decoder := &bencodeDecoder{data: []byte("iabce"), position: 0}
	_, err := decoder.decodeNumber()
	assert.Error(t, err)
}

func TestIsNumber(t *testing.T) {
	decoder := &bencodeDecoder{}
	assert.True(t, decoder.isNumber('i'))
	assert.False(t, decoder.isNumber('a'))
}

func TestIsString(t *testing.T) {
	decoder := &bencodeDecoder{}
	assert.True(t, decoder.isString('0'))
	assert.True(t, decoder.isString('9'))
	assert.False(t, decoder.isString('a'))
}

func TestDecodeString(t *testing.T) {
	decoder := &bencodeDecoder{data: []byte("5:value"), position: 0}
	str, err := decoder.decodeString()
	assert.NoError(t, err)
	assert.Equal(t, "value", str)
}

func TestIsDict(t *testing.T) {
	decoder := &bencodeDecoder{}
	assert.True(t, decoder.isDict('d'))
	assert.False(t, decoder.isDict('l'))
}

func TestDecodeDict(t *testing.T) {
	decoder := &bencodeDecoder{data: []byte("d4:dictd3:1234:test3:4565:thingee"), position: 0}
	dict, err := decoder.decodeDict()
	assert.NoError(t, err)
	expected := map[string]interface{}{
		"dict": map[string]interface{}{
			"123": "test",
			"456": "thing",
		},
	}
	assert.Equal(t, expected, dict)
}

func TestIsList(t *testing.T) {
	decoder := &bencodeDecoder{}
	assert.True(t, decoder.isList('l'))
	assert.False(t, decoder.isList('d'))
}

func TestDecodeList(t *testing.T) {
	decoder := &bencodeDecoder{data: []byte("l11:list-item-1i123ee"), position: 0}
	list, err := decoder.decodeList()
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"list-item-1", 123}, list)
}

