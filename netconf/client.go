package netconf

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/metajar/netconf/logging"
	"github.com/metajar/netconf/netconf/message"
)

const (
	CISCOTYPE = iota
	ARISTATYPE
)

type Client struct {
	s         *Session
	SessionID string
	target    string
	DevType   int
	timing    time.Time
}

func NewClient(target string, user, pass string, devicetype int) (Client, error) {
	var sess *Session
	var err error
	id, err := uuid.NewV4()

	switch devicetype {
	case 0:
		sess, err = dialXR(target, user, pass)
	case 1:
		sess, err = dialArista(target, user, pass)
	}
	if err != nil {
		return Client{}, err
	}
	cs := DefaultCapabilities
	err = sess.SendHello(&message.Hello{Capabilities: cs})
	if err != nil {
		return Client{}, err
	}
	return Client{
		s:         sess,
		SessionID: id.String(),
		target:    target,
		timing:    time.Now(),
	}, nil
}

func (c *Client) Get(filter string) (string, error) {
	logging.Logger.Infow("Get", "filter", filter, "session", c.SessionID, "target", c.target, "duration", time.Since(c.timing).String())
	return c.executeRPC(message.NewGet(message.FilterTypeSubtree, filter), 30)
}

func (c *Client) GetRunning() (string, error) {
	logging.Logger.Infow("running GetRunning RPC", "session", c.SessionID, "target", c.target, "duration", time.Since(c.timing).String())
	return c.executeRPC(message.NewGetConfig(message.DatastoreRunning, message.FilterTypeSubtree, ""), 30)
}

func (c *Client) Lock() (string, error) {
	logging.Logger.Infow("running Lock RPC", "session", c.SessionID, "target", c.target, "duration", time.Since(c.timing).String())
	return c.executeRPC(message.NewLock(message.DatastoreCandidate), 10)
}

func (c *Client) UnLock() (string, error) {
	logging.Logger.Infow("running Unlock RPC", "session", c.SessionID, "target", c.target, "duration", time.Since(c.timing).String())
	return c.executeRPC(message.NewUnlock(message.DatastoreCandidate), 10)
}

func (c *Client) Commit() (string, error) {
	logging.Logger.Infow("running Commit RPC", "session", c.SessionID, "target", c.target, "duration", time.Since(c.timing).String())
	return c.executeRPC(message.NewCommit(), 10)
}

func (c *Client) Edit(payload string) (string, error) {
	logging.Logger.Infow("running Edit RPC", "session", c.SessionID, "target", c.target, "duration", time.Since(c.timing).String())
	return c.executeRPC(message.NewEditConfig(message.DatastoreCandidate, message.DefaultOperationTypeMerge, payload), 10)
}

func (c *Client) executeRPC(rpc message.RPCMethod, timeout int32) (string, error) {
	response, err := c.s.SyncRPC(rpc, timeout)
	return response.Data, err
}

func (c *Client) Close() (string, error) {
	logging.Logger.Infow("running Close RPC", "session", c.SessionID, "target", c.target, "duration", time.Since(c.timing).String())
	return c.executeRPC(message.NewCloseSession(), 10)
}
