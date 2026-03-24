// Package keyva provides a client for the Keyva Credential management server.
//
// Auto-generated from keyva protocol spec. Do not edit.
//
// Connect using a Keyva URI:
//
//	client, err := keyva.Connect("keyva://://localhost")
//	result, err := client.Issue(ctx, "my-keyspace", &keyva.IssueOptions{TTL: 3600})
//	fmt.Println(result.CredentialID, result.Token)
//	client.Close()
package keyva

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Client is a Keyva client backed by a connection pool.
type Client struct {
	pool *pool
}

// Connect creates a new Client from a Keyva URI.
//
// Supported URI formats:
//
//	keyva://://localhost
//	keyva://://localhost:6399
//	keyva://+tls://prod.example.com
//	keyva://://mytoken@localhost:6399
//	keyva://://mytoken@localhost/sessions
//	keyva://+tls://tok@host:6399/keys
func Connect(uri string, opts ...PoolConfig) (*Client, error) {
	cfg, err := parseURI(uri)
	if err != nil {
		return nil, err
	}
	var poolCfg PoolConfig
	if len(opts) > 0 {
		poolCfg = opts[0]
	}
	p := newPool(cfg.host, cfg.port, cfg.tls, cfg.authToken, poolCfg)
	// Verify connectivity by getting and returning one connection
	c, err := p.get()
	if err != nil {
		return nil, err
	}
	p.put(c)
	return &Client{pool: p}, nil
}

// Close shuts down the client and all pooled connections.
func (c *Client) Close() {
	c.pool.close()
}

func (c *Client) exec(args ...string) (any, error) {
	conn, err := c.pool.get()
	if err != nil {
		return nil, err
	}
	result, err := conn.execute(args...)
	if err != nil {
		conn.close()
		return nil, err
	}
	c.pool.put(conn)
	return result, nil
}

func (c *Client) execMap(args ...string) (map[string]any, error) {
	result, err := c.exec(args...)
	if err != nil {
		return nil, err
	}
	m, ok := result.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("keyva: expected map response, got %T", result)
	}
	return m, nil
}

type uriConfig struct {
	host      string
	port      int
	tls       bool
	authToken string
	keyspace  string
}

func parseURI(uri string) (*uriConfig, error) {
	cfg := &uriConfig{port: defaultPort}

	switch {
	case strings.HasPrefix(uri, "keyva://+tls://"):
		cfg.tls = true
		uri = "keyva://://" + uri[len("keyva://+tls://"):]
	case strings.HasPrefix(uri, "keyva://://"):
		// ok
	default:
		return nil, fmt.Errorf("keyva: invalid URI scheme (expected keyva://:// or keyva://+tls://): %s", uri)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("keyva: invalid URI: %w", err)
	}

	cfg.host = u.Hostname()
	if cfg.host == "" {
		cfg.host = "localhost"
	}

	if p := u.Port(); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			cfg.port = n
		}
	}

	if u.User != nil {
		cfg.authToken = u.User.Username()
	}

	cfg.keyspace = strings.TrimPrefix(u.Path, "/")

	return cfg, nil
}

// Auth — Authenticate the current connection
func (c *Client) Auth(token string) error {
	args := []string{
		"AUTH",
		token,
	}
	_, err := c.exec(args...)
	return err
}

