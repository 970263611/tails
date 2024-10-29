package charsettool

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Charset uint32

const (
	UTF8 Charset = iota
	GBK
	GB18030
	HZGB2312
)

/*
*
字符集转换
*/
func Convert(data []byte, from, to Charset) ([]byte, error) {
	if from == to {
		return data, nil
	}
	var resultBytes []byte
	var decoder *encoding.Decoder
	switch from {
	case GBK:
		decoder = simplifiedchinese.GBK.NewDecoder()
	case GB18030:
		decoder = simplifiedchinese.GB18030.NewDecoder()
	case HZGB2312:
		decoder = simplifiedchinese.HZGB2312.NewDecoder()
	default:
	}
	if decoder != nil {
		bytes, _, err := transform.Bytes(decoder, data)
		if err != nil {
			return data, err
		}
		resultBytes = bytes
	}
	var encoder *encoding.Encoder
	switch to {
	case GBK:
		encoder = simplifiedchinese.GBK.NewEncoder()
	case GB18030:
		encoder = simplifiedchinese.GB18030.NewEncoder()
	case HZGB2312:
		encoder = simplifiedchinese.HZGB2312.NewEncoder()
	default:
	}
	if encoder != nil {
		bytes, _, err := transform.Bytes(encoder, resultBytes)
		if err != nil {
			return data, err
		}
		resultBytes = bytes
	}
	return resultBytes, nil
}
