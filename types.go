// Keyva response types.
//
// Auto-generated from keyva protocol spec. Do not edit.

package keyva

// ConfigGetResponse is the response from the CONFIG command.
type ConfigGetResponse struct {
	Value string // Current config value
}

func parseConfigGetResponse(m map[string]any) *ConfigGetResponse {
	r := &ConfigGetResponse{}
	if v, ok := m["value"].(string); ok { r.Value = v }
	return r
}

// HealthResponse is the response from the HEALTH command.
type HealthResponse struct {
	State string // Engine state (e.g. 'ready')
	Keyspaces map[string]any // Per-keyspace credential counts
	Count *int64 // Credential count (keyspace-level)
}

func parseHealthResponse(m map[string]any) *HealthResponse {
	r := &HealthResponse{}
	if v, ok := m["state"].(string); ok { r.State = v }
	if v, ok := m["keyspaces"].(map[string]any); ok { r.Keyspaces = v }
	if v, ok := m["count"].(int64); ok { r.Count = &v }
	return r
}

// InspectResponse is the response from the INSPECT command.
type InspectResponse struct {
	CredentialId string // Credential identifier
	State string // active, suspended, or revoked
	CreatedAt int64 // Creation timestamp
	ExpiresAt *int64 // Expiration timestamp
	LastVerifiedAt *int64 // Last verification timestamp
	Meta map[string]any // Attached metadata
	FamilyId string // Family ID (refresh_token only)
}

func parseInspectResponse(m map[string]any) *InspectResponse {
	r := &InspectResponse{}
	if v, ok := m["credential_id"].(string); ok { r.CredentialId = v }
	if v, ok := m["state"].(string); ok { r.State = v }
	if v, ok := m["created_at"].(int64); ok { r.CreatedAt = v }
	if v, ok := m["expires_at"].(int64); ok { r.ExpiresAt = &v }
	if v, ok := m["last_verified_at"].(int64); ok { r.LastVerifiedAt = &v }
	if v, ok := m["meta"].(map[string]any); ok { r.Meta = v }
	if v, ok := m["family_id"].(string); ok { r.FamilyId = v }
	return r
}

// IssueResponse is the response from the ISSUE command.
type IssueResponse struct {
	CredentialId string // Unique credential identifier
	Token string // The issued token/key
	ExpiresAt *int64 // Expiration timestamp (if TTL set)
	FamilyId string // Refresh token family ID (refresh_token keyspaces only)
}

func parseIssueResponse(m map[string]any) *IssueResponse {
	r := &IssueResponse{}
	if v, ok := m["credential_id"].(string); ok { r.CredentialId = v }
	if v, ok := m["token"].(string); ok { r.Token = v }
	if v, ok := m["expires_at"].(int64); ok { r.ExpiresAt = &v }
	if v, ok := m["family_id"].(string); ok { r.FamilyId = v }
	return r
}

// JwksResponse is the response from the JWKS command.
type JwksResponse struct {
	Jwks string // JSON-encoded JWKS document
}

func parseJwksResponse(m map[string]any) *JwksResponse {
	r := &JwksResponse{}
	if v, ok := m["jwks"].(string); ok { r.Jwks = v }
	return r
}

// KeysResponse is the response from the KEYS command.
type KeysResponse struct {
	Cursor string // Cursor for next page ('0' when complete)
	Keys map[string]any // Array of credential ID strings
}

func parseKeysResponse(m map[string]any) *KeysResponse {
	r := &KeysResponse{}
	if v, ok := m["cursor"].(string); ok { r.Cursor = v }
	if v, ok := m["keys"].(map[string]any); ok { r.Keys = v }
	return r
}

// KeystateResponse is the response from the KEYSTATE command.
type KeystateResponse struct {
	Keys map[string]any // Array of key info maps (key_id, state, created_at)
}

func parseKeystateResponse(m map[string]any) *KeystateResponse {
	r := &KeystateResponse{}
	if v, ok := m["keys"].(map[string]any); ok { r.Keys = v }
	return r
}

// PasswordChangeResponse is the response from the PASSWORD command.
type PasswordChangeResponse struct {
	CredentialId string // Credential identifier
	UpdatedAt int64 // Update timestamp
}

func parsePasswordChangeResponse(m map[string]any) *PasswordChangeResponse {
	r := &PasswordChangeResponse{}
	if v, ok := m["credential_id"].(string); ok { r.CredentialId = v }
	if v, ok := m["updated_at"].(int64); ok { r.UpdatedAt = v }
	return r
}

// PasswordImportResponse is the response from the PASSWORD command.
type PasswordImportResponse struct {
	CredentialId string // Unique credential identifier
	UserId string // User identifier
	Algorithm string // Detected hash algorithm (argon2id, bcrypt, scrypt, etc.)
	CreatedAt int64 // Creation timestamp
}

