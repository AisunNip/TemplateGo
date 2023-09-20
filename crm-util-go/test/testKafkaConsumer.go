package test

/*
import (
	"crm-util-go/common"
	"crm-util-go/kafkautil"
	"crm-util-go/logging"
	"sync"
)

var schLogger = logging.InitScheduleLogger("crm-util-go", logging.CrmDatabase)

func onMessage(mainTransID string, kafkaKey string, kafkaMsg string) error {
	var err error

	transID, _ := common.NewUUID()
	schLogger.Info(transID, "Start onMessage. MainTransID: " + mainTransID)

	// TODO: Logic

	return err
}

func ConsumerController() {
	transID, _ := common.NewUUID()

	schLogger.Info(transID, "Start ConsumerController")

	kc := new(kafkautil.KafkaConfig)
	kc.BootstrapServers = "kkts01:9094,kkts02:9094,kkts03:9094"
	kc.TopicName = "TestTopicName"
	kc.SecurityProtocol = "SASL_SSL"
	kc.SaslMechanism = "GSSAPI"
	kc.KerberosPrincipalName = "kfuat_crmuser01@KAFKA.SECURE"
	kc.KerberosServiceName = "bigfoot"
	kc.KerberosKeytab = "/etc/kfuat_crmuser01.client.keytab"
	kc.SslCALocation = "/etc/server.crt"
	kc.SslCipherSuites = "DHE-DSS-AES256-GCM-SHA384"
	kc.Logger = schLogger

	schLogger.Info(transID, "Initial KafkaConfig success")

	var noOfConsumer int = 10
	var wg sync.WaitGroup

	kc.StartConsumer()

	for i := 0; i < noOfConsumer; i++ {
		wg.Add(1)
		go kc.Consumer(&wg, transID, onMessage)
	}

	wg.Wait()
	schLogger.Info(transID, "End ConsumerController")
}
*/