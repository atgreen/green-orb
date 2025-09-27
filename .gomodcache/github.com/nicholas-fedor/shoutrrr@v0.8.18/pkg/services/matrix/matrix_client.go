package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nicholas-fedor/shoutrrr/pkg/types"
	"github.com/nicholas-fedor/shoutrrr/pkg/util"
)

// schemeHTTPPrefixLength is the length of "http" in "https", used to strip TLS suffix.
const (
	schemeHTTPPrefixLength = 4
	tokenHintLength        = 3
	minSliceLength         = 1
	httpClientErrorStatus  = 400
	defaultHTTPTimeout     = 10 * time.Second // defaultHTTPTimeout is the timeout for HTTP requests.
)

// ErrUnsupportedLoginFlows indicates that none of the server login flows are supported.
var (
	ErrUnsupportedLoginFlows = errors.New("none of the server login flows are supported")
	ErrUnexpectedStatus      = errors.New("unexpected HTTP status")
)

// client manages interactions with the Matrix API.
type client struct {
	apiURL      url.URL
	accessToken string
	logger      types.StdLogger
	httpClient  *http.Client
}

// newClient creates a new Matrix client with the specified host and TLS settings.
func newClient(host string, disableTLS bool, logger types.StdLogger) *client {
	client := &client{
		logger: logger,
		apiURL: url.URL{
			Host:   host,
			Scheme: "https",
		},
		httpClient: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
	}

	if client.logger == nil {
		client.logger = util.DiscardLogger
	}

	if disableTLS {
		client.apiURL.Scheme = client.apiURL.Scheme[:schemeHTTPPrefixLength] // "https" -> "http"
	}

	client.logger.Printf("Using server: %v\n", client.apiURL.String())

	return client
}

// useToken sets the access token for the client.
func (c *client) useToken(token string) {
	c.accessToken = token
	c.updateAccessToken()
}

// login authenticates the client using a username and password.
func (c *client) login(user string, password string) error {
	c.apiURL.RawQuery = ""
	defer c.updateAccessToken()

	resLogin := apiResLoginFlows{}
	if err := c.apiGet(apiLogin, &resLogin); err != nil {
		return fmt.Errorf("failed to get login flows: %w", err)
	}

	flows := make([]string, 0, len(resLogin.Flows))
	for _, flow := range resLogin.Flows {
		flows = append(flows, string(flow.Type))

		if flow.Type == flowLoginPassword {
			c.logf("Using login flow '%v'", flow.Type)

			return c.loginPassword(user, password)
		}
	}

	return fmt.Errorf("%w: %v", ErrUnsupportedLoginFlows, strings.Join(flows, ", "))
}

// loginPassword performs a password-based login to the Matrix server.
func (c *client) loginPassword(user string, password string) error {
	response := apiResLogin{}
	if err := c.apiPost(apiLogin, apiReqLogin{
		Type:       flowLoginPassword,
		Password:   password,
		Identifier: newUserIdentifier(user),
	}, &response); err != nil {
		return fmt.Errorf("failed to log in: %w", err)
	}

	c.accessToken = response.AccessToken

	tokenHint := ""
	if len(response.AccessToken) > tokenHintLength {
		tokenHint = response.AccessToken[:tokenHintLength]
	}

	c.logf("AccessToken: %v...\n", tokenHint)
	c.logf("HomeServer: %v\n", response.HomeServer)
	c.logf("User: %v\n", response.UserID)

	return nil
}

// sendMessage sends a message to the specified rooms or all joined rooms if none are specified.
func (c *client) sendMessage(message string, rooms []string) []error {
	if len(rooms) >= minSliceLength {
		return c.sendToExplicitRooms(rooms, message)
	}

	return c.sendToJoinedRooms(message)
}

