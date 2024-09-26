package bencode

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
)

func Decode(data string) (interface{}, error) {
	decoder := bencodeDecoder{data: []byte(data), position: 0}
	return decoder.readString()
}

type bencodeDecoder struct {
	data     []byte
	position int
}

func (b *bencodeDecoder) readString() (interface{}, error) {
	char := b.data[b.position]
	switch {
	case b.isNumber(char):
		return b.decodeNumber()
	case b.isString(char):
		return b.decodeString()
	case b.isDict(char):
		return b.decodeDict()
	case b.isList(char):
		return b.decodeList()
	}

	return nil, errors.New("error decoding")
}
func (b *bencodeDecoder) isNumber(char byte) bool {
	return char == 'i'
}

func (b *bencodeDecoder) decodeNumber() (int, error) {
	for i := b.position; i < len(b.data); i++ {
		if b.data[i] == 'e' {
			numStr := string(b.data[b.position+1 : i])
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", numStr)
			}

			b.position = i + 1
			return num, nil
		}
	}

	return 0, errors.New("error decoding number")
}

func (b *bencodeDecoder) isString(char byte) bool {
	return char >= '0' && char <= '9'
}

func (b *bencodeDecoder) decodeString() (string, error) {
	colonIndex := slices.Index(b.data[b.position:], ':')
	if colonIndex == -1 {
		return "", fmt.Errorf("malformed string")
	}
	colonIndex += b.position
	lenghtStr := string(b.data[b.position:colonIndex])
	lenght, err := strconv.Atoi(lenghtStr)
	if err != nil {
		return "", fmt.Errorf("invalid lenght: %s", lenghtStr)
	}
	b.position = colonIndex + 1 // skip :
	endPosition := b.position + lenght

	if endPosition > len(b.data) {
		return "", fmt.Errorf("string lenght exceeds input size")
	}

	str := string(b.data[b.position:endPosition])

	b.position += lenght
	return str, nil
}

func (b *bencodeDecoder) isDict(char byte) bool {
	return char == 'd'
}

func (b *bencodeDecoder) decodeDict() (map[string]interface{}, error) {
	b.position++ // skip the d
	dict := make(map[string]interface{})
	for b.data[b.position] != 'e' {
		key, error := b.decodeString()
		if error != nil {
			return nil, error
		}

		value, error := b.readString()

		if error != nil {
			return nil, error
		}
		dict[key] = value
	}
	b.position++ // skip the e

	return dict, nil
}

func (b *bencodeDecoder) isList(character byte) bool {
	return character == 'l'
}

func (b *bencodeDecoder) decodeList() ([]interface{}, error) {
	b.position++ // skip the l
	var list []interface{}
	for b.data[b.position] != 'e' {
		item, error := b.readString()
		if error != nil {
			return nil, error
		}

		list = append(list, item)
	}
	b.position++
	return list, nil
}
