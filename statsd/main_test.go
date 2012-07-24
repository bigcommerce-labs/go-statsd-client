package statsd

import (
	"testing"
	"bytes"
	"bufio"
	//"fmt"
)

func NewTestClient(prefix string) (*Client, *bytes.Buffer, chan bool) {
	b := &bytes.Buffer{}
	buf := bufio.NewReadWriter(bufio.NewReader(b), bufio.NewWriter(b))
	// make data syncronous for testing
	data := make(chan string)
	quit := make(chan bool)
	f := &Client{buf: buf, prefix: prefix, data: data, quit: quit}
	f.StartSender()
	return f, b, quit
}


func TestGuage(t *testing.T) {
	f, buf, q := NewTestClient("test")

	err := f.Guage("guage", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	q <- true
	b := buf.String()
	buf.Reset()
	expected := "test.guage:1|g"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestIncRatio(t *testing.T) {
	f, buf, q := NewTestClient("test")

	err := f.Inc("count", 1, 0.999999)
	if err != nil {
		t.Fatal(err)
	}

	q <- true
	b := buf.String()
	buf.Reset()
	expected := "test.count:1|c|@0.999999"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestInc(t *testing.T) {
	f, buf, q := NewTestClient("test")

	err := f.Inc("count", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	q <- true
	b := buf.String()
	buf.Reset()
	expected := "test.count:1|c"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestDec(t *testing.T) {
	f, buf, q := NewTestClient("test")

	err := f.Dec("count", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	q <- true
	b := buf.String()
	buf.Reset()
	expected := "test.count:-1|c"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestTiming(t *testing.T) {
	f, buf, q := NewTestClient("test")

	err := f.Timing("timing", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	q <- true
	b := buf.String()
	buf.Reset()
	expected := "test.timing:1|ms"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestEmptyPrefix(t *testing.T) {
	f, buf, q := NewTestClient("")

	err := f.Inc("count", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	q <- true
	b := buf.String()
	buf.Reset()
	expected := "count:1|c"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

