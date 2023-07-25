package reciever

import (
	"bytes"
	"reflect"
	"testing"
)

type MockSerialManagement struct {
	data []byte
	pos  int
}

func (m *MockSerialManagement) Close() error {
	return nil
}

func (m *MockSerialManagement) Read(b []byte) (int, error) {
	if m.pos >= len(m.data) {
		return 0, nil
	}
	n := copy(b, m.data[m.pos:])
	m.pos += n
	return n, nil
}

func (m *MockSerialManagement) Write(b []byte) (int, error) {
	return len(b), nil
}

func TestListen(t *testing.T) {
	packetsData := make([]byte, 0)
	packetsData = append(packetsData, append(append(GetHeader(), []byte{0x01, 0x02, 0x03}...), GetFooter()...)...)
	packetsData = append(packetsData, append(append(GetHeader(), []byte{0x04, 0x05, 0x06}...), GetFooter()...)...)

	sm := &MockSerialManagement{
		data: packetsData,
	}

	packets := make(chan Packet, 2)
	go Listen(packets, sm)

	p := <-packets
	if !bytes.Equal(p.Header, GetHeader()) {
		t.Errorf("expected header to be %v, got %v", GetHeader(), p.Header)
	}
	if !bytes.Equal(p.Data, []byte{0x01, 0x02, 0x03}) {
		t.Errorf("expected data to be %v, got %v", []byte{0x01, 0x02, 0x03}, p.Data)
	}
	if !bytes.Equal(p.Footer, GetFooter()) {
		t.Errorf("expected footer to be %v, got %v", GetFooter(), p.Footer)
	}

	p = <-packets
	if !bytes.Equal(p.Header, GetHeader()) {
		t.Errorf("expected header to be %v, got %v", GetHeader(), p.Header)
	}
	if !bytes.Equal(p.Data, []byte{0x04, 0x05, 0x06}) {
		t.Errorf("expected data to be %v, got %v", []byte{0x04, 0x05, 0x06}, p.Data)
	}

	if !bytes.Equal(p.Footer, GetFooter()) {
		t.Errorf("expected footer to be %v, got %v", GetFooter(), p.Footer)
	}
}

func TestTryParsePacket(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  Packet
		rest  []byte
		found bool
	}{
		{
			name:  "valid packet",
			input: append(append(GetHeader(), []byte{0x01, 0x02, 0x03}...), GetFooter()...),
			want: Packet{
				Header: GetHeader(),
				Data:   []byte{0x01, 0x02, 0x03},
				Footer: GetFooter(),
			},
			rest:  []byte{},
			found: true,
		},
		{
			name:  "packet with trailing data",
			input: append(append(GetHeader(), []byte{0x01, 0x02, 0x03}...), append(GetFooter(), []byte{0x04, 0x05, 0x06}...)...),
			want: Packet{
				Header: GetHeader(),
				Data:   []byte{0x01, 0x02, 0x03},
				Footer: GetFooter(),
			},
			rest:  []byte{0x04, 0x05, 0x06},
			found: true,
		},
		{
			name:  "no header",
			input: append([]byte{0x01, 0x02, 0x03}, GetFooter()...),
			want:  Packet{},
			rest:  append([]byte{0x01, 0x02, 0x03}, GetFooter()...),
			found: false,
		},
		{
			name:  "no footer",
			input: append(GetHeader(), []byte{0x01, 0x02, 0x03}...),
			want:  Packet{},
			rest:  append(GetHeader(), []byte{0x01, 0x02, 0x03}...),
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, rest, found := TryParsePacket(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TryParsePacket() got = %v, want %v", got, tt.want)
			}
			if !bytes.Equal(rest, tt.rest) {
				t.Errorf("TryParsePacket() rest = %v, want %v", rest, tt.rest)
			}
			if found != tt.found {
				t.Errorf("TryParsePacket() found = %v, want %v", found, tt.found)
			}
		})
	}
}
