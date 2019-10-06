package main

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"testing"
)

const TEST_XML string = `<?xml version="1.0" encoding="UTF-8"?>
<books>
  <book seq="1">
    <author>William H. Gaddis</author>
    <review>One of the great seminal American novels of the 20th century.</review>
    <title>The Recognitions</title>
  </book>
  <book seq="2">
    <author>Austin Tappan Wright</author>
    <review>An example of earlier 20th century American utopian fiction.</review>
    <title>Islandia</title>
  </book>
  <book seq="3">
    <author>John Hawkes</author>
    <review>A lyrical novel about the construction of Ft. Peck Dam in Montana.</review>
    <title>The Beetle Leg</title>
  </book>
  <book seq="4">
    <author>T.E. Porter</author>
    <review>A magical novella.</review>
    <title>King's Day</title>
  </book>
</books>`

const TEST_JSON string = `{"books":{"book":[{"-seq":"1","author":"William H. Gaddis","review":"One of the great seminal American novels of the 20th century.","title":"The Recognitions"},{"-seq":"2","author":"Austin Tappan Wright","review":"An example of earlier 20th century American utopian fiction.","title":"Islandia"},{"-seq":"3","author":"John Hawkes","review":"A lyrical novel about the construction of Ft. Peck Dam in Montana.","title":"The Beetle Leg"},{"-seq":"4","author":"T.E. Porter","review":"A magical novella.","title":"King's Day"}]}}`

func TestNewConv(t *testing.T) {

	conv := NewConv()
	conv_dummy := &Converter{}

	if !reflect.DeepEqual(conv, conv_dummy) {
		t.Errorf("NewConv() wrong return %d != %d", conv, conv_dummy)
	}
}

func TestCopyInput(t *testing.T) {

	var check_xml bytes.Buffer
	var dummy bytes.Buffer
	check_xml.WriteString(TEST_XML)
	conv_dummy := &Converter{check_xml, XML, dummy, Unsupported}

	req := httptest.NewRequest("POST", "/", &check_xml)

	conv := NewConv()
	conv.CopyInput(req)

	if !reflect.DeepEqual(conv, conv_dummy) {
		t.Errorf("SwapFormat() bad state %d != %d", conv, conv_dummy)
	}
}

func TestSwapFormat(t *testing.T) {

	conv := NewConv()
	conv.Payload.WriteString(TEST_XML)
	conv.ContentType = XML

	conv.SwapFormat()

	if conv.ResultContentType != JSON {
		t.Error("Content-type key incorrect")
	}

	if conv.Result.String() != TEST_JSON {
		t.Error("Recieved payload data wrong")
	}

}

func TestMakeReply(t *testing.T) {

	conv := NewConv()

	conv.Result.WriteString(`{"Some": "JSON"}`)
	conv.ContentType = XML // need to pass MakeReply checks
	conv.ResultContentType = JSON

	w := httptest.NewRecorder()
	conv.MakeReply(w)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Error("Bad http status code")
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-type key incorrect")
	}

	if string(body) != `{"Some": "JSON"}` {
		t.Error("Wrong body")
	}

}

func TestXml2Json(t *testing.T) {

	conv := NewConv()
	conv.Payload.WriteString(TEST_XML)
	conv.ContentType = XML

	conv.xml2json()

	if conv.ResultContentType != JSON {
		t.Error("Content-type key incorrect")
	}

	if conv.Result.String() != TEST_JSON {
		t.Error("Recieved payload data wrong")
	}

}

func TestJson2Xml(t *testing.T) {

	conv := NewConv()
	conv.Payload.WriteString(TEST_JSON)
	conv.ContentType = JSON

	conv.json2xml()

	if conv.ResultContentType != XML {
		t.Error("Content-type key incorrect")
	}

	if conv.Result.String() != TEST_XML {
		t.Error("Recieved payload data wrong")
	}

}

func BenchmarkJson2Xml(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conv := NewConv()
		conv.Payload.WriteString(TEST_JSON)
		conv.ContentType = JSON
		conv.json2xml()
	}
}

func BenchmarkXml2Json(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conv := NewConv()
		conv.Payload.WriteString(TEST_XML)
		conv.ContentType = XML
		conv.xml2json()
	}
}
