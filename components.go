package ics

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type Component interface {
	UnknownPropertiesIANAProperties() []IANAProperty
	SubComponents() []Component
	serialize(b io.Writer)
}

type ComponentBase struct {
	Properties []IANAProperty
	Components []Component
}

func (cb *ComponentBase) UnknownPropertiesIANAProperties() []IANAProperty {
	return cb.Properties
}

func (cb *ComponentBase) SubComponents() []Component {
	return cb.Components
}
func (base ComponentBase) serializeThis(writer io.Writer, componentType string) {
	fmt.Fprint(writer, "BEGIN:"+componentType, "\r\n")
	for _, p := range base.Properties {
		p.serialize(writer)
	}
	for _, c := range base.Components {
		c.serialize(writer)
	}
	fmt.Fprint(writer, "END:"+componentType, "\r\n")
}

type VEvent struct {
	ComponentBase
}

func (c *VEvent) serialize(w io.Writer) {
	c.ComponentBase.serializeThis(w, "VEVENT")
}

func (c *VEvent) Serialize() string {
	b := &bytes.Buffer{}
	c.ComponentBase.serializeThis(b, "VEVENT")
	return b.String()
}

const (
	icalTimeFormat       = "20060102T150405Z"
	icalAllDayTimeFormat = "20060102"
)

func (event *VEvent) SetCreatedTime(t time.Time, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyCreated, t.UTC().Format(icalTimeFormat), props...)
}

func (event *VEvent) SetDtStampTime(t time.Time, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyDtstamp, t.UTC().Format(icalTimeFormat), props...)
}

func (event *VEvent) SetModifiedAt(t time.Time, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyLastModified, t.UTC().Format(icalTimeFormat), props...)
}

func (event *VEvent) SetStartAt(t time.Time, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyDtStart, t.UTC().Format(icalTimeFormat), props...)
}

func (event *VEvent) SetAllDayStartAt(t time.Time, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyDtStart, t.UTC().Format(icalAllDayTimeFormat), props...)
}

func (event *VEvent) SetEndAt(t time.Time, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyDtEnd, t.UTC().Format(icalTimeFormat), props...)
}

func (event *VEvent) SetAllDayEndAt(t time.Time, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyDtEnd, t.UTC().Format(icalAllDayTimeFormat), props...)
}

func (event *VEvent) SetProperty(property ComponentProperty, value string, props ...PropertyParameter) {
	for i := range event.Properties {
		if event.Properties[i].IANAToken == string(property) {
			event.Properties[i].Value = value
			event.Properties[i].ICalParameters = map[string][]string{}
			for _, p := range props {
				k, v := p.KeyValue()
				event.Properties[i].ICalParameters[k] = v
			}
			return
		}
	}
	event.AddProperty(property, value, props...)
}

func (event *VEvent) AddProperty(property ComponentProperty, value string, props ...PropertyParameter) {
	r := IANAProperty{
		BaseProperty{
			IANAToken:      string(property),
			Value:          value,
			ICalParameters: map[string][]string{},
		},
	}
	for _, p := range props {
		k, v := p.KeyValue()
		r.ICalParameters[k] = v
	}
	event.Properties = append(event.Properties, r)
}

func (event *VEvent) SetSummary(s string, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertySummary, s, props...)
}

func (event *VEvent) SetStatus(s ObjectStatus, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyStatus, string(s), props...)
}

func (event *VEvent) SetDescription(s string, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyDescription, s, props...)
}

func (event *VEvent) SetLocation(s string, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyLocation, s, props...)
}

func (event *VEvent) SetURL(s string, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyUrl, s, props...)
}

func (event *VEvent) SetOrganizer(s string, props ...PropertyParameter) {
	event.SetProperty(ComponentPropertyOrganizer, s, props...)
}

func (event *VEvent) AddAttendee(s string, props ...PropertyParameter) {
	event.AddProperty(ComponentPropertyAttendee, "mailto:"+s, props...)
}

type Attendee struct {
	IANAProperty
}

func (attendee *Attendee) Email() string {
	if strings.HasPrefix(attendee.Value, "mailto:") {
		return attendee.Value[len("mailto:"):]
	}
	return attendee.Value
}

