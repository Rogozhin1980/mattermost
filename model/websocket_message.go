// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	WEBSOCKET_EVENT_TYPING                  = "typing"
	WEBSOCKET_EVENT_POSTED                  = "posted"
	WEBSOCKET_EVENT_POST_EDITED             = "post_edited"
	WEBSOCKET_EVENT_POST_DELETED            = "post_deleted"
	WEBSOCKET_EVENT_POST_UNREAD             = "post_unread"
	WEBSOCKET_EVENT_CHANNEL_CONVERTED       = "channel_converted"
	WEBSOCKET_EVENT_CHANNEL_CREATED         = "channel_created"
	WEBSOCKET_EVENT_CHANNEL_DELETED         = "channel_deleted"
	WEBSOCKET_EVENT_CHANNEL_UPDATED         = "channel_updated"
	WEBSOCKET_EVENT_CHANNEL_MEMBER_UPDATED  = "channel_member_updated"
	WEBSOCKET_EVENT_DIRECT_ADDED            = "direct_added"
	WEBSOCKET_EVENT_GROUP_ADDED             = "group_added"
	WEBSOCKET_EVENT_NEW_USER                = "new_user"
	WEBSOCKET_EVENT_ADDED_TO_TEAM           = "added_to_team"
	WEBSOCKET_EVENT_LEAVE_TEAM              = "leave_team"
	WEBSOCKET_EVENT_UPDATE_TEAM             = "update_team"
	WEBSOCKET_EVENT_DELETE_TEAM             = "delete_team"
	WEBSOCKET_EVENT_RESTORE_TEAM            = "restore_team"
	WEBSOCKET_EVENT_USER_ADDED              = "user_added"
	WEBSOCKET_EVENT_USER_UPDATED            = "user_updated"
	WEBSOCKET_EVENT_USER_ROLE_UPDATED       = "user_role_updated"
	WEBSOCKET_EVENT_MEMBERROLE_UPDATED      = "memberrole_updated"
	WEBSOCKET_EVENT_USER_REMOVED            = "user_removed"
	WEBSOCKET_EVENT_PREFERENCE_CHANGED      = "preference_changed"
	WEBSOCKET_EVENT_PREFERENCES_CHANGED     = "preferences_changed"
	WEBSOCKET_EVENT_PREFERENCES_DELETED     = "preferences_deleted"
	WEBSOCKET_EVENT_EPHEMERAL_MESSAGE       = "ephemeral_message"
	WEBSOCKET_EVENT_STATUS_CHANGE           = "status_change"
	WEBSOCKET_EVENT_HELLO                   = "hello"
	WEBSOCKET_AUTHENTICATION_CHALLENGE      = "authentication_challenge"
	WEBSOCKET_EVENT_REACTION_ADDED          = "reaction_added"
	WEBSOCKET_EVENT_REACTION_REMOVED        = "reaction_removed"
	WEBSOCKET_EVENT_RESPONSE                = "response"
	WEBSOCKET_EVENT_EMOJI_ADDED             = "emoji_added"
	WEBSOCKET_EVENT_CHANNEL_VIEWED          = "channel_viewed"
	WEBSOCKET_EVENT_PLUGIN_STATUSES_CHANGED = "plugin_statuses_changed"
	WEBSOCKET_EVENT_PLUGIN_ENABLED          = "plugin_enabled"
	WEBSOCKET_EVENT_PLUGIN_DISABLED         = "plugin_disabled"
	WEBSOCKET_EVENT_ROLE_UPDATED            = "role_updated"
	WEBSOCKET_EVENT_LICENSE_CHANGED         = "license_changed"
	WEBSOCKET_EVENT_CONFIG_CHANGED          = "config_changed"
	WEBSOCKET_EVENT_OPEN_DIALOG             = "open_dialog"
	WEBSOCKET_EVENT_GUESTS_DEACTIVATED      = "guests_deactivated"
)

type WebSocketMessage interface {
	ToJson() string
	IsValid() bool
	EventType() string
}

type WebsocketBroadcast struct {
	OmitUsers             map[string]bool `json:"omit_users"` // broadcast is omitted for users listed here
	UserId                string          `json:"user_id"`    // broadcast only occurs for this user
	ChannelId             string          `json:"channel_id"` // broadcast only occurs for users in this channel
	TeamId                string          `json:"team_id"`    // broadcast only occurs for users in this team
	ContainsSanitizedData bool            `json:"-"`
	ContainsSensitiveData bool            `json:"-"`
}

type precomputedWebSocketEventJSON struct {
	Event     json.RawMessage
	Data      json.RawMessage
	Broadcast json.RawMessage
}

// webSocketEventJSON mirrors WebSocketEvent to make some of its unexported fields serializable
type webSocketEventJSON struct {
	Event     string                 `json:"event"`
	Data      map[string]interface{} `json:"data"`
	Broadcast *WebsocketBroadcast    `json:"broadcast"`
	Sequence  int64                  `json:"seq"`
}

