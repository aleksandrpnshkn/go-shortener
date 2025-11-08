package audit

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/services/batcher"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

type DeleteCode struct {
	Code types.Code
	User users.User
}

type RemoteBatchSender struct {
	remoteLogger *RemoteLogger
}

func (r *RemoteBatchSender) GetName() string {
	return "remote_audit_sender"
}

func (r *RemoteBatchSender) Execute(
	ctx context.Context,
	params []batcher.BatchParam,
) error {
	entries := make([]Entry, 0, len(params))

	for _, param := range params {
		entry, ok := param.(Entry)
		if !ok {
			return errors.New("passed param is not audit.Entry")
		}
		entries = append(entries, entry)
	}

	return r.remoteLogger.SendEntries(ctx, entries)
}

func NewRemoteBatchSender(remoteLogger *RemoteLogger) *RemoteBatchSender {
	return &RemoteBatchSender{
		remoteLogger: remoteLogger,
	}
}
