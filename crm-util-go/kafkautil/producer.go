package kafkautil

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	kafkaAcksSyncReplicas   int = -1
	kafkaFlushTimeoutMs     int = 10000
	kafkaFlushBulkTimeoutMs int = 10000
	kafkaRequestTimeoutMs   int = 30000
	kafkaLingerMs           int = 0
	kafkaImmediatePublish   int = 0
)

func (kc KafkaConfig) Publish(transID string, key string, value string, headers []kafka.Header) (*kafka.Message, error) {
	kc.Logger.Info(transID, "Starting Publish key: "+key+", value: "+value)

	var kafkaMessage *kafka.Message

	// "compression.type":           kc.CompressionType,
	// "enable.ssl.certificate.verification": "false",
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":                     kc.BootstrapServers,
		"request.required.acks":                 kafkaAcksSyncReplicas,
		"request.timeout.ms":                    kafkaRequestTimeoutMs,
		"linger.ms":                             kafkaLingerMs,
		"queue.buffering.max.ms":                kafkaImmediatePublish,
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
		kc.Logger.Error(transID, "Producer can not connect to kafka server", err)
		return kafkaMessage, err
	}

	defer producer.Close()

	// Optional delivery channel, if not specified the Producer object's
	// .Events channel is used.
	deliveryChan := make(chan kafka.Event)
	// closes a channel
	defer close(deliveryChan)

	// Produce messages to topic (asynchronously)
	// []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}}
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kc.TopicName[0], Partition: kafka.PartitionAny},
		Value:          []byte(value),
		Key:            []byte(key),
		Timestamp:      time.Now(),
		Headers:        headers,
	}, deliveryChan)

	if err != nil {
		kc.Logger.Error(transID, "Producer.Produce() failed because "+err.Error())
		return kafkaMessage, err
	}

	kafkaEvent := <-deliveryChan
	kafkaMessage = kafkaEvent.(*kafka.Message)

	if kafkaMessage.TopicPartition.Error != nil {
		errMsg := "Delivery message failed because " + kafkaMessage.TopicPartition.Error.Error()
		kc.Logger.Error(transID, errMsg)
		return kafkaMessage, errors.New(errMsg)
	}

	kc.Logger.Info(transID, fmt.Sprintf("Delivered message success to Topic: %s Partition: %d Offset: %v",
		*kafkaMessage.TopicPartition.Topic, kafkaMessage.TopicPartition.Partition, kafkaMessage.TopicPartition.Offset))

	// Wait for message deliveries before shutting down
	producer.Flush(kafkaFlushTimeoutMs)

	return kafkaMessage, err
}

func (kc KafkaConfig) PublishBulk(transID string, kafkaMessageList []KafkaMessage, maxThread int,
	onProducedMessage func(wg *sync.WaitGroup, transID string, key string, value string, err error)) error {
	kc.Logger.Info(transID, "Starting PublishBulk size:", len(kafkaMessageList))

	// "compression.type":           kc.CompressionType,
	// "enable.ssl.certificate.verification": "false",
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":                     kc.BootstrapServers,
		"request.required.acks":                 kafkaAcksSyncReplicas,
		"request.timeout.ms":                    kafkaRequestTimeoutMs,
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
		kc.Logger.Error(transID, "Producer can not connect to kafka server", err)
		return err
	}

	defer producer.Close()

	deliveryChan := make(chan kafka.Event)
	// closes a channel
	defer close(deliveryChan)

	// Produce messages
	totalProduceMessage := 0
	topicPartition := kafka.TopicPartition{Topic: &kc.TopicName[0], Partition: kafka.PartitionAny}
	for _, kafkaMessage := range kafkaMessageList {
		// Produce single message to topic (asynchronously)
		errProduce := producer.Produce(&kafka.Message{
			TopicPartition: topicPartition,
			Value:          []byte(kafkaMessage.Value),
			Key:            []byte(kafkaMessage.Key),
			Timestamp:      time.Now(),
			Headers:        kafkaMessage.Headers,
		}, deliveryChan)

		if errProduce != nil {
			kc.Logger.Error(transID, "producer.Produce() failed because "+errProduce.Error())
			break
		} else {
			totalProduceMessage++
		}
	}
	kc.Logger.Info(transID, "Finished producer.Produce()")

	// Delivery report handler for produced messages
	// NOTE: Message headers are not available on producer delivery report messages.
	countThread := 0
	var wg sync.WaitGroup
	for i := 0; i < totalProduceMessage; i++ {
		kafkaEvent := <-deliveryChan

		switch kafkaEvent.(type) {
		case *kafka.Message:
			countThread++
			kafkaMessage := kafkaEvent.(*kafka.Message)
			if kafkaMessage.TopicPartition.Error != nil {
				kc.Logger.Error(transID, "Delivery message failed because "+kafkaMessage.TopicPartition.Error.Error())

				wg.Add(1)
				go onProducedMessage(&wg, transID, string(kafkaMessage.Key), string(kafkaMessage.Value), kafkaMessage.TopicPartition.Error)
			} else {
				kc.Logger.Info(transID, fmt.Sprintf("Delivered message success to Topic: %s Partition: %d Offset: %v",
					*kafkaMessage.TopicPartition.Topic, kafkaMessage.TopicPartition.Partition, kafkaMessage.TopicPartition.Offset))

				wg.Add(1)
				go onProducedMessage(&wg, transID, string(kafkaMessage.Key), string(kafkaMessage.Value), nil)
			}

			if countThread == maxThread {
				kc.Logger.Info(transID, "Waiting thread in loop ...")
				wg.Wait()
				countThread = 0
			}
		case kafka.Error:
			kc.Logger.Error(transID, fmt.Sprintf("Error: %v", kafkaEvent))
		default:
			kc.Logger.Info(transID, fmt.Sprintf("Ignored event: %v", kafkaEvent))
		}
	}

	if countThread > 0 {
		kc.Logger.Info(transID, "Waiting thread ...")
		wg.Wait()
	}

	producer.Flush(kafkaFlushBulkTimeoutMs)

	return err
}