func (attendee *Attendee) ParticipationStatus() ParticipationStatus {
	return ParticipationStatus(attendee.getPropertyFirst(ParameterParticipationStatus))
}

func (attendee *Attendee) getPropertyFirst(parameter Parameter) string {
	vs := attendee.getProperty(parameter)
	if vs != nil && len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (attendee *Attendee) getProperty(parameter Parameter) []string {
	if vs, ok := attendee.ICalParameters[string(parameter)]; ok {
		return vs
	}
	return nil
}

func (event *VEvent) Attendees() (r []*Attendee) {
	r = []*Attendee{}
	for i := range event.Properties {
		switch event.Properties[i].IANAToken {
		case string(ComponentPropertyAttendee):
			a := &Attendee{
				event.Properties[i],
			}
			r = append(r, a)
		}
	}
	return
}

func (event *VEvent) Id() string {
	p := event.GetProperty(ComponentPropertyUniqueId)
	if p != nil {
		return p.Value
	}
	return ""
}

func (event *VEvent) GetProperty(componentProperty ComponentProperty) *IANAProperty {
	for i := range event.Properties {
		if event.Properties[i].IANAToken == string(componentProperty) {
			return &event.Properties[i]
		}
	}
	return nil
}

type VTodo struct {
	ComponentBase
}

func (c *VTodo) serialize(w io.Writer) {
	c.ComponentBase.serializeThis(w, "VTODO")
}

func (c *VTodo) Serialize() string {
	b := &bytes.Buffer{}
	c.ComponentBase.serializeThis(b, "VTODO")
	return b.String()
}

type VJournal struct {
	ComponentBase
}

func (c *VJournal) serialize(w io.Writer) {
	c.ComponentBase.serializeThis(w, "VJOURNAL")
}

func (c *VJournal) Serialize() string {
	b := &bytes.Buffer{}
	c.ComponentBase.serializeThis(b, "VJOURNAL")
	return b.String()
}

type VBusy struct {
	ComponentBase
}

func (c *VBusy) Serialize() string {
	b := &bytes.Buffer{}
	c.ComponentBase.serializeThis(b, "VBUSY")
	return b.String()
}

func (c *VBusy) serialize(w io.Writer) {
	c.ComponentBase.serializeThis(w, "VBUSY")
}

type VTimezone struct {
	ComponentBase
}

func (c *VTimezone) Serialize() string {
	b := &bytes.Buffer{}
	c.ComponentBase.serializeThis(b, "VTIMEZONE")
	return b.String()
}

func (c *VTimezone) serialize(w io.Writer) {
	c.ComponentBase.serializeThis(w, "VTIMEZONE")
}

type VAlarm struct {
	ComponentBase
}

func (c *VAlarm) Serialize() string {
	b := &bytes.Buffer{}
	c.ComponentBase.serializeThis(b, "VALARM")
	return b.String()
}

func (c *VAlarm) serialize(w io.Writer) {
	c.ComponentBase.serializeThis(w, "VALARM")
}

type Standard struct {
	ComponentBase
}

func (c *Standard) Serialize() string {
	b := &bytes.Buffer{}
	c.ComponentBase.serializeThis(b, "STANDARD")
	return b.String()
}

func (c *Standard) serialize(w io.Writer) {
	c.ComponentBase.serializeThis(w, "STANDARD")
}

type Daylight struct {
	ComponentBase
}

func (c *Daylight) Serialize() string {
	b := &bytes.Buffer{}
	c.ComponentBase.serializeThis(b, "DAYLIGHT")
	return b.String()
}

func (c *Daylight) serialize(w io.Writer) {
	c.ComponentBase.serializeThis(w, "DAYLIGHT")
}

type GeneralComponent struct {
	ComponentBase
	Token string
}

func (c *GeneralComponent) Serialize() string {
	b := &bytes.Buffer{}
	c.ComponentBase.serializeThis(b, c.Token)
	return b.String()
}

func (c *GeneralComponent) serialize(w io.Writer) {
	c.ComponentBase.serializeThis(w, c.Token)
}

func GeneralParseComponent(cs *CalendarStream, startLine *BaseProperty) (Component, error) {
	var co Component
	switch startLine.Value {
	case "VCALENDAR":
		return nil, errors.New("Malformed calendar")
	case "VEVENT":
		co = ParseVEvent(cs, startLine)
	case "VTODO":
		co = ParseVTodo(cs, startLine)
	case "VJOURNAL":
		co = ParseVJournal(cs, startLine)
	case "VFREEBUSY":
		co = ParseVBusy(cs, startLine)
	case "VTIMEZONE":
		co = ParseVTimezone(cs, startLine)
	case "VALARM":
		co = ParseVAlarm(cs, startLine)
	case "STANDARD":
		co = ParseStandard(cs, startLine)
	case "DAYLIGHT":
		co = ParseDaylight(cs, startLine)
	default:
		co = ParseGeneralComponent(cs, startLine)
	}
	return co, nil
}

func ParseVEvent(cs *CalendarStream, startLine *BaseProperty) *VEvent {
	r, err := ParseComponent(cs, startLine)
	if err != nil {
		return nil
	}
	rr := &VEvent{
		ComponentBase: r,
	}
	return rr
}

func ParseVTodo(cs *CalendarStream, startLine *BaseProperty) *VTodo {
	r, err := ParseComponent(cs, startLine)
	if err != nil {
		return nil
	}
	rr := &VTodo{
		ComponentBase: r,
	}
	return rr
}

func ParseVJournal(cs *CalendarStream, startLine *BaseProperty) *VJournal {
	r, err := ParseComponent(cs, startLine)
	if err != nil {
		return nil
	}
	rr := &VJournal{
		ComponentBase: r,
	}
	return rr
}

func ParseVBusy(cs *CalendarStream, startLine *BaseProperty) *VBusy {
	r, err := ParseComponent(cs, startLine)
	if err != nil {
		return nil
	}
	rr := &VBusy{
		ComponentBase: r,
	}
	return rr
}

func ParseVTimezone(cs *CalendarStream, startLine *BaseProperty) *VTimezone {
	r, err := ParseComponent(cs, startLine)
	if err != nil {
		return nil
	}
	rr := &VTimezone{
		ComponentBase: r,
	}
	return rr
}

func ParseVAlarm(cs *CalendarStream, startLine *BaseProperty) *VAlarm {
	r, err := ParseComponent(cs, startLine)
	if err != nil {
		return nil
	}
	rr := &VAlarm{
		ComponentBase: r,
	}
	return rr
}

func ParseStandard(cs *CalendarStream, startLine *BaseProperty) *Standard {
	r, err := ParseComponent(cs, startLine)
	if err != nil {
		return nil
	}
	rr := &Standard{
		ComponentBase: r,
	}
	return rr
}

func ParseDaylight(cs *CalendarStream, startLine *BaseProperty) *Daylight {
	r, err := ParseComponent(cs, startLine)
	if err != nil {
		return nil
	}
	rr := &Daylight{
		ComponentBase: r,
	}
	return rr
}

func ParseGeneralComponent(cs *CalendarStream, startLine *BaseProperty) *GeneralComponent {
	r, err := ParseComponent(cs, startLine)
	if err != nil {
		return nil
	}
	rr := &GeneralComponent{
		ComponentBase: r,
		Token:         startLine.Value,
	}
	return rr
}

func ParseComponent(cs *CalendarStream, startLine *BaseProperty) (ComponentBase, error) {
	cb := ComponentBase{}
	cont := true
	for i := 0; cont; i++ {
		l, err := cs.ReadLine()
		if err != nil {
			switch err {
			case io.EOF:
				cont = false
			default:
				return cb, err
			}
		}
		if l == nil || len(*l) == 0 {
			continue
		}
		line := ParseProperty(*l)
		if line == nil {
			return cb, errors.New("Error parsing line")
		}
		switch line.IANAToken {
		case "END":
			switch line.Value {
			case startLine.Value:
				return cb, nil
			default:
				return cb, errors.New("Unbalanched end!")
			}
		case "BEGIN":
			co, err := GeneralParseComponent(cs, line)
			if err != nil {
				return cb, err
			}
			if co != nil {
				cb.Components = append(cb.Components, co)
			}
		default: // TODO put in all the supported types for type switching etc.
			cb.Properties = append(cb.Properties, IANAProperty{*line})
		}
	}
	return cb, errors.New("Ran out of lines")
}
