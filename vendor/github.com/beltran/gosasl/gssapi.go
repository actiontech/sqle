// +build kerberos

package gosasl

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/beltran/gssapi"
)

// GSSAPIMechanism corresponds to GSSAPI SASL mechanism
type GSSAPIMechanism struct {
	config           *MechanismConfig
	host             string
	user             string
	service          string
	negotiationStage int
	context          *GSSAPIContext
	qop              byte
	supportedQop     uint8
	serverMaxLength  int
	UserSelectQop    uint8
	MaxLength        int
}

// NewGSSAPIMechanism returns a new GSSAPIMechanism
func NewGSSAPIMechanism(service string) (mechanism *GSSAPIMechanism, err error) {
	context := newGSSAPIContext()
	mechanism = &GSSAPIMechanism{
		config:           newDefaultConfig("GSSAPI"),
		service:          service,
		negotiationStage: 0,
		context:          context,
		supportedQop:     QOP_TO_FLAG[AUTH] | QOP_TO_FLAG[AUTH_CONF] | QOP_TO_FLAG[AUTH_INT],
		MaxLength:        DEFAULT_MAX_LENGTH,
		UserSelectQop:    QOP_TO_FLAG[AUTH] | QOP_TO_FLAG[AUTH_INT] | QOP_TO_FLAG[AUTH_CONF],
	}
	return
}

func (m *GSSAPIMechanism) start() ([]byte, error) {
	return m.step(nil)
}

func (m *GSSAPIMechanism) step(challenge []byte) ([]byte, error) {
	var serviceHostQualified string
	var fullServiceName string
	// Allows to use a service principal designated for another host to still be used.
	// Useful for containerized environments.
	serviceHostQualified = os.Getenv("SERVICE_HOST_QUALIFIED")
	if len(serviceHostQualified) > 0 {
		fullServiceName = m.service + "/" + serviceHostQualified
	} else {
		fullServiceName = m.service + "/" + m.host
	}

	if m.negotiationStage == 0 {
		err := initClientContext(m.context, fullServiceName, nil)
		m.negotiationStage = 1
		return m.context.token, err

	} else if m.negotiationStage == 1 {
		err := initClientContext(m.context, fullServiceName, challenge)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		var srcName *gssapi.Name
		if m.context.contextId != nil {
			srcName, _, _, _, _, _, _, _ = m.context.contextId.InquireContext()
			if srcName != nil {
				m.user = srcName.String()
			}
		}
		if m.user != "" {
			// Check if the context is available. If the user has set the flags
			// it will fail, although at this point we could know that the negotiation won't succeed
			if !m.context.integAvail() && !m.context.confAvail() {
				log.Println("No security layer can be established, authentication is still possible")
			}
			m.negotiationStage = 2
		}
		return m.context.token, nil
	} else if m.negotiationStage == 2 {
		data, err := m.context.unwrap(challenge)
		if err != nil {
			return nil, err
		}
		if len(data) != 4 {
			return nil, fmt.Errorf("Decoded data should have length for at this stage")
		}
		qopBits := data[0]
		data[0] = 0
		m.serverMaxLength = int(binary.BigEndian.Uint32(data))

		m.qop, err = m.selectQop(qopBits)
		// The client doesn't support or want any of the security layers offered by the server
		if err != nil {
			m.MaxLength = 0
		}

		header := make([]byte, 4)
		maxLength := m.serverMaxLength
		if m.MaxLength < m.serverMaxLength {
			maxLength = m.MaxLength
		}

		headerInt := (uint(m.qop) << 24) | uint(maxLength)

		binary.BigEndian.PutUint32(header, uint32(headerInt))

		// FLAG_BYTE + 3 bytes of length + user or authority
		var name string
		if name = m.user; m.config.AuthorizationID != "" {
			name = m.config.AuthorizationID
		}
		out := append(header, []byte(name)...)
		wrappedOut, err := m.context.wrap(out, false)

		m.config.complete = true
		return wrappedOut, err
	}
	return nil, fmt.Errorf("Error, this code should be unreachable")
}

func (m *GSSAPIMechanism) selectQop(qopByte byte) (byte, error) {
	availableQops := m.UserSelectQop & m.supportedQop & qopByte
	for _, qop := range []byte{QOP_TO_FLAG[AUTH_CONF], QOP_TO_FLAG[AUTH_INT], QOP_TO_FLAG[AUTH]} {
		if qop&availableQops != 0 {
			return qop, nil
		}
	}
	return byte(0), fmt.Errorf("No qop satisfying all the conditions where found")
}

// replaceSPNHostWildcard substitutes the special string '_HOST' in the given
// SPN for the given (current) host.
func replaceSPNHostWildcard(spn, host string) string {
	res := krbSPNHost.FindStringSubmatchIndex(spn)
	if res == nil || res[2] == -1 {
		return spn
	}
	return spn[:res[2]] + host + spn[res[3]:]
}

func (m GSSAPIMechanism) encode(outgoing []byte) ([]byte, error) {
	if m.qop == QOP_TO_FLAG[AUTH] {
		return outgoing, nil
	} else {
		var conf_flag bool = false
		if m.qop == QOP_TO_FLAG[AUTH_CONF] {
			conf_flag = true
		}
		return m.context.wrap(deepCopy(outgoing), conf_flag)
	}
}

func (m GSSAPIMechanism) decode(incoming []byte) ([]byte, error) {
	if m.qop == QOP_TO_FLAG[AUTH] {
		return incoming, nil
	}
	return m.context.unwrap(deepCopy(incoming))
}

