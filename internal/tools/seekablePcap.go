package tools

import (
	"math"
	"os"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/pcapgo"
)

type (
	SeekablePcapHolder struct {
		filename    string
		handle      *os.File
		reader      *pcapgo.NgReader
		source      *gopacket.PacketSource
		packetIndex uint64
	}
)

func NewSeekablePcapHolder(filename string) *SeekablePcapHolder {
	return &SeekablePcapHolder{
		filename:    filename,
		packetIndex: math.MaxUint32,
	}
}

func (s *SeekablePcapHolder) Close() {
	if s.handle != nil {
		s.handle.Close()
		s.handle = nil
	}
}

func (s *SeekablePcapHolder) Packet(packetIndex uint64) (gopacket.Packet, error) {
	if s.packetIndex > packetIndex {
		s.Close()
		handle, err := os.Open(s.filename)
		if err != nil {
			return nil, err
		}
		reader, err := pcapgo.NewNgReader(handle, pcapgo.DefaultNgReaderOptions)
		if err != nil {
			return nil, err
		}
		s.handle = handle
		s.reader = reader
		s.source = gopacket.NewPacketSource(reader, reader.LinkType())
		s.packetIndex = 0
	}
	for s.packetIndex < packetIndex {
		_, err := s.source.NextPacket()
		if err != nil {
			return nil, err
		}
		s.packetIndex++
	}
	pkt, err := s.source.NextPacket()
	if err != nil {
		return nil, err
	}
	s.packetIndex++
	return pkt, nil
}
