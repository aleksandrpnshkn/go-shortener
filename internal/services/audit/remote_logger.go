package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RemoteLogger struct {
	client  *http.Client
	logsURL string
}

func (r *RemoteLogger) SendEntry(ctx context.Context, entry Entry) error {
	rawEntry, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	body := bytes.NewReader(rawEntry)
	req, err := http.NewRequestWithContext(ctx, "POST", r.logsURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("bad error code %d", res.StatusCode)
	}

	return nil
}

func NewRemoteLogger(
	client *http.Client,
	logsURL string,
) *RemoteLogger {
	return &RemoteLogger{
		client:  client,
		logsURL: logsURL,
	}
}
