// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package feed

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	types "main.go/internal/types"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD77e0694DecodeMainGoInternalFeed(in *jlexer.Lexer, out *MsgProperties) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = int64(in.Int64())
		case "data":
			out.Data = string(in.String())
		case "sender":
			out.Sender = types.UserID(in.Int64())
		case "receiver":
			out.Receiver = types.UserID(in.Int64())
		case "time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Time).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed(out *jwriter.Writer, in MsgProperties) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Id))
	}
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		out.String(string(in.Data))
	}
	{
		const prefix string = ",\"sender\":"
		out.RawString(prefix)
		out.Int64(int64(in.Sender))
	}
	{
		const prefix string = ",\"receiver\":"
		out.RawString(prefix)
		out.Int64(int64(in.Receiver))
	}
	{
		const prefix string = ",\"time\":"
		out.RawString(prefix)
		out.Raw((in.Time).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MsgProperties) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MsgProperties) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MsgProperties) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MsgProperties) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed(l, v)
}
func easyjsonD77e0694DecodeMainGoInternalFeed1(in *jlexer.Lexer, out *MessagesToSend) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "messages":
			if in.IsNull() {
				in.Skip()
				out.Messages = nil
			} else {
				in.Delim('[')
				if out.Messages == nil {
					if !in.IsDelim(']') {
						out.Messages = make([]Message, 0, 0)
					} else {
						out.Messages = []Message{}
					}
				} else {
					out.Messages = (out.Messages)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Message
					(v1).UnmarshalEasyJSON(in)
					out.Messages = append(out.Messages, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed1(out *jwriter.Writer, in MessagesToSend) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"messages\":"
		out.RawString(prefix[1:])
		if in.Messages == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Messages {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessagesToSend) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessagesToSend) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessagesToSend) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessagesToSend) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed1(l, v)
}
func easyjsonD77e0694DecodeMainGoInternalFeed2(in *jlexer.Lexer, out *MessageToReceive) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Type":
			out.MsgType = string(in.String())
		case "properties":
			easyjsonD77e0694Decode(in, &out.Properties)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed2(out *jwriter.Writer, in MessageToReceive) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Type\":"
		out.RawString(prefix[1:])
		out.String(string(in.MsgType))
	}
	{
		const prefix string = ",\"properties\":"
		out.RawString(prefix)
		easyjsonD77e0694Encode(out, in.Properties)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MessageToReceive) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MessageToReceive) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MessageToReceive) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MessageToReceive) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed2(l, v)
}
func easyjsonD77e0694Decode(in *jlexer.Lexer, out *struct {
	Id       int64        `json:"id"`
	Data     string       `json:"data"`
	Sender   types.UserID `json:"sender"`
	Receiver types.UserID `json:"receiver"`
	Time     int64        `json:"time"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = int64(in.Int64())
		case "data":
			out.Data = string(in.String())
		case "sender":
			out.Sender = types.UserID(in.Int64())
		case "receiver":
			out.Receiver = types.UserID(in.Int64())
		case "time":
			out.Time = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694Encode(out *jwriter.Writer, in struct {
	Id       int64        `json:"id"`
	Data     string       `json:"data"`
	Sender   types.UserID `json:"sender"`
	Receiver types.UserID `json:"receiver"`
	Time     int64        `json:"time"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Id))
	}
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		out.String(string(in.Data))
	}
	{
		const prefix string = ",\"sender\":"
		out.RawString(prefix)
		out.Int64(int64(in.Sender))
	}
	{
		const prefix string = ",\"receiver\":"
		out.RawString(prefix)
		out.Int64(int64(in.Receiver))
	}
	{
		const prefix string = ",\"time\":"
		out.RawString(prefix)
		out.Int64(int64(in.Time))
	}
	out.RawByte('}')
}
func easyjsonD77e0694DecodeMainGoInternalFeed3(in *jlexer.Lexer, out *Message) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Type":
			out.MsgType = string(in.String())
		case "Properties":
			(out.Properties).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed3(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Type\":"
		out.RawString(prefix[1:])
		out.String(string(in.MsgType))
	}
	{
		const prefix string = ",\"Properties\":"
		out.RawString(prefix)
		(in.Properties).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed3(l, v)
}
func easyjsonD77e0694DecodeMainGoInternalFeed4(in *jlexer.Lexer, out *GetChatRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "person":
			out.Person = types.UserID(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed4(out *jwriter.Writer, in GetChatRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"person\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Person))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetChatRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetChatRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetChatRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetChatRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed4(l, v)
}
func easyjsonD77e0694DecodeMainGoInternalFeed5(in *jlexer.Lexer, out *CreateLikeRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Profile2 = types.UserID(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed5(out *jwriter.Writer, in CreateLikeRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Profile2))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateLikeRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateLikeRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateLikeRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateLikeRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed5(l, v)
}
func easyjsonD77e0694DecodeMainGoInternalFeed6(in *jlexer.Lexer, out *CreateClaimRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "type":
			out.Type = int64(in.Int64())
		case "receiverID":
			out.ReceiverID = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed6(out *jwriter.Writer, in CreateClaimRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Type))
	}
	{
		const prefix string = ",\"receiverID\":"
		out.RawString(prefix)
		out.Int64(int64(in.ReceiverID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CreateClaimRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CreateClaimRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CreateClaimRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CreateClaimRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed6(l, v)
}
func easyjsonD77e0694DecodeMainGoInternalFeed7(in *jlexer.Lexer, out *ClaimsToSend) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "claims":
			if in.IsNull() {
				in.Skip()
				out.Claims = nil
			} else {
				in.Delim('[')
				if out.Claims == nil {
					if !in.IsDelim(']') {
						out.Claims = make([]PureClaim, 0, 2)
					} else {
						out.Claims = []PureClaim{}
					}
				} else {
					out.Claims = (out.Claims)[:0]
				}
				for !in.IsDelim(']') {
					var v4 PureClaim
					easyjsonD77e0694DecodeMainGoInternalFeed8(in, &v4)
					out.Claims = append(out.Claims, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed7(out *jwriter.Writer, in ClaimsToSend) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"claims\":"
		out.RawString(prefix[1:])
		if in.Claims == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Claims {
				if v5 > 0 {
					out.RawByte(',')
				}
				easyjsonD77e0694EncodeMainGoInternalFeed8(out, v6)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ClaimsToSend) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ClaimsToSend) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ClaimsToSend) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ClaimsToSend) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed7(l, v)
}
func easyjsonD77e0694DecodeMainGoInternalFeed8(in *jlexer.Lexer, out *PureClaim) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = int64(in.Int64())
		case "title":
			out.Title = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed8(out *jwriter.Writer, in PureClaim) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Id))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	out.RawByte('}')
}
func easyjsonD77e0694DecodeMainGoInternalFeed9(in *jlexer.Lexer, out *CardsToSend) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "cards":
			if in.IsNull() {
				in.Skip()
				out.Cards = nil
			} else {
				in.Delim('[')
				if out.Cards == nil {
					if !in.IsDelim(']') {
						out.Cards = make([]Card, 0, 0)
					} else {
						out.Cards = []Card{}
					}
				} else {
					out.Cards = (out.Cards)[:0]
				}
				for !in.IsDelim(']') {
					var v7 Card
					easyjsonD77e0694DecodeMainGoInternalFeed10(in, &v7)
					out.Cards = append(out.Cards, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed9(out *jwriter.Writer, in CardsToSend) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cards\":"
		out.RawString(prefix[1:])
		if in.Cards == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Cards {
				if v8 > 0 {
					out.RawByte(',')
				}
				easyjsonD77e0694EncodeMainGoInternalFeed10(out, v9)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CardsToSend) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CardsToSend) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CardsToSend) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CardsToSend) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed9(l, v)
}
func easyjsonD77e0694DecodeMainGoInternalFeed10(in *jlexer.Lexer, out *Card) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = types.UserID(in.Int64())
		case "name":
			out.Name = string(in.String())
		case "birthday":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Birthday).UnmarshalJSON(data))
			}
		case "description":
			out.Description = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "interests":
			if in.IsNull() {
				in.Skip()
				out.Interests = nil
			} else {
				in.Delim('[')
				if out.Interests == nil {
					if !in.IsDelim(']') {
						out.Interests = make([]*Interest, 0, 8)
					} else {
						out.Interests = []*Interest{}
					}
				} else {
					out.Interests = (out.Interests)[:0]
				}
				for !in.IsDelim(']') {
					var v10 *Interest
					if in.IsNull() {
						in.Skip()
						v10 = nil
					} else {
						if v10 == nil {
							v10 = new(Interest)
						}
						easyjsonD77e0694DecodeMainGoInternalFeed11(in, v10)
					}
					out.Interests = append(out.Interests, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "photos":
			if in.IsNull() {
				in.Skip()
				out.Photos = nil
			} else {
				in.Delim('[')
				if out.Photos == nil {
					if !in.IsDelim(']') {
						out.Photos = make([]ImageToSend, 0, 2)
					} else {
						out.Photos = []ImageToSend{}
					}
				} else {
					out.Photos = (out.Photos)[:0]
				}
				for !in.IsDelim(']') {
					var v11 ImageToSend
					easyjsonD77e0694DecodeMainGoInternalFeed12(in, &v11)
					out.Photos = append(out.Photos, v11)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed10(out *jwriter.Writer, in Card) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"birthday\":"
		out.RawString(prefix)
		out.Raw((in.Birthday).MarshalJSON())
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"interests\":"
		out.RawString(prefix)
		if in.Interests == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v12, v13 := range in.Interests {
				if v12 > 0 {
					out.RawByte(',')
				}
				if v13 == nil {
					out.RawString("null")
				} else {
					easyjsonD77e0694EncodeMainGoInternalFeed11(out, *v13)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"photos\":"
		out.RawString(prefix)
		if in.Photos == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v14, v15 := range in.Photos {
				if v14 > 0 {
					out.RawByte(',')
				}
				easyjsonD77e0694EncodeMainGoInternalFeed12(out, v15)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
func easyjsonD77e0694DecodeMainGoInternalFeed12(in *jlexer.Lexer, out *ImageToSend) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "cell":
			out.Cell = string(in.String())
		case "url":
			out.Url = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed12(out *jwriter.Writer, in ImageToSend) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cell\":"
		out.RawString(prefix[1:])
		out.String(string(in.Cell))
	}
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix)
		out.String(string(in.Url))
	}
	out.RawByte('}')
}
func easyjsonD77e0694DecodeMainGoInternalFeed11(in *jlexer.Lexer, out *Interest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ID":
			out.ID = types.InterestID(in.Int64())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed11(out *jwriter.Writer, in Interest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ID\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}
func easyjsonD77e0694DecodeMainGoInternalFeed13(in *jlexer.Lexer, out *AllChats) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "chats":
			if in.IsNull() {
				in.Skip()
				out.Chats = nil
			} else {
				in.Delim('[')
				if out.Chats == nil {
					if !in.IsDelim(']') {
						out.Chats = make([]ChatPreview, 0, 0)
					} else {
						out.Chats = []ChatPreview{}
					}
				} else {
					out.Chats = (out.Chats)[:0]
				}
				for !in.IsDelim(']') {
					var v16 ChatPreview
					easyjsonD77e0694DecodeMainGoInternalFeed14(in, &v16)
					out.Chats = append(out.Chats, v16)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed13(out *jwriter.Writer, in AllChats) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"chats\":"
		out.RawString(prefix[1:])
		if in.Chats == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v17, v18 := range in.Chats {
				if v17 > 0 {
					out.RawByte(',')
				}
				easyjsonD77e0694EncodeMainGoInternalFeed14(out, v18)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AllChats) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD77e0694EncodeMainGoInternalFeed13(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AllChats) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD77e0694EncodeMainGoInternalFeed13(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AllChats) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD77e0694DecodeMainGoInternalFeed13(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AllChats) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD77e0694DecodeMainGoInternalFeed13(l, v)
}
func easyjsonD77e0694DecodeMainGoInternalFeed14(in *jlexer.Lexer, out *ChatPreview) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "personID":
			out.PersonID = int64(in.Int64())
		case "name":
			out.Name = string(in.String())
		case "photo":
			out.Photo = string(in.String())
		case "lastMessage":
			(out.LastMessage).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD77e0694EncodeMainGoInternalFeed14(out *jwriter.Writer, in ChatPreview) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"personID\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.PersonID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"photo\":"
		out.RawString(prefix)
		out.String(string(in.Photo))
	}
	{
		const prefix string = ",\"lastMessage\":"
		out.RawString(prefix)
		(in.LastMessage).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}
