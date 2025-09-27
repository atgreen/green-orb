package discord

import (
	"errors"
	"time"

	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/nicholas-fedor/shoutrrr/pkg/util"
)

const (
	MaxEmbeds = 9
)

// Static error definition.
var ErrEmptyMessage = errors.New("message is empty")

// WebhookPayload is the webhook endpoint payload.
type WebhookPayload struct {
	Embeds    []embedItem `json:"embeds"`
	Username  string      `json:"username,omitempty"`
	AvatarURL string      `json:"avatar_url,omitempty"`
}

// JSON is the actual notification payload.
type embedItem struct {
	Title     string       `json:"title,omitempty"`
	Content   string       `json:"description,omitempty"`
	URL       string       `json:"url,omitempty"`
	Timestamp string       `json:"timestamp,omitempty"`
	Color     uint         `json:"color,omitempty"`
	Footer    *embedFooter `json:"footer,omitempty"`
}

type embedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// CreatePayloadFromItems creates a JSON payload to be sent to the discord webhook API.
func CreatePayloadFromItems(
	items []types.MessageItem,
	title string,
	colors [types.MessageLevelCount]uint,
) (WebhookPayload, error) {
	if len(items) < 1 {
		return WebhookPayload{}, ErrEmptyMessage
	}

	itemCount := util.Min(MaxEmbeds, len(items))

	embeds := make([]embedItem, 0, itemCount)

	for _, item := range items {
		color := uint(0)
		if item.Level >= types.Unknown && int(item.Level) < len(colors) {
			color = colors[item.Level]
		}

		embeddedItem := embedItem{
			Content: item.Text,
			Color:   color,
		}

		if item.Level != types.Unknown {
			embeddedItem.Footer = &embedFooter{
				Text: item.Level.String(),
			}
		}

		if !item.Timestamp.IsZero() {
			embeddedItem.Timestamp = item.Timestamp.UTC().Format(time.RFC3339)
		}

		embeds = append(embeds, embeddedItem)
	}

	// This should not happen, but it's better to leave the index check before dereferencing the array
	if len(embeds) > 0 {
		embeds[0].Title = title
	}

	return WebhookPayload{
		Embeds: embeds,
	}, nil
}
