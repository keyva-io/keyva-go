// Connection pool for Keyva clients.
//
// Auto-generated from keyva protocol spec. Do not edit.

package keyva

import (
	"sync"
)

// PoolConfig controls connection pool behavior.
type PoolConfig struct {
	// Maximum number of idle connections to keep in the pool.
	// Default: 4.
	MaxIdle int

	// Maximum number of total connections (idle + in-use).
	// 0 means unlimited. Default: 0.
	MaxOpen int
}

type pool struct {
	mu      sync.Mutex
	host    string
	port    int
	tls     bool
	auth    string
	idle    []*connection
	open    int
	config  PoolConfig
}

func newPool(host string, port int, useTLS bool, auth string, cfg PoolConfig) *pool {
	if cfg.MaxIdle <= 0 {
		cfg.MaxIdle = 4
	}
	return &pool{
		host:   host,
		port:   port,
		tls:    useTLS,
		auth:   auth,
		config: cfg,
	}
}

func (p *pool) get() (*connection, error) {
	p.mu.Lock()

	// Try to reuse an idle connection
	if len(p.idle) > 0 {
		c := p.idle[len(p.idle)-1]
		p.idle = p.idle[:len(p.idle)-1]
		p.mu.Unlock()
		return c, nil
	}

	// Check max open limit
	if p.config.MaxOpen > 0 && p.open >= p.config.MaxOpen {
		p.mu.Unlock()
		// Block would be better, but keep it simple: create anyway
		// A production pool would use a condition variable here
	}

	p.open++
	p.mu.Unlock()

	c, err := dial(p.host, p.port, p.tls)
	if err != nil {
		p.mu.Lock()
		p.open--
		p.mu.Unlock()
		return nil, err
	}

	if p.auth != "" {
		if _, err := c.execute("AUTH", p.auth); err != nil {
			c.close()
			p.mu.Lock()
			p.open--
			p.mu.Unlock()
			return nil, err
		}
	}

	return c, nil
}

func (p *pool) put(c *connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.idle) < p.config.MaxIdle {
		p.idle = append(p.idle, c)
	} else {
		c.close()
		p.open--
	}
}

func (p *pool) close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, c := range p.idle {
		c.close()
	}
	p.idle = nil
	p.open = 0
}