func parsePasswordImportResponse(m map[string]any) *PasswordImportResponse {
	r := &PasswordImportResponse{}
	if v, ok := m["credential_id"].(string); ok { r.CredentialId = v }
	if v, ok := m["user_id"].(string); ok { r.UserId = v }
	if v, ok := m["algorithm"].(string); ok { r.Algorithm = v }
	if v, ok := m["created_at"].(int64); ok { r.CreatedAt = v }
	return r
}

// PasswordSetResponse is the response from the PASSWORD command.
type PasswordSetResponse struct {
	CredentialId string // Unique credential identifier
	UserId string // User identifier
	Algorithm string // Hash algorithm used
	CreatedAt int64 // Creation timestamp
}

func parsePasswordSetResponse(m map[string]any) *PasswordSetResponse {
	r := &PasswordSetResponse{}
	if v, ok := m["credential_id"].(string); ok { r.CredentialId = v }
	if v, ok := m["user_id"].(string); ok { r.UserId = v }
	if v, ok := m["algorithm"].(string); ok { r.Algorithm = v }
	if v, ok := m["created_at"].(int64); ok { r.CreatedAt = v }
	return r
}

// PasswordVerifyResponse is the response from the PASSWORD command.
type PasswordVerifyResponse struct {
	Valid bool // Whether the password is correct
	CredentialId string // Credential identifier
	Metadata map[string]any // Attached metadata
}

func parsePasswordVerifyResponse(m map[string]any) *PasswordVerifyResponse {
	r := &PasswordVerifyResponse{}
	r.Valid = m["valid"]
	if v, ok := m["credential_id"].(string); ok { r.CredentialId = v }
	if v, ok := m["metadata"].(map[string]any); ok { r.Metadata = v }
	return r
}

// RefreshResponse is the response from the REFRESH command.
type RefreshResponse struct {
	CredentialId string // New credential identifier
	Token string // New refresh token
	FamilyId string // Family ID (unchanged)
	ExpiresAt int64 // New token expiration
}

func parseRefreshResponse(m map[string]any) *RefreshResponse {
	r := &RefreshResponse{}
	if v, ok := m["credential_id"].(string); ok { r.CredentialId = v }
	if v, ok := m["token"].(string); ok { r.Token = v }
	if v, ok := m["family_id"].(string); ok { r.FamilyId = v }
	if v, ok := m["expires_at"].(int64); ok { r.ExpiresAt = v }
	return r
}

// RevokeResponse is the response from the REVOKE command.
type RevokeResponse struct {
	Revoked *int64 // Number of credentials revoked
}

func parseRevokeResponse(m map[string]any) *RevokeResponse {
	r := &RevokeResponse{}
	if v, ok := m["revoked"].(int64); ok { r.Revoked = &v }
	return r
}

// RevokeBulkResponse is the response from the REVOKE command.
type RevokeBulkResponse struct {
	Revoked int64 // Number of credentials revoked
}

func parseRevokeBulkResponse(m map[string]any) *RevokeBulkResponse {
	r := &RevokeBulkResponse{}
	if v, ok := m["revoked"].(int64); ok { r.Revoked = v }
	return r
}

// RevokeFamilyResponse is the response from the REVOKE command.
type RevokeFamilyResponse struct {
	Revoked int64 // Number of credentials revoked
}

func parseRevokeFamilyResponse(m map[string]any) *RevokeFamilyResponse {
	r := &RevokeFamilyResponse{}
	if v, ok := m["revoked"].(int64); ok { r.Revoked = v }
	return r
}

// RotateResponse is the response from the ROTATE command.
type RotateResponse struct {
	NewKeyId string // ID of the newly created key
	OldKeyId string // ID of the key that entered drain mode
	Dryrun string // 'true' if this was a dry run
}

func parseRotateResponse(m map[string]any) *RotateResponse {
	r := &RotateResponse{}
	if v, ok := m["new_key_id"].(string); ok { r.NewKeyId = v }
	if v, ok := m["old_key_id"].(string); ok { r.OldKeyId = v }
	if v, ok := m["dryrun"].(string); ok { r.Dryrun = v }
	return r
}

// SchemaResponse is the response from the SCHEMA command.
type SchemaResponse struct {
	Schema map[string]any // Schema definition (fields, types, constraints)
}

func parseSchemaResponse(m map[string]any) *SchemaResponse {
	r := &SchemaResponse{}
	if v, ok := m["schema"].(map[string]any); ok { r.Schema = v }
	return r
}

// VerifyResponse is the response from the VERIFY command.
type VerifyResponse struct {
	CredentialId string // Credential identifier
	Claims map[string]any // Decoded JWT claims (JWT keyspaces only)
	Meta map[string]any // Attached metadata
	State string // Credential state (active, suspended)
}

func parseVerifyResponse(m map[string]any) *VerifyResponse {
	r := &VerifyResponse{}
	if v, ok := m["credential_id"].(string); ok { r.CredentialId = v }
	if v, ok := m["claims"].(map[string]any); ok { r.Claims = v }
	if v, ok := m["meta"].(map[string]any); ok { r.Meta = v }
	if v, ok := m["state"].(string); ok { r.State = v }
	return r
}

