package gosasl

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var (
	krbSPNHost = regexp.MustCompile(`\A[^/]+/(_HOST)([@/]|\z)`)
)

// DEFAULT_MAX_LENGTH is the max length that will be requested in the negotiation
// It can be set with gssapiMechanism.MaxLength = 1000
const DEFAULT_MAX_LENGTH = 16384000

// AUTH if the flag used for just basic auth, no confidentiality
var AUTH = "auth"

// AUTH_INT is the flag for authentication and integrety
var AUTH_INT = "auth-int"

// AUTH_CONF is the flag for authentication and confidentiality. It
// the most secure option.
var AUTH_CONF = "auth-conf"

//QOP_TO_FLAG is a dict that translate the string flag name into the actual bit
// It can be used wiht gssapiMechanism.UserSelectQop = QOP_TO_FLAG[AUTH_CONF] | QOP_TO_FLAG[AUTH_INT]
var QOP_TO_FLAG = map[string]byte{
	AUTH:      1,
	AUTH_INT:  2,
	AUTH_CONF: 4,
}

// QOP is the byte that holds the QOP flags
type QOP []byte

// MechanismConfig is the configuration to use for mechanisms
type MechanismConfig struct {
	name               string
	score              int
	complete           bool
	hasInitialResponse bool
	allowsAnonymous    bool
	usesPlaintext      bool
	activeSafe         bool
	dictionarySafe     bool
	qop                QOP
	// It can be set with mechanism.getConfig().AuthorizationID = "authorizationId"
	AuthorizationID string
}

// Mechanism is the common interface for all mechanisms
type Mechanism interface {
	start() ([]byte, error)
	step(challenge []byte) ([]byte, error)
	encode(outgoing []byte) ([]byte, error)
	decode(incoming []byte) ([]byte, error)
	dispose()
	getConfig() *MechanismConfig
}

// AnonymousMechanism corresponds to NONE/ Anonymous SASL mechanism
type AnonymousMechanism struct {
	config *MechanismConfig
}

// NewAnonymousMechanism returns a new AnonymousMechanism
func NewAnonymousMechanism() *AnonymousMechanism {
	return &AnonymousMechanism{
		config: newDefaultConfig("Anonymous"),
	}
}

func (m *AnonymousMechanism) start() ([]byte, error) {
	return m.step(nil)
}

func (m *AnonymousMechanism) step([]byte) ([]byte, error) {
	m.config.complete = true
	return []byte("Anonymous, None"), nil
}

func (m *AnonymousMechanism) encode([]byte) ([]byte, error) {
	return nil, nil
}

func (m *AnonymousMechanism) decode([]byte) ([]byte, error) {
	return nil, nil
}

func (m *AnonymousMechanism) dispose() {}

func (m *AnonymousMechanism) getConfig() *MechanismConfig {
	return m.config
}

// PlainMechanism corresponds to PLAIN SASL mechanism
type PlainMechanism struct {
	mechanismConfig *MechanismConfig
	identity        string
	username        string
	password        string
}

// NewPlainMechanism returns a new PlainMechanism
func NewPlainMechanism(username string, password string) *PlainMechanism {
	return &PlainMechanism{
		mechanismConfig: newDefaultConfig("PLAIN"),
		username:        username,
		password:        password,
	}
}

func (m *PlainMechanism) start() ([]byte, error) {
	return m.step(nil)
}

func (m *PlainMechanism) step(challenge []byte) ([]byte, error) {
	m.mechanismConfig.complete = true
	var authID string

	if m.mechanismConfig.AuthorizationID != "" {
		authID = m.mechanismConfig.AuthorizationID
	} else {
		authID = m.identity
	}
	NULL := "\x00"
	return []byte(fmt.Sprintf("%s%s%s%s%s", authID, NULL, m.username, NULL, m.password)), nil
}

func (m *PlainMechanism) encode(outgoing []byte) ([]byte, error) {
	return outgoing, nil
}

func (m *PlainMechanism) decode(incoming []byte) ([]byte, error) {
	return incoming, nil
}

func (m *PlainMechanism) dispose() {
	m.password = ""
}

