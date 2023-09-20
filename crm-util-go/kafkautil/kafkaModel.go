package kafkautil

import (
	"crm-util-go/logging"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

/*
Default Kerberos on Linux:
/etc/krb5.conf
/etc/krb5.keytab

Reference:
https://github.com/confluentinc/confluent-kafka-go/tree/v2.1.1
https://github.com/confluentinc/librdkafka/tree/v2.1.1
https://github.com/edenhill/librdkafka/wiki/Using-SSL-with-librdkafka
https://github.com/edenhill/librdkafka/wiki/Using-SASL-with-librdkafka
https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md

SecurityProtocol: PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL

Example: TrueIT KafkaProducerConfig
BootstrapServers: "kkts01:9094,kkts02:9094,kkts03:9094"
SecurityProtocol : "SASL_SSL"
SaslMechanism: "GSSAPI"
KerberosPrincipalName: "kfuat_crmuser01@KAFKA.SECURE"
KerberosServiceName: "bigfoot"
KerberosKeytab: "/etc/krb5.keytab"
kerberosReloginMS: 180000
SslCALocation:  "/etc/ssl/certs/ca-certificates.crt"
SslCipherSuites: "DHE-DSS-AES256-GCM-SHA384"
CompressionType: "lz4"
*/

type KafkaConfig struct {
	BootstrapServers      string
	TopicName             []string
	SecurityProtocol      string
	SaslMechanism         string
	KerberosPrincipalName string
	KerberosServiceName   string
	KerberosKeytab        string
	KerberosReloginMS     int
	SslCALocation         string
	SslCipherSuites       string
	CompressionType       string
	GroupID               string
	IsShutdown            bool
	IsLineNotify          bool
	Logger                *logging.PatternLogger
}

type KafkaMessage struct {
	Key     string
	Value   string
	Headers []kafka.Header
}
