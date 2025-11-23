package audit

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/services/batcher"
)

type remoteBatchSender struct {
	remoteLogger *remoteLogger
}

// GetName возвращает имя.
func (r *remoteBatchSender) GetName() string {
	return "remote_audit_sender"
}

// Execute запускает обработку пачки событий.
func (r *remoteBatchSender) Execute(
	ctx context.Context,
	params []batcher.BatchParam,
) error {
	entries := make([]entry, 0, len(params))

	for _, param := range params {
		entry, ok := param.(entry)
		if !ok {
			return errors.New("passed param is not audit.Entry")
		}
		entries = append(entries, entry)
	}

	return r.remoteLogger.sendEntries(ctx, entries)
}

func newRemoteBatchSender(remoteLogger *remoteLogger) *remoteBatchSender {
	return &remoteBatchSender{
		remoteLogger: remoteLogger,
	}
}