func (m *PlainMechanism) getConfig() *MechanismConfig {
	return m.mechanismConfig
}

// CramMD5Mechanism corresponds to PLAIN SASL mechanism
type CramMD5Mechanism struct {
	*PlainMechanism
}

// NewCramMD5Mechanism returns a new PlainMechanism
func NewCramMD5Mechanism(username string, password string) *CramMD5Mechanism {
	plain := NewPlainMechanism(username, password)
	return &CramMD5Mechanism{
		plain,
	}
}

func (m *CramMD5Mechanism) step(challenge []byte) ([]byte, error) {
	if challenge == nil {
		return nil, nil
	}
	m.mechanismConfig.complete = true
	hash := hmac.New(md5.New, []byte(m.password))
	// hashed := make([]byte, hash.Size())
	_, err := hash.Write(challenge)
	if err != nil {
		return nil, err
	}
	return append([]byte(fmt.Sprintf("%s ", m.username)), hash.Sum(nil)...), nil
}

// DigestMD5Mechanism corresponds to PLAIN SASL mechanism
type DigestMD5Mechanism struct {
	mechanismConfig *MechanismConfig
	service         string
	identity        string
	username        string
	password        string
	host            string
	nonceCount      int
	cnonce          string
	nonce           string
	keyHash         string
	auth            string
}

// parseChallenge turns the challenge string into a map
func parseChallenge(challenge []byte) map[string]string {
	s := string(challenge)

	c := make(map[string]string)

	for len(s) > 0 {
		eq := strings.Index(s, "=")
		key := s[:eq]
		s = s[eq+1:]
		isQuoted := false
		search := ","
		if s[0:1] == "\"" {
			isQuoted = true
			search = "\""
			s = s[1:]
		}
		co := strings.Index(s, search)
		if co == -1 {
			co = len(s)
		}
		val := s[:co]
		if isQuoted && len(s) > len(val)+1 {
			s = s[co+2:]
		} else if co < len(s) {
			s = s[co+1:]
		} else {
			s = ""
		}
		c[key] = val
	}

	return c
}

// NewDigestMD5Mechanism returns a new PlainMechanism
func NewDigestMD5Mechanism(service string, username string, password string) *DigestMD5Mechanism {
	return &DigestMD5Mechanism{
		mechanismConfig: newDefaultConfig("DIGEST-MD5"),
		service:         service,
		username:        username,
		password:        password,
	}
}

