package kafkautil

import (
	"crm-util-go/common"
	"crm-util-go/pointer"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"strconv"
	"strings"
	"sync"
)

func (kc *KafkaConfig) StartConsumer() {
	kc.IsShutdown = false
}

func (kc *KafkaConfig) ShutdownConsumer() {
	kc.IsShutdown = true
}

func (kc *KafkaConfig) postLineNotify(transID string, message string) {
	kc.Logger.Debug(transID, message)
	/*
		lineToken := config.GetString("lineToken")

		httpClient := httpclient.NewHttpClient()
		httpClient.Logger = kc.Logger
		_, errLine := httpClient.PostLineNotify(transID, lineToken, message)

		if errLine != nil {
			kc.Logger.Error(transID, "Error PostLineNotify: "+errLine.Error())
		}
	*/
}

func getTopicInfo(tp kafka.TopicPartition) string {
	topicName := pointer.GetStringValue(tp.Topic)
	partition := strconv.FormatInt(int64(tp.Partition), 10)
	offset := tp.Offset.String()

	return fmt.Sprintf("Topic: %s, Partition: %s, Offset: %s", topicName, partition, offset)
}

func (kc *KafkaConfig) logKafkaHeaders(msgTransID string, kafkaHeaders []kafka.Header) {
	if len(kafkaHeaders) > 0 {
		var headersBuilder strings.Builder
		headersBuilder.WriteString("Kafka Header.")

		for _, header := range kafkaHeaders {
			headersBuilder.WriteString("\nKey: ")
			headersBuilder.WriteString(header.Key)

			if header.Value == nil {
				headersBuilder.WriteString(", Value: nil")
			} else {
				headerValue := string(header.Value)
				headersBuilder.WriteString(", Value: ")
				headersBuilder.WriteString(headerValue)
			}

			headersBuilder.WriteString("\n")
		}

		kc.Logger.Info(msgTransID, headersBuilder.String())
	}
}

func (kc *KafkaConfig) Consumer(wg *sync.WaitGroup, transID string,
	onMessage func(transID string, topicName string, kafkaKey string, kafkaMsg string) error) {

	defer wg.Done()
	defer kc.ShutdownConsumer()

	kc.Logger.Info(transID, fmt.Sprintf("Starting Consumer with KafkaConfig: %#v", kc))

	// "enable.ssl.certificate.verification": false,
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":                     kc.BootstrapServers,
		"group.id":                              kc.GroupID,
		"auto.offset.reset":                     "earliest",
		"enable.auto.commit":                    false,
		"enable.partition.eof":                  true,
		"security.protocol":                     kc.SecurityProtocol,
		"sasl.mechanisms":                       kc.SaslMechanism,
		"sasl.kerberos.principal":               kc.KerberosPrincipalName,
		"sasl.kerberos.service.name":            kc.KerberosServiceName,
		"sasl.kerberos.keytab":                  kc.KerberosKeytab,
		"sasl.kerberos.min.time.before.relogin": kc.KerberosReloginMS,
		"ssl.ca.location":                       kc.SslCALocation,
		"ssl.cipher.suites":                     kc.SslCipherSuites,
	})

	if err != nil {
		kc.Logger.Error(transID, fmt.Sprintf("Consumer.NewConsumer Error: %v", err))
		panic(err)
	}

	err = consumer.SubscribeTopics(kc.TopicName, nil)

	if err != nil {
		kc.Logger.Error(transID, fmt.Sprintf("Consumer.SubscribeTopics Error: %v", err))
		panic(err)
	}

	defer consumer.Close()

	// Synchronous commits
	var isRun bool = true
	var pollTimeoutMs int = 10000

	for isRun {
		if kc.IsShutdown {
			isRun = false
			break
		}

		event := consumer.Poll(pollTimeoutMs)

		switch e := event.(type) {
		case *kafka.Message:
			msgTransID := common.NewUUID()
			kafkaKey := string(e.Key)
			kafkaMsg := string(e.Value)
			kafkaHeaders := e.Headers

			tp := e.TopicPartition
			topicInfo := getTopicInfo(tp)

			kc.Logger.Info(msgTransID, fmt.Sprintf("Start message on %s, Key: %s, Value: %s",
				topicInfo, kafkaKey, kafkaMsg))

			kc.logKafkaHeaders(msgTransID, kafkaHeaders)

			err = onMessage(msgTransID, pointer.GetStringValue(tp.Topic), kafkaKey, kafkaMsg)

			if err == nil {
				_, err = consumer.CommitMessage(e)

				if err != nil {
					kc.IsShutdown = true
					kc.Logger.Error(msgTransID, topicInfo+" Consumer.CommitMessage Error: "+err.Error())
				} else {
					kc.Logger.Info(msgTransID, topicInfo+" Consumer.CommitMessage Success")
				}
			} else {
				notifyMessage := "Error Kafka Consumer. App: " + kc.Logger.ApplicationName + ", " +
					topicInfo + ", Error: " + err.Error()

				kc.Logger.Error(msgTransID, topicInfo+" Consumer Rollback Message")

				_, err = consumer.SeekPartitions([]kafka.TopicPartition{e.TopicPartition})

				if err != nil {
					kc.IsShutdown = true
					kc.Logger.Error(msgTransID, topicInfo+" Consumer.SeekPartitions Error: "+err.Error())
				}

				kc.postLineNotify(msgTransID, notifyMessage)
			}
		case kafka.PartitionEOF:
			kc.Logger.Info(transID, fmt.Sprintf("Reached the end of a partition. %v", e))
		case kafka.Error:
			kc.IsShutdown = true
			kc.Logger.Error(transID, fmt.Sprintf("Consumer.Poll Error: %v", e))
		default:
			kc.Logger.Info(transID, fmt.Sprintf("Ignored %v", e))
		}
	}
}
