package review

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Severity string

const (
	SeverityBlocking    Severity = "blocking"
	SeverityNonBlocking Severity = "non_blocking"
)

const (
	VerdictPass = "pass"
	VerdictFail = "fail"
)

var ErrInvalidPacket = errors.New("invalid review packet")

type Finding struct {
	ID       string   `json:"id"`
	Severity Severity `json:"severity"`
	Summary  string   `json:"summary"`
}

type Packet struct {
	Verdict  string    `json:"verdict"`
	Findings []Finding `json:"findings,omitempty"`
	Raw      string    `json:"-"`
}

func ParseNDJSON(text string) (Packet, error) {
	var packet Packet
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var frame struct {
			Type     string   `json:"type"`
			Verdict  string   `json:"verdict"`
			ID       string   `json:"id"`
			Severity Severity `json:"severity"`
			Summary  string   `json:"summary"`
		}
		if err := json.Unmarshal([]byte(line), &frame); err != nil {
			return Packet{}, fmt.Errorf("%w: %v", ErrInvalidPacket, err)
		}
		switch frame.Type {
		case "verdict", "done":
			if frame.Verdict != "" {
				if !ValidVerdict(frame.Verdict) {
					return Packet{}, fmt.Errorf("%w: invalid verdict %q", ErrInvalidPacket, frame.Verdict)
				}
				packet.Verdict = frame.Verdict
			}
		case "finding":
			if frame.ID == "" {
				frame.ID = fmt.Sprintf("finding-%d", len(packet.Findings)+1)
			}
			if frame.Severity == "" {
				frame.Severity = SeverityNonBlocking
			}
			if !ValidSeverity(frame.Severity) {
				return Packet{}, fmt.Errorf("%w: invalid severity %q", ErrInvalidPacket, frame.Severity)
			}
			packet.Findings = append(packet.Findings, Finding{ID: frame.ID, Severity: frame.Severity, Summary: frame.Summary})
		case "workspace_mutation":
			packet.Findings = append(packet.Findings, Finding{ID: "workspace_mutation", Severity: SeverityBlocking, Summary: "provider mutated workspace during review"})
		case "tick", "partial":
		default:
			return Packet{}, fmt.Errorf("%w: unknown frame type %q", ErrInvalidPacket, frame.Type)
		}
	}
	if err := scanner.Err(); err != nil {
		return Packet{}, fmt.Errorf("%w: %v", ErrInvalidPacket, err)
	}
	packet.Raw = text
	if packet.Verdict == "" {
		packet.Verdict = VerdictFromFindings(packet.Findings)
	}
	return packet, nil
}

func ValidatePacket(packet Packet) error {
	if !ValidVerdict(packet.Verdict) {
		return fmt.Errorf("%w: invalid verdict %q", ErrInvalidPacket, packet.Verdict)
	}
	for _, finding := range packet.Findings {
		if !ValidSeverity(finding.Severity) {
			return fmt.Errorf("%w: invalid severity %q", ErrInvalidPacket, finding.Severity)
		}
	}
	return nil
}

func VerdictFromFindings(findings []Finding) string {
	for _, finding := range findings {
		if finding.Severity == SeverityBlocking {
			return VerdictFail
		}
	}
	return VerdictPass
}

func ValidVerdict(verdict string) bool {
	switch verdict {
	case VerdictPass, VerdictFail:
		return true
	default:
		return false
	}
}

func ValidSeverity(severity Severity) bool {
	switch severity {
	case SeverityBlocking, SeverityNonBlocking:
		return true
	default:
		return false
	}
}
