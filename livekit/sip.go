package livekit

import (
	"errors"
	"fmt"
	"strings"
)

// ToProto implements DataPacket in Go SDK.
func (p *SipDTMF) ToProto() *DataPacket {
	return &DataPacket{
		Value: &DataPacket_SipDtmf{
			SipDtmf: p,
		},
	}
}

// AsInbound converts legacy SIPTrunkInfo to SIPInboundTrunkInfo.
func (p *SIPTrunkInfo) AsInbound() *SIPInboundTrunkInfo {
	if p == nil || p.Kind == SIPTrunkInfo_TRUNK_OUTBOUND {
		return nil
	}
	var nums []string
	if p.OutboundNumber != "" {
		nums = []string{p.OutboundNumber}
	}
	return &SIPInboundTrunkInfo{
		SipTrunkId:       p.SipTrunkId,
		Name:             p.Name,
		Metadata:         p.Metadata,
		Numbers:          nums,
		AllowedAddresses: p.InboundAddresses,
		AllowedNumbers:   p.InboundNumbers,
		AuthUsername:     p.InboundUsername,
		AuthPassword:     p.InboundPassword,
	}
}

// AsTrunkInfo converts SIPInboundTrunkInfo to legacy SIPTrunkInfo.
func (p *SIPInboundTrunkInfo) AsTrunkInfo() *SIPTrunkInfo {
	if p == nil {
		return nil
	}
	var num string
	if len(p.Numbers) != 0 {
		num = p.Numbers[0]
	}
	return &SIPTrunkInfo{
		SipTrunkId:       p.SipTrunkId,
		Kind:             SIPTrunkInfo_TRUNK_INBOUND,
		Name:             p.Name,
		Metadata:         p.Metadata,
		OutboundNumber:   num,
		InboundAddresses: p.AllowedAddresses,
		InboundNumbers:   p.AllowedNumbers,
		InboundUsername:  p.AuthUsername,
		InboundPassword:  p.AuthPassword,
	}
}

// AsOutbound converts legacy SIPTrunkInfo to SIPOutboundTrunkInfo.
func (p *SIPTrunkInfo) AsOutbound() *SIPOutboundTrunkInfo {
	if p == nil || p.Kind == SIPTrunkInfo_TRUNK_INBOUND {
		return nil
	}
	var nums []string
	if p.OutboundNumber != "" {
		nums = []string{p.OutboundNumber}
	}
	return &SIPOutboundTrunkInfo{
		SipTrunkId:   p.SipTrunkId,
		Name:         p.Name,
		Metadata:     p.Metadata,
		Address:      p.OutboundAddress,
		Transport:    p.Transport,
		Numbers:      nums,
		AuthUsername: p.OutboundUsername,
		AuthPassword: p.OutboundPassword,
	}
}

// AsTrunkInfo converts SIPOutboundTrunkInfo to legacy SIPTrunkInfo.
func (p *SIPOutboundTrunkInfo) AsTrunkInfo() *SIPTrunkInfo {
	if p == nil {
		return nil
	}
	var num string
	if len(p.Numbers) != 0 {
		num = p.Numbers[0]
	}
	return &SIPTrunkInfo{
		SipTrunkId:       p.SipTrunkId,
		Kind:             SIPTrunkInfo_TRUNK_OUTBOUND,
		Name:             p.Name,
		Metadata:         p.Metadata,
		OutboundAddress:  p.Address,
		Transport:        p.Transport,
		OutboundNumber:   num,
		OutboundUsername: p.AuthUsername,
		OutboundPassword: p.AuthPassword,
	}
}

func validateHeader(header string) error {
	v := strings.ToLower(header)
	if !strings.HasPrefix(v, "x-") {
		return fmt.Errorf("only X-* headers are allowed: %s", header)
	}
	return nil
}

func validateHeaderKeys(headers map[string]string) error {
	for k := range headers {
		if err := validateHeader(k); err != nil {
			return err
		}
	}
	return nil
}

func validateHeaderValues(headers map[string]string) error {
	for _, v := range headers {
		if err := validateHeader(v); err != nil {
			return err
		}
	}
	return nil
}

func (p *SIPTrunkInfo) Validate() error {
	if len(p.InboundNumbersRegex) != 0 {
		return fmt.Errorf("trunks with InboundNumbersRegex are deprecated")
	}
	return nil
}

func (p *CreateSIPOutboundTrunkRequest) Validate() error {
	if p.Trunk == nil {
		return errors.New("missing trunk")
	}
	if p.Trunk.SipTrunkId != "" {
		return errors.New("trunk id must not be set")
	}
	if err := p.Trunk.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *CreateSIPInboundTrunkRequest) Validate() error {
	if p.Trunk == nil {
		return errors.New("missing trunk")
	}
	if p.Trunk.SipTrunkId != "" {
		return errors.New("trunk id must not be set")
	}
	if err := p.Trunk.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *SIPInboundTrunkInfo) Validate() error {
	hasAuth := p.AuthUsername != "" || p.AuthPassword != ""
	hasCIDR := len(p.AllowedAddresses) != 0
	hasNumbers := len(p.Numbers) != 0 // TODO: remove this condition, it doesn't really help with security
	if !hasAuth && !hasCIDR && !hasNumbers {
		return errors.New("for security, one of the fields must be set: AuthUsername+AuthPassword, AllowedAddresses or Numbers")
	}
	if err := validateHeaderKeys(p.Headers); err != nil {
		return err
	}
	if err := validateHeaderKeys(p.HeadersToAttributes); err != nil {
		return err
	}
	if err := validateHeaderValues(p.AttributesToHeaders); err != nil {
		return err
	}
	return nil
}

func (p *SIPOutboundTrunkInfo) Validate() error {
	if len(p.Numbers) == 0 {
		return errors.New("no trunk numbers specified")
	}
	if p.Address == "" {
		return errors.New("no outbound address specified")
	} else if strings.Contains(p.Address, "transport=") {
		return errors.New("trunk transport should be set as a field, not a URI parameter")
	} else if strings.ContainsAny(p.Address, "@;") || strings.HasPrefix(p.Address, "sip:") || strings.HasPrefix(p.Address, "sips:") {
		return errors.New("trunk address should be a hostname or IP, not SIP URI")
	}
	if err := validateHeaderKeys(p.Headers); err != nil {
		return err
	}
	if err := validateHeaderKeys(p.HeadersToAttributes); err != nil {
		return err
	}
	if err := validateHeaderValues(p.AttributesToHeaders); err != nil {
		return err
	}
	return nil
}

func (p *CreateSIPDispatchRuleRequest) Validate() error {
	if p.Rule == nil {
		return errors.New("missing rule")
	}
	return nil
}

func (p *CreateSIPParticipantRequest) Validate() error {
	if p.SipTrunkId == "" {
		return errors.New("missing sip trunk id")
	}
	if p.SipCallTo == "" {
		return errors.New("missing sip callee number")
	} else if strings.Contains(p.SipCallTo, "@") {
		return errors.New("SipCallTo should be a phone number or SIP user, not a full SIP URI")
	}
	if p.RoomName == "" {
		return errors.New("missing room name")
	}
	return nil
}

func (p *TransferSIPParticipantRequest) Validate() error {
	if p.RoomName == "" {
		return errors.New("missing room name")
	}
	if p.ParticipantIdentity == "" {
		return errors.New("missing participant identity")
	}
	if p.TransferTo == "" {
		return errors.New("missing transfer to")
	}

	return nil
}