type WebSocketEvent struct {
	event           string
	data            map[string]interface{}
	broadcast       *WebsocketBroadcast
	sequence        int64
	precomputedJSON *precomputedWebSocketEventJSON
}

// PrecomputeJSON precomputes and stores the serialized JSON for all fields other than Sequence.
// This makes ToJson much more efficient when sending the same event to multiple connections.
func (ev *WebSocketEvent) PrecomputeJSON() *WebSocketEvent {
	copy := ev.Copy()
	event, _ := json.Marshal(copy.event)
	data, _ := json.Marshal(copy.data)
	broadcast, _ := json.Marshal(copy.broadcast)
	copy.precomputedJSON = &precomputedWebSocketEventJSON{
		Event:     json.RawMessage(event),
		Data:      json.RawMessage(data),
		Broadcast: json.RawMessage(broadcast),
	}
	return copy
}

func (ev *WebSocketEvent) Add(key string, value interface{}) {
	ev.data[key] = value
}

func NewWebSocketEvent(event, teamId, channelId, userId string, omitUsers map[string]bool) *WebSocketEvent {
	return &WebSocketEvent{event: event, data: make(map[string]interface{}),
		broadcast: &WebsocketBroadcast{TeamId: teamId, ChannelId: channelId, UserId: userId, OmitUsers: omitUsers}}
}

func (ev *WebSocketEvent) Copy() *WebSocketEvent {
	copy := &WebSocketEvent{
		event:           ev.event,
		data:            ev.data,
		broadcast:       ev.broadcast,
		sequence:        ev.sequence,
		precomputedJSON: ev.precomputedJSON,
	}
	return copy
}

func (ev *WebSocketEvent) Data() map[string]interface{} {
	return ev.data
}

func (ev *WebSocketEvent) Broadcast() *WebsocketBroadcast {
	return ev.broadcast
}

func (ev *WebSocketEvent) Sequence() int64 {
	return ev.sequence
}

func (ev *WebSocketEvent) SetEvent(event string) *WebSocketEvent {
	copy := ev.Copy()
	copy.event = event
	return copy
}

func (ev *WebSocketEvent) SetData(data map[string]interface{}) *WebSocketEvent {
	copy := ev.Copy()
	copy.data = data
	return copy
}

func (ev *WebSocketEvent) SetBroadcast(broadcast *WebsocketBroadcast) *WebSocketEvent {
	copy := ev.Copy()
	copy.broadcast = broadcast
	return copy
}

func (ev *WebSocketEvent) SetSequence(seq int64) *WebSocketEvent {
	copy := ev.Copy()
	copy.sequence = seq
	return copy
}

func (ev *WebSocketEvent) IsValid() bool {
	return ev.event != ""
}

func (ev *WebSocketEvent) EventType() string {
	return ev.event
}

func (ev *WebSocketEvent) ToJson() string {
	if ev.precomputedJSON != nil {
		return fmt.Sprintf(`{"event": %s, "data": %s, "broadcast": %s, "seq": %d}`, ev.precomputedJSON.Event, ev.precomputedJSON.Data, ev.precomputedJSON.Broadcast, ev.sequence)
	}
	b, _ := json.Marshal(webSocketEventJSON{
		ev.event,
		ev.data,
		ev.broadcast,
		ev.sequence,
	})
	return string(b)
}

func WebSocketEventFromJson(data io.Reader) *WebSocketEvent {
	var ev WebSocketEvent
	var o webSocketEventJSON
	if err := json.NewDecoder(data).Decode(&o); err != nil {
		return nil
	}
	ev.event = o.Event
	ev.data = o.Data
	ev.broadcast = o.Broadcast
	ev.sequence = o.Sequence
	return &ev
}

type WebSocketResponse struct {
	Status   string                 `json:"status"`
	SeqReply int64                  `json:"seq_reply,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Error    *AppError              `json:"error,omitempty"`
}

func (m *WebSocketResponse) Add(key string, value interface{}) {
	m.Data[key] = value
}

func NewWebSocketResponse(status string, seqReply int64, data map[string]interface{}) *WebSocketResponse {
	return &WebSocketResponse{Status: status, SeqReply: seqReply, Data: data}
}

func NewWebSocketError(seqReply int64, err *AppError) *WebSocketResponse {
	return &WebSocketResponse{Status: STATUS_FAIL, SeqReply: seqReply, Error: err}
}

func (o *WebSocketResponse) IsValid() bool {
	return o.Status != ""
}

func (o *WebSocketResponse) EventType() string {
	return WEBSOCKET_EVENT_RESPONSE
}

func (o *WebSocketResponse) ToJson() string {
	b, _ := json.Marshal(o)
	return string(b)
}

func WebSocketResponseFromJson(data io.Reader) *WebSocketResponse {
	var o *WebSocketResponse
	json.NewDecoder(data).Decode(&o)
	return o
}
