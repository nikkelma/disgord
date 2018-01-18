package disgord

import (
	"encoding/json"

	"github.com/andersfylling/disgord/lvl"
	"github.com/andersfylling/snowflake"
)

// Guild Guilds in Discord represent an isolated collection of users and channels,
//  and are often referred to as "servers" in the UI.
// https://discordapp.com/developers/docs/resources/guild#guild-object
// Fields with `*` are only sent within the GUILD_CREATE event
type Guild struct {
	ID                          snowflake.ID                   `json:"id,string"`
	ApplicationID               *snowflake.ID                  `json:"application_id"` //   |?
	Name                        string                         `json:"name"`
	Icon                        *string                        `json:"icon"`            //  |?, icon hash
	Splash                      *string                        `json:"splash"`          //  |?, image hash
	Owner                       bool                           `json:"owner,omitempty"` // ?|
	OwnerID                     snowflake.ID                   `json:"owner_id,string"`
	Permissions                 uint64                         `json:"permissions,omitempty"` // ?|, permission flags for connected user `/users/@me/guilds`
	Region                      string                         `json:"region"`
	AfkChannelID                snowflake.ID                   `json:"afk_channel_id,string"`
	AfkTimeout                  uint                           `json:"afk_timeout"`
	EmbedEnabled                bool                           `json:"embed_enabled"`
	EmbedChannelID              snowflake.ID                   `json:"embed_channel_id,string"`
	VerificationLevel           lvl.Verification               `json:"verification_level"`
	DefaultMessageNotifications lvl.DefaultMessageNotification `json:"default_message_notifications"`
	ExplicitContentFilter       lvl.ExplicitContentFilter      `json:"explicit_content_filter"`
	MFALevel                    lvl.MFA                        `json:"mfa_level"`
	WidgetEnabled               bool                           `json:"widget_enabled"`           //   |
	WidgetChannelID             snowflake.ID                   `json:"widget_channel_id,string"` //   |
	Roles                       []*Role                        `json:"roles"`
	Emojis                      []*Emoji                       `json:"emojis"`
	Features                    []string                       `json:"features"`
	SystemChannelID             *snowflake.ID                  `json:"system_channel_id,string,omitempty"` //   |?
	JoinedAt                    DiscordTimestamp               `json:"joined_at,omitempty"`                // ?*|
	Large                       bool                           `json:"large,omitempty"`                    // ?*|
	Unavailable                 bool                           `json:"unavailable"`                        // ?*|
	MemberCount                 uint                           `json:"member_count,omitempty"`             // ?*|
	VoiceStates                 []*VoiceState                  `json:"voice_states,omitempty"`             // ?*|
	Members                     []*GuildMember                 `json:"members,omitempty"`                  // ?*|
	Channels                    []*Channel                     `json:"channels,omitempty"`                 // ?*|
	Presences                   []*Presence                    `json:"presences,omitempty"`                // ?*|
}
type GuildUnavailable struct {
	ID          snowflake.ID `json:"id,string"`
	Unavailable bool         `json:"unavailable"` // ?*|
}

// Compare two guild objects
func (guild *Guild) Compare(g *Guild) bool {
	return (guild == nil && g == nil) || (g != nil && guild.ID == g.ID)
}

func (guild *Guild) MarshalJSON() ([]byte, error) {
	var jsonData []byte
	var err error
	if guild.Unavailable {
		guildUnavailable := GuildUnavailable{ID: guild.ID, Unavailable: true}
		jsonData, err = json.Marshal(&guildUnavailable)
		if err != nil {
			return []byte(""), nil
		}
	} else {
		g := Guild(*guild) // avoid stack overflow by recursive call of Marshal
		jsonData, err = json.Marshal(g)
		if err != nil {
			return []byte(""), nil
		}
	}

	return jsonData, nil
}
