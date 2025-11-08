package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type remoteLogger struct {
	client  *http.Client
	logsURL string
}

func (r *remoteLogger) sendEntries(ctx context.Context, entries []entry) error {
	rawEntries, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	body := bytes.NewReader(rawEntries)
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

// newRemoteLogger cоздаёт новый логгер для отправки событий во внешний сервис аудита.
func newRemoteLogger(
	client *http.Client,
	logsURL string,
) *remoteLogger {
	return &remoteLogger{
		client:  client,
		logsURL: logsURL,
	}
}
