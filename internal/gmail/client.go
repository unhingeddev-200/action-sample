package gmail

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

const apiBase = "https://gmail.googleapis.com/gmail/v1"

// Client calls Gmail REST with a bearer access token.
type Client struct {
	HTTP  *http.Client
	Token string
}

func NewClient(token string) *Client {
	return &Client{
		HTTP:  &http.Client{Timeout: 60 * time.Second},
		Token: strings.TrimSpace(token),
	}
}

func normalizeUserID(userID string) string {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "me"
	}
	return userID
}

// BuildRawMessage builds a Gmail API messages.send "raw" payload (base64url, no pad).
func BuildRawMessage(from, to, subject, textBody, htmlBody string) (string, error) {
	to = strings.TrimSpace(to)
	if to == "" {
		return "", fmt.Errorf("to is required")
	}
	if _, err := mail.ParseAddress(to); err != nil {
		return "", fmt.Errorf("to: %w", err)
	}
	from = strings.TrimSpace(from)
	if from != "" {
		if _, err := mail.ParseAddress(from); err != nil {
			return "", fmt.Errorf("from: %w", err)
		}
	}
	subject = strings.TrimSpace(subject)
	textBody = strings.TrimSpace(textBody)
	htmlBody = strings.TrimSpace(htmlBody)
	if textBody == "" && htmlBody == "" {
		return "", fmt.Errorf("textBody or htmlBody is required")
	}

	var buf bytes.Buffer
	if from != "" {
		fmt.Fprintf(&buf, "From: %s\r\n", from)
	}
	fmt.Fprintf(&buf, "To: %s\r\n", to)
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")

	switch {
	case textBody != "" && htmlBody != "":
		boundary := "genai_boundary_7a3f"
		fmt.Fprintf(&buf, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", boundary)
		fmt.Fprintf(&buf, "--%s\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s\r\n", boundary, textBody)
		fmt.Fprintf(&buf, "--%s\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s\r\n", boundary, htmlBody)
		fmt.Fprintf(&buf, "--%s--\r\n", boundary)
	case htmlBody != "":
		fmt.Fprintf(&buf, "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s", htmlBody)
	default:
		fmt.Fprintf(&buf, "Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s", textBody)
	}

	encoded := base64.RawURLEncoding.EncodeToString(buf.Bytes())
	return encoded, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, query map[string]string, body any, out any) error {
	if c == nil || c.Token == "" {
		return fmt.Errorf("accessToken is required")
	}
	u, err := url.Parse(apiBase + path)
	if err != nil {
		return err
	}
	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			if strings.TrimSpace(v) == "" {
				continue
			}
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
	var bodyReader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gmail api %s %s: %s: %s", method, path, resp.Status, strings.TrimSpace(string(respBody)))
	}
	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode gmail response: %w", err)
	}
	return nil
}

type MessageRef struct {
	ID       string   `json:"id"`
	ThreadID string   `json:"threadId,omitempty"`
	LabelIDs []string `json:"labelIds,omitempty"`
}

type ListResult struct {
	Messages           []MessageRef `json:"messages"`
	NextPageToken      string       `json:"nextPageToken,omitempty"`
	ResultSizeEstimate int64        `json:"resultSizeEstimate,omitempty"`
}

type Message struct {
	ID           string          `json:"id"`
	ThreadID     string          `json:"threadId,omitempty"`
	LabelIDs     []string        `json:"labelIds,omitempty"`
	Snippet      string          `json:"snippet,omitempty"`
	InternalDate string          `json:"internalDate,omitempty"`
	Payload      json.RawMessage `json:"payload,omitempty"`
	SizeEstimate int64           `json:"sizeEstimate,omitempty"`
	Raw          string          `json:"raw,omitempty"`
}

func (c *Client) Send(ctx context.Context, userID, raw string) (MessageRef, error) {
	userID = normalizeUserID(userID)
	var out MessageRef
	err := c.doJSON(ctx, http.MethodPost, "/users/"+userID+"/messages/send", nil, map[string]string{"raw": raw}, &out)
	return out, err
}

func (c *Client) CreateDraft(ctx context.Context, userID, raw string) (map[string]any, error) {
	userID = normalizeUserID(userID)
	var out map[string]any
	err := c.doJSON(ctx, http.MethodPost, "/users/"+userID+"/drafts", nil, map[string]any{
		"message": map[string]string{"raw": raw},
	}, &out)
	return out, err
}

func (c *Client) List(ctx context.Context, userID, q string, maxResults int64, pageToken string) (ListResult, error) {
	userID = normalizeUserID(userID)
	if maxResults <= 0 {
		maxResults = 25
	}
	query := map[string]string{
		"maxResults": fmt.Sprintf("%d", maxResults),
	}
	if strings.TrimSpace(q) != "" {
		query["q"] = q
	}
	if strings.TrimSpace(pageToken) != "" {
		query["pageToken"] = pageToken
	}
	var out ListResult
	err := c.doJSON(ctx, http.MethodGet, "/users/"+userID+"/messages", query, nil, &out)
	if out.Messages == nil {
		out.Messages = []MessageRef{}
	}
	return out, err
}

func (c *Client) Get(ctx context.Context, userID, messageID, format string) (Message, error) {
	userID = normalizeUserID(userID)
	messageID = strings.TrimSpace(messageID)
	if messageID == "" {
		return Message{}, fmt.Errorf("messageId is required")
	}
	format = strings.TrimSpace(format)
	if format == "" {
		format = "full"
	}
	var out Message
	err := c.doJSON(ctx, http.MethodGet, "/users/"+userID+"/messages/"+messageID, map[string]string{"format": format}, nil, &out)
	return out, err
}
