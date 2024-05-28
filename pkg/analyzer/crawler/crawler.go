package crawler

import "C"
import (
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"os"
	"path"
	"strings"
	"time"
)

type Crawler struct {
	Path   string
	Result chan Record
}

type Record struct {
	Timestamp  time.Time
	File       string
	Seq        int
	RelativeID int
	Src        string
	Dst        string
	Data       string
}

func (c Crawler) Process(ctx context.Context) error {
	csvPath := fmt.Sprintf("%s.csv", c.Path)

	dataRecords, err := makeDataRecords(csvPath)
	if err != nil {
		return err
	}

	handle, err := pcap.OpenOffline(c.Path)
	if err != nil {
		return err
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	return c.processPackets(
		ctx,
		packetSource.Packets(),
		dataRecords,
	)
}

func (c Crawler) processPackets(ctx context.Context, packets <-chan gopacket.Packet, dataRecords []string) error {
	i := 0

	for {
		var (
			packet gopacket.Packet
			ok     bool
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case packet, ok = <-packets:
			if !ok {
				return nil
			}
			i++
		}

		if packet == nil {
			continue
		}

		tcp, ok := findLayer[*layers.TCP](packet)
		if !ok || tcp == nil {
			// tcp layer not found
			continue
		}
		ipv4, ok := findLayer[*layers.IPv4](packet)
		if !ok || ipv4 == nil {
			// ipv4 layer not found
			continue
		}

		data := matchData(tcp, dataRecords[i-1])

		if data == "" {
			// this packet not satisfied
			continue
		}

		c.Result <- Record{
			Timestamp:  packet.Metadata().Timestamp,
			File:       path.Base(c.Path),
			Seq:        int(tcp.Seq),
			RelativeID: i,
			Src:        ipv4.SrcIP.String(),
			Dst:        ipv4.DstIP.String(),
			Data:       data,
		}
	}
}

func matchData(tcp *layers.TCP, data string) string {
	if tcp.RST {
		return "RST"
	}

	return dataMatchesLabels(data)
}

func dataMatchesLabels(data string) string {
	labels := []string{
		"TCP Retransmission",
		"TCP Spurious Retransmission",
		"TCP Fast Retransmission",
	}

	for _, label := range labels {
		if strings.Contains(data, label) {
			return label
		}
	}

	return ""
}

func makeDataRecords(p string) ([]string, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil
}

func findLayer[T any](packet gopacket.Packet) (t T, ok bool) {
	for _, layer := range packet.Layers() {
		t, ok = layer.(T)
		if ok {
			return t, ok
		}
	}
	return t, ok
}