// sendToExplicitRooms sends a message to explicitly specified rooms and collects any errors.
func (c *client) sendToExplicitRooms(rooms []string, message string) []error {
	var errors []error

	for _, room := range rooms {
		c.logf("Sending message to '%v'...\n", room)

		roomID, err := c.joinRoom(room)
		if err != nil {
			errors = append(errors, fmt.Errorf("error joining room %v: %w", roomID, err))

			continue
		}

		if room != roomID {
			c.logf("Resolved room alias '%v' to ID '%v'", room, roomID)
		}

		if err := c.sendMessageToRoom(message, roomID); err != nil {
			errors = append(
				errors,
				fmt.Errorf("failed to send message to room '%v': %w", roomID, err),
			)
		}
	}

	return errors
}

// sendToJoinedRooms sends a message to all joined rooms and collects any errors.
func (c *client) sendToJoinedRooms(message string) []error {
	var errors []error

	joinedRooms, err := c.getJoinedRooms()
	if err != nil {
		return append(errors, fmt.Errorf("failed to get joined rooms: %w", err))
	}

	for _, roomID := range joinedRooms {
		c.logf("Sending message to '%v'...\n", roomID)

		if err := c.sendMessageToRoom(message, roomID); err != nil {
			errors = append(
				errors,
				fmt.Errorf("failed to send message to room '%v': %w", roomID, err),
			)
		}
	}

	return errors
}

// joinRoom joins a specified room and returns its ID.
func (c *client) joinRoom(room string) (string, error) {
	resRoom := apiResRoom{}
	if err := c.apiPost(fmt.Sprintf(apiRoomJoin, room), nil, &resRoom); err != nil {
		return "", err
	}

	return resRoom.RoomID, nil
}

// sendMessageToRoom sends a message to a specific room.
func (c *client) sendMessageToRoom(message string, roomID string) error {
	resEvent := apiResEvent{}

	return c.apiPost(fmt.Sprintf(apiSendMessage, roomID), apiReqSend{
		MsgType: msgTypeText,
		Body:    message,
	}, &resEvent)
}

// apiGet performs a GET request to the Matrix API.
func (c *client) apiGet(path string, response any) error {
	c.apiURL.Path = path

	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating GET request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing GET request: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading GET response body: %w", err)
	}

	if res.StatusCode >= httpClientErrorStatus {
		resError := &apiResError{}
		if err = json.Unmarshal(body, resError); err == nil {
			return resError
		}

		return fmt.Errorf("%w: %v (unmarshal error: %w)", ErrUnexpectedStatus, res.Status, err)
	}

	if err = json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("unmarshaling GET response: %w", err)
	}

	return nil
}

// apiPost performs a POST request to the Matrix API.
func (c *client) apiPost(path string, request any, response any) error {
	c.apiURL.Path = path

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshaling POST request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.apiURL.String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("creating POST request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing POST request: %w", err)
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading POST response body: %w", err)
	}

	if res.StatusCode >= httpClientErrorStatus {
		resError := &apiResError{}
		if err = json.Unmarshal(body, resError); err == nil {
			return resError
		}

		return fmt.Errorf("%w: %v (unmarshal error: %w)", ErrUnexpectedStatus, res.Status, err)
	}

	if err = json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("unmarshaling POST response: %w", err)
	}

	return nil
}

// updateAccessToken updates the API URL query with the current access token.
func (c *client) updateAccessToken() {
	query := c.apiURL.Query()
	query.Set(accessTokenKey, c.accessToken)
	c.apiURL.RawQuery = query.Encode()
}

// logf logs a formatted message using the client's logger.
func (c *client) logf(format string, v ...any) {
	c.logger.Printf(format, v...)
}

// getJoinedRooms retrieves the list of rooms the client has joined.
func (c *client) getJoinedRooms() ([]string, error) {
	response := apiResJoinedRooms{}
	if err := c.apiGet(apiJoinedRooms, &response); err != nil {
		return []string{}, err
	}

	return response.Rooms, nil
}
