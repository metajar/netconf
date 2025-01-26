package netconf

import (
	"time"

	"github.com/lithammer/shortuuid/v4"
	"github.com/metajar/netconf/logging"
	"github.com/metajar/netconf/netconf/message"
)

type Client struct {
	s         *Session
	SessionID string
	target    string
	DevType   int
	timing    time.Time
}

type Options struct {
	keyboardAuth bool
}

type Option func(*Options)

func WithKeyboardAuthentication() Option {
	return func(o *Options) {
		o.keyboardAuth = true
	}
}

func NewClient(target string, user, pass string, opts ...Option) (Client, error) {
	id := shortuuid.New()
	o := &Options{}

	for _, opt := range opts {
		opt(o)
	}
	sess, err := NewNetconfConnection(target, user, pass, o.keyboardAuth)

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
		SessionID: id,
		target:    target,
		timing:    time.Now(),
	}, nil
}

func (c *Client) Get(filter string) (string, error) {
	logging.Logger.Infow("Get", "filter", filter, "session", c.SessionID, "target", c.target, "duration", time.Since(c.timing).String())
	return c.executeRPC(message.NewGet(message.FilterTypeSubtree, filter), 30)
}

func (c *Client) Raw(data string) (string, error) {
	logging.Logger.Infow("Raw", "data", data)
	return c.executeRPC(message.NewRaw(data), 30)
}

func (c *Client) GetConfig(datastore string, filterType, filter string) (string, error) {
	logging.Logger.Infow("Get", "filter", filter, "session", c.SessionID, "target", c.target, "duration", time.Since(c.timing).String())
	if filterType == "" {
		filterType = "subtree"
	}
	return c.executeRPC(message.NewGetConfig(datastore, filterType, filter), 30)
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