func deepCopy(original []byte) []byte {
	copied := make([]byte, len(original))
	for i, el := range original {
		copied[i] = el
	}
	return copied
}

func (m GSSAPIMechanism) dispose() {
	m.context.dispose()
}

func (m GSSAPIMechanism) getConfig() *MechanismConfig {
	return m.config
}

type GSSAPIContext struct {
	DebugLog       bool
	RunAsService   bool
	ServiceName    string
	ServiceAddress string

	gssapi.Options

	*gssapi.Lib `json:"-"`
	loadonce    sync.Once

	// Service credentials loaded from keytab
	credential     *gssapi.CredId
	token          []byte
	continueNeeded bool
	contextId      *gssapi.CtxId
	reqFlags       uint32
	availFlags     uint32
}

//
func newGSSAPIContext() *GSSAPIContext {
	var c = &GSSAPIContext{
		reqFlags: uint32(gssapi.GSS_C_INTEG_FLAG) + uint32(gssapi.GSS_C_MUTUAL_FLAG) + uint32(gssapi.GSS_C_SEQUENCE_FLAG) + uint32(gssapi.GSS_C_CONF_FLAG),
	}
	prefix := "gosasl-client"
	err := loadlib(c.DebugLog, prefix, c)
	if err != nil {
		log.Fatal(err)
	}

	j, _ := json.MarshalIndent(c, "", "  ")
	c.Debug(fmt.Sprintf("Config: %s", string(j)))
	return c
}

// InitClientContext initializes the context and gets the response(token)
// to send to the server
func initClientContext(c *GSSAPIContext, service string, inputToken []byte) error {
	c.ServiceName = service

	var _inputToken *gssapi.Buffer
	var err error
	if inputToken == nil {
		_inputToken = c.GSS_C_NO_BUFFER
	} else {
		_inputToken, err = c.MakeBufferBytes(inputToken)
		defer _inputToken.Release()
		if err != nil {
			return err
		}
	}

	preparedName := prepareServiceName(c)
	defer preparedName.Release()

	contextId, _, token, outputRetFlags, _, err := c.InitSecContext(
		nil,
		c.contextId,
		preparedName,
		c.GSS_MECH_KRB5,
		c.reqFlags,
		0,
		c.GSS_C_NO_CHANNEL_BINDINGS,
		_inputToken)
	defer token.Release()

	c.token = token.Bytes()
	c.contextId = contextId
	c.availFlags = outputRetFlags
	return nil
}

// Wrap calls GSS_Wrap
func (c *GSSAPIContext) wrap(original []byte, conf_flag bool) (wrapped []byte, err error) {
	if original == nil {
		return
	}
	_original, err := c.MakeBufferBytes(original)
	defer _original.Release()

	if err != nil {
		return nil, err
	}
	_, wrappedBuffer, err := c.contextId.Wrap(conf_flag, gssapi.GSS_C_QOP_DEFAULT, _original)
	defer wrappedBuffer.Release()
	if err != nil {
		return nil, err
	}
	return wrappedBuffer.Bytes(), nil
}

// Unwrap calls GSS_Unwrap
func (c *GSSAPIContext) unwrap(original []byte) (unwrapped []byte, err error) {
	if original == nil {
		return
	}
	_original, err := c.MakeBufferBytes(original)
	defer _original.Release()

	if err != nil {
		return nil, err
	}
	unwrappedBuffer, _, _, err := c.contextId.Unwrap(_original)
	defer unwrappedBuffer.Release()
	if err != nil {
		return nil, err
	}
	return unwrappedBuffer.Bytes(), nil
}

// Dispose releases the acquired memory and destroys sensitive information
func (c *GSSAPIContext) dispose() error {
	if c.contextId != nil {
		return c.contextId.Unload()
	}
	return nil
}

// IntegAvail returns true in the integ_flag is available and therefore a security layer can be established
func (c *GSSAPIContext) integAvail() bool {
	return c.availFlags&uint32(gssapi.GSS_C_INTEG_FLAG) != 0
}

// ConfAvail returns true in the conf_flag is available and therefore a confidentiality layer can be established
func (c *GSSAPIContext) confAvail() bool {
	return c.availFlags&uint32(gssapi.GSS_C_CONF_FLAG) != 0
}

func loadlib(debug bool, prefix string, c *GSSAPIContext) error {
	max := gssapi.Err + 1
	if debug {
		max = gssapi.MaxSeverity
	}
	pp := make([]gssapi.Printer, 0, max)
	for i := gssapi.Severity(0); i < max; i++ {
		p := log.New(os.Stderr,
			fmt.Sprintf("%s: %s\t", prefix, i),
			log.LstdFlags)
		pp = append(pp, p)
	}
	c.Options.Printers = pp

	lib, err := gssapi.Load(&c.Options)
	if err != nil {
		return err
	}
	c.Lib = lib
	return nil
}

func prepareServiceName(c *GSSAPIContext) *gssapi.Name {
	if c.ServiceName == "" {
		log.Fatal("Need a --service-name")
	}

	nameBuf, err := c.MakeBufferString(c.ServiceName)
	defer nameBuf.Release()
	if err != nil {
		log.Fatal(err)
	}

	name, err := nameBuf.Name(c.GSS_KRB5_NT_PRINCIPAL_NAME)
	if err != nil {
		log.Fatal(err)
	}
	if name.String() != c.ServiceName {
		log.Fatalf("name: got %q, expected %q", name.String(), c.ServiceName)
	}

	return name
}
