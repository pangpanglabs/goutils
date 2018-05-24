package httpreq

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
)

type formatType int

const (
	JsonType      formatType = iota //0
	FormType                        //1
	XmlType                         //2
	ByteArrayType                   //3
)

const (
	charsetUTF8 = "charset=UTF-8"
)
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = "application/json" + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript" + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = "application/xml" + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = "text/xml" + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationFormUTF8              = "application/x-www-form-urlencoded" + "; " + charsetUTF8
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = "text/html" + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = "text/plain" + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

type reqFormatter interface {
	marshal(data interface{}) ([]byte, error)
	contentType() string
	unMarshal(data []byte, v interface{}) error
}

type DataTypeFactory struct {
}

func (DataTypeFactory) New(dataType formatType) reqFormatter {
	switch dataType {
	case XmlType:
		return XmlFormat{}
	case FormType:
		return FormFormat{}
	case JsonType:
		return JsonFormat{}
	case ByteArrayType:
		return ByteArrayFormat{}
	default:
		panic("dataType only supports xml,form,json")
	}
}

type XmlFormat struct {
}

func (XmlFormat) marshal(param interface{}) ([]byte, error) {
	paramStr, ok := param.(string)
	if !ok {
		return nil, errors.New("param is expected to string. ")
	}
	byteData := []byte(paramStr)
	return byteData, nil
}
func (XmlFormat) contentType() string {
	return "application/xml"
}
func (XmlFormat) unMarshal(data []byte, v interface{}) error {
	if err := xml.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

type JsonFormat struct {
}

func (JsonFormat) marshal(param interface{}) ([]byte, error) {
	return json.Marshal(param)
}
func (JsonFormat) contentType() string {
	return "application/json"
}
func (JsonFormat) unMarshal(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

type FormFormat struct {
}

func (FormFormat) marshal(param interface{}) ([]byte, error) {
	paramStr, ok := param.(string)
	if !ok {
		return nil, errors.New("param is expected to string. ")
	}
	byteData := []byte(paramStr)
	return byteData, nil
}
func (FormFormat) contentType() string {
	return "application/x-www-form-urlencoded" + "; " + "charset=UTF-8"
}
func (FormFormat) unMarshal(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

type ByteArrayFormat struct {
}

func (ByteArrayFormat) marshal(param interface{}) ([]byte, error) {
	paramData, ok := param.([]byte)
	if !ok {
		return nil, errors.New("param is expected to []byte. ")
	}
	return paramData, nil
}
func (ByteArrayFormat) contentType() string {
	return "application/x-www-form-urlencoded" + "; " + "charset=UTF-8"
}
func (ByteArrayFormat) unMarshal(data []byte, v interface{}) error {
	rawData, ok := v.(*[]byte)
	if !ok {
		return errors.New("param is expected to *[]byte. ")
	}
	*rawData = append(*rawData, data...)
	return nil
}

func CertTransport(certFile string, keyFile string, caFile string) (*http.Transport, error) {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	var caCertPool *(x509.CertPool)
	if caFile != "" {
		// Load CA cert
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}
		caCertPool = x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	return &http.Transport{TLSClientConfig: tlsConfig}, nil
}
