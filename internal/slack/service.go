package slack

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
)

// SlackClient defines the operations for incident Slack management.
type SlackClient interface {
	CreateIncidentChannel(name string) (string, error)
	SetChannelPurpose(channelID, purpose string) error
	PostInitialMessage(channelID string, incID string, title string) error
}

type client struct {
	api *slack.Client
}

// NewClient creates a new Slack client.
func NewClient(token string) SlackClient {
	return &client{
		api: slack.New(token),
	}
}

// CreateIncidentChannel creates a private channel for the incident.
func (c *client) CreateIncidentChannel(name string) (string, error) {
	sanitizedName := Slugify(name)
	
	// Create private channel (isPrivate = true)
	channel, err := c.api.CreateConversation(slack.CreateConversationParams{
		ChannelName: sanitizedName,
		IsPrivate:   true,
	})
	if err != nil {
		if strings.Contains(err.Error(), "invalid_auth") {
			return "", fmt.Errorf("authentication failed: please check your Slack token")
		}
		return "", fmt.Errorf("failed to create channel '%s': %w", sanitizedName, err)
	}

	return channel.ID, nil
}

// SetChannelPurpose sets the purpose of the incident channel.
func (c *client) SetChannelPurpose(channelID, purpose string) error {
	_, err := c.api.SetPurposeOfConversation(channelID, purpose)
	if err != nil {
		return fmt.Errorf("failed to set channel purpose: %w", err)
	}
	return nil
}

// PostInitialMessage sends the initial incident response template.
func (c *client) PostInitialMessage(channelID string, incID string, title string) error {
	attachment := slack.Attachment{
		Color: "#f2c744", // Warning color
		Title: fmt.Sprintf("🚨 Novo Incidente Declarado: %s", title),
		Text:  fmt.Sprintf("ID: %s\n\n*Próximos Passos:*\n1. Identificar o Incident Commander (IC).\n2. Estabelecer canal de comunicação (este).\n3. Investigar a causa raiz.", incID),
		Fields: []slack.AttachmentField{
			{
				Title: "Status",
				Value: "Investigating",
				Short: true,
			},
			{
				Title: "Severity",
				Value: "Unknown",
				Short: true,
			},
		},
	}

	_, _, err := c.api.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		return fmt.Errorf("failed to post initial message: %w", err)
	}

	return nil
}

// Slugify sanitizes a string to be a valid Slack channel name.
func Slugify(s string) string {
	// Lowercase
	s = strings.ToLower(s)
	
	// Replace non-alphanumeric (except hyphens) with hyphens
	reg := regexp.MustCompile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")
	
	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")
	
	// Limit to 80 chars
	if len(s) > 80 {
		s = s[:80]
		s = strings.TrimRight(s, "-")
	}
	
	return s
}