// ConfigGet — Retrieve a runtime configuration value
func (c *Client) ConfigGet(key string) (*ConfigGetResponse, error) {
	args := []string{
		"CONFIG", "GET",
		key,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseConfigGetResponse(m), nil
}

// ConfigSet — Set a runtime configuration value
func (c *Client) ConfigSet(key string, value string) error {
	args := []string{
		"CONFIG", "SET",
		key,
		value,
	}
	_, err := c.exec(args...)
	return err
}

// Health — Check server or keyspace health
func (c *Client) Health(keyspace string) (*HealthResponse, error) {
	args := []string{
		"HEALTH",
	}
	if keyspace != "" {
		args = append(args, keyspace)
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseHealthResponse(m), nil
}

// Inspect — Retrieve full details about a credential
func (c *Client) Inspect(keyspace string, credential_id string) (*InspectResponse, error) {
	args := []string{
		"INSPECT",
		keyspace,
		credential_id,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseInspectResponse(m), nil
}

// IssueOptions are optional parameters for Issue.
type IssueOptions struct {
	Claims map[string]any
	Metadata map[string]any
	TtlSecs int64
	IdempotencyKey string
}

// Issue — Issue a new credential in the given keyspace
func (c *Client) Issue(keyspace string, opts *IssueOptions) (*IssueResponse, error) {
	args := []string{
		"ISSUE",
		keyspace,
	}
	if opts != nil {
		if opts.Claims != nil {
			b, _ := json.Marshal(opts.Claims)
			args = append(args, "CLAIMS", string(b))
		}
		if opts.Metadata != nil {
			b, _ := json.Marshal(opts.Metadata)
			args = append(args, "META", string(b))
		}
		if opts.TtlSecs != 0 {
			args = append(args, "TTL", strconv.FormatInt(opts.TtlSecs, 10))
		}
		if opts.IdempotencyKey != "" {
			args = append(args, "IDEMPOTENCY_KEY", opts.IdempotencyKey)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseIssueResponse(m), nil
}

// Jwks — Return the JSON Web Key Set for a JWT keyspace
func (c *Client) Jwks(keyspace string) (*JwksResponse, error) {
	args := []string{
		"JWKS",
		keyspace,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseJwksResponse(m), nil
}

// KeysOptions are optional parameters for Keys.
type KeysOptions struct {
	Cursor string
	Pattern string
	StateFilter string
	Count int64
}

// Keys — List credential IDs with optional filtering and pagination
func (c *Client) Keys(keyspace string, opts *KeysOptions) (*KeysResponse, error) {
	args := []string{
		"KEYS",
		keyspace,
	}
	if opts != nil {
		if opts.Cursor != "" {
			args = append(args, "CURSOR", opts.Cursor)
		}
		if opts.Pattern != "" {
			args = append(args, "MATCH", opts.Pattern)
		}
		if opts.StateFilter != "" {
			args = append(args, "STATE", opts.StateFilter)
		}
		if opts.Count != 0 {
			args = append(args, "COUNT", strconv.FormatInt(opts.Count, 10))
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseKeysResponse(m), nil
}

// Keystate — Show the current key ring state for a keyspace
func (c *Client) Keystate(keyspace string) (*KeystateResponse, error) {
	args := []string{
		"KEYSTATE",
		keyspace,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseKeystateResponse(m), nil
}

// PasswordChange — Change a user's password (requires old password)
func (c *Client) PasswordChange(keyspace string, user_id string, old_password string, new_password string) (*PasswordChangeResponse, error) {
	args := []string{
		"PASSWORD", "CHANGE",
		keyspace,
		user_id,
		old_password,
		new_password,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parsePasswordChangeResponse(m), nil
}

// PasswordImportOptions are optional parameters for PasswordImport.
type PasswordImportOptions struct {
	Metadata map[string]any
}

// PasswordImport — Import a pre-hashed password for migration from another system (argon2, bcrypt, scrypt)
func (c *Client) PasswordImport(keyspace string, user_id string, hash string, opts *PasswordImportOptions) (*PasswordImportResponse, error) {
	args := []string{
		"PASSWORD", "IMPORT",
		keyspace,
		user_id,
		hash,
	}
	if opts != nil {
		if opts.Metadata != nil {
			b, _ := json.Marshal(opts.Metadata)
			args = append(args, "META", string(b))
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parsePasswordImportResponse(m), nil
}

// PasswordSetOptions are optional parameters for PasswordSet.
type PasswordSetOptions struct {
	Metadata map[string]any
}

// PasswordSet — Set a password for a user in a password keyspace
func (c *Client) PasswordSet(keyspace string, user_id string, password string, opts *PasswordSetOptions) (*PasswordSetResponse, error) {
	args := []string{
		"PASSWORD", "SET",
		keyspace,
		user_id,
		password,
	}
	if opts != nil {
		if opts.Metadata != nil {
			b, _ := json.Marshal(opts.Metadata)
			args = append(args, "META", string(b))
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parsePasswordSetResponse(m), nil
}

// PasswordVerify — Verify a user's password
func (c *Client) PasswordVerify(keyspace string, user_id string, password string) (*PasswordVerifyResponse, error) {
	args := []string{
		"PASSWORD", "VERIFY",
		keyspace,
		user_id,
		password,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parsePasswordVerifyResponse(m), nil
}

// Refresh — Exchange a refresh token for a new one
func (c *Client) Refresh(keyspace string, token string) (*RefreshResponse, error) {
	args := []string{
		"REFRESH",
		keyspace,
		token,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseRefreshResponse(m), nil
}

// Revoke — Revoke a credential by ID
func (c *Client) Revoke(keyspace string, credential_id string) (*RevokeResponse, error) {
	args := []string{
		"REVOKE",
		keyspace,
		credential_id,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseRevokeResponse(m), nil
}

// RevokeBulkOptions are optional parameters for RevokeBulk.
type RevokeBulkOptions struct {
	Ids []string
}

// RevokeBulk — Bulk-revoke multiple credentials
func (c *Client) RevokeBulk(keyspace string, opts *RevokeBulkOptions) (*RevokeBulkResponse, error) {
	args := []string{
		"REVOKE",
		keyspace,
	}
	if opts != nil {
		if len(opts.Ids) > 0 {
			args = append(args, "BULK")
			args = append(args, opts.Ids...)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseRevokeBulkResponse(m), nil
}

// RevokeFamilyOptions are optional parameters for RevokeFamily.
type RevokeFamilyOptions struct {
	FamilyId string
}

// RevokeFamily — Revoke all credentials in a refresh token family
func (c *Client) RevokeFamily(keyspace string, opts *RevokeFamilyOptions) (*RevokeFamilyResponse, error) {
	args := []string{
		"REVOKE",
		keyspace,
	}
	if opts != nil {
		if opts.FamilyId != "" {
			args = append(args, "FAMILY", opts.FamilyId)
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseRevokeFamilyResponse(m), nil
}

// RotateOptions are optional parameters for Rotate.
type RotateOptions struct {
	Force bool
	Nowait bool
	Dryrun bool
}

// Rotate — Trigger signing key rotation for a keyspace
func (c *Client) Rotate(keyspace string, opts *RotateOptions) (*RotateResponse, error) {
	args := []string{
		"ROTATE",
		keyspace,
	}
	if opts != nil {
		if opts.Force {
			args = append(args, "FORCE")
		}
		if opts.Nowait {
			args = append(args, "NOWAIT")
		}
		if opts.Dryrun {
			args = append(args, "DRYRUN")
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseRotateResponse(m), nil
}

// Schema — Display the metadata schema for a keyspace
func (c *Client) Schema(keyspace string) (*SchemaResponse, error) {
	args := []string{
		"SCHEMA",
		keyspace,
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseSchemaResponse(m), nil
}

// Subscribe is not yet supported in the Go client (requires streaming).

// Suspend — Temporarily suspend a credential
func (c *Client) Suspend(keyspace string, credential_id string) error {
	args := []string{
		"SUSPEND",
		keyspace,
		credential_id,
	}
	_, err := c.exec(args...)
	return err
}

// Unsuspend — Reactivate a previously suspended credential
func (c *Client) Unsuspend(keyspace string, credential_id string) error {
	args := []string{
		"UNSUSPEND",
		keyspace,
		credential_id,
	}
	_, err := c.exec(args...)
	return err
}

// UpdateOptions are optional parameters for Update.
type UpdateOptions struct {
	Metadata map[string]any
}

// Update — Update metadata on an existing credential
func (c *Client) Update(keyspace string, credential_id string, opts *UpdateOptions) error {
	args := []string{
		"UPDATE",
		keyspace,
		credential_id,
	}
	if opts != nil {
		if opts.Metadata != nil {
			b, _ := json.Marshal(opts.Metadata)
			args = append(args, "META", string(b))
		}
	}
	_, err := c.exec(args...)
	return err
}

// VerifyOptions are optional parameters for Verify.
type VerifyOptions struct {
	Payload string
	CheckRevoked bool
}

// Verify — Verify a credential (JWT, API key, or HMAC signature)
func (c *Client) Verify(keyspace string, token string, opts *VerifyOptions) (*VerifyResponse, error) {
	args := []string{
		"VERIFY",
		keyspace,
		token,
	}
	if opts != nil {
		if opts.Payload != "" {
			args = append(args, "PAYLOAD", opts.Payload)
		}
		if opts.CheckRevoked {
			args = append(args, "CHECKREV")
		}
	}
	m, err := c.execMap(args...)
	if err != nil {
		return nil, err
	}
	return parseVerifyResponse(m), nil
}