func (m *DigestMD5Mechanism) start() ([]byte, error) {
	return m.step(nil)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (m *DigestMD5Mechanism) authenticate(digestUri string, challengeMap map[string]string) error {
	a2String := ":" + digestUri

	if m.auth != "auth" {
		a2String += ":00000000000000000000000000000000"
	}

	if m.getHash(digestUri, a2String, challengeMap) != challengeMap["rspauth"] {
		return fmt.Errorf("authenticate failed")
	}
	return nil
}

func (m *DigestMD5Mechanism) getHash(digestUri string, a2String string, c map[string]string) string {
	// Create a1: HEX(H(H(username:realm:password):nonce:cnonce:authid))
	if m.keyHash == "" {
		x := m.username + ":" + c["realm"] + ":" + m.password
		byteKeyHash := md5.Sum([]byte(x))
		m.keyHash = string(byteKeyHash[:])
	}
	a1String := []string{
		m.keyHash,
		m.nonce,
		m.cnonce,
	}

	if len(m.mechanismConfig.AuthorizationID) != 0 {
		a1String = append(a1String, m.mechanismConfig.AuthorizationID)
	}

	h1 := md5.Sum([]byte(strings.Join(a1String, ":")))
	a1 := hex.EncodeToString(h1[:])

	h2 := md5.Sum([]byte(a2String))
	a2 := hex.EncodeToString(h2[:])

	// Set nonce count nc
	nc := fmt.Sprintf("%08x", m.nonceCount)

	// Create response: H(a1:nonce:nc:cnonce:qop:a2)
	r := strings.ToLower(a1) + ":" + m.nonce + ":" + nc + ":" + m.cnonce + ":" + m.auth + ":" + strings.ToLower(a2)
	hr := md5.Sum([]byte(r))

	// Convert response to hex
	response := strings.ToLower(hex.EncodeToString(hr[:]))
	return string(response)

}

func (m *DigestMD5Mechanism) step(challenge []byte) ([]byte, error) {
	if challenge == nil {
		return nil, nil
	}

	// Create map of challenge
	c := parseChallenge(challenge)
	digestUri := m.service + "/" + m.host

	if _, ok := c["rspauth"]; ok {
		m.mechanismConfig.complete = true
		return nil, m.authenticate(digestUri, c)
	}

	// Prepare response variables
	m.nonce = c["nonce"]
	m.auth = c["qop"]
	if m.nonceCount == 0 {
		m.cnonce = randSeq(14)
	}
	m.nonceCount++

	// Create a2: HEX(H(AUTHENTICATE:digest-uri-value:00000000000000000000000000000000))
	a2String := "AUTHENTICATE:" + digestUri

	maxBuf := ""
	if c["qop"] != AUTH {
		a2String += ":00000000000000000000000000000000"
		maxBuf = ",maxbuf=16777215"
	}
	// Set nonce count nc
	nc := fmt.Sprintf("%08x", m.nonceCount)
	// Create final response sent to server
	resp := "qop=" + c["qop"] + ",realm=" + strconv.Quote(c["realm"]) + ",username=" + strconv.Quote(m.username) + ",nonce=" + strconv.Quote(m.nonce) +
		",cnonce=" + strconv.Quote(m.cnonce) + ",nc=" + nc + ",digest-uri=" + strconv.Quote(digestUri) + ",response=" + m.getHash(digestUri, a2String, c) + maxBuf

	return []byte(resp), nil
}

func (m *DigestMD5Mechanism) encode(outgoing []byte) ([]byte, error) {
	return outgoing, nil
}

func (m *DigestMD5Mechanism) decode(incoming []byte) ([]byte, error) {
	return incoming, nil
}

func (m *DigestMD5Mechanism) dispose() {
	m.password = ""
}

func (m *DigestMD5Mechanism) getConfig() *MechanismConfig {
	return m.mechanismConfig
}

// Client is the entry point for usage of this library
type Client struct {
	host            string
	authorizationID string
	mechanism       Mechanism
}

func newDefaultConfig(name string) *MechanismConfig {
	return &MechanismConfig{
		name:               name,
		score:              0,
		complete:           false,
		hasInitialResponse: false,
		allowsAnonymous:    true,
		usesPlaintext:      true,
		activeSafe:         false,
		dictionarySafe:     false,
		qop:                nil,
		AuthorizationID:    "",
	}
}

// NewSaslClient creates a new client given a host and a mechanism
func NewSaslClient(host string, mechanism Mechanism) *Client {
	mech, ok := mechanism.(*GSSAPIMechanism)
	if ok {
		mech.host = host
	}
	mechDigest, ok := mechanism.(*DigestMD5Mechanism)
	if ok {
		mechDigest.host = host
	}
	return &Client{
		host:      host,
		mechanism: mechanism,
	}
}

// Start initializes the client and may generate the first challenge
func (client *Client) Start() ([]byte, error) {
	return client.mechanism.start()
}

// Step is used for the initial handshake
func (client *Client) Step(challenge []byte) ([]byte, error) {
	return client.mechanism.step(challenge)
}

// Complete returns true if the handshake has ended
func (client *Client) Complete() bool {
	return client.mechanism.getConfig().complete
}

// GetConfig returns the configuration of the mechanism
func (client *Client) GetConfig() *MechanismConfig {
	return client.mechanism.getConfig()
}

// Encode is applied on the outgoing bytes to secure them usually
func (client *Client) Encode(outgoing []byte) ([]byte, error) {
	return client.mechanism.encode(outgoing)
}

// Decode is used on the incoming data to produce the usable bytes
func (client *Client) Decode(incoming []byte) ([]byte, error) {
	return client.mechanism.decode(incoming)
}

// Dispose eliminates sensitive information
func (client *Client) Dispose() {
	client.mechanism.dispose()
}
