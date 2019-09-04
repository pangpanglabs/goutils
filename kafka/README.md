# goutils/kafka

Wrapper of [sarama](https://github.com/Shopify/sarama)

## Getting Started

### Producer

```golang
producer, err := kafka.NewProducer(brokers, topic, func(c *sarama.Config) {
        c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
        c.Producer.Compression = sarama.CompressionGZIP     // Compress messages
        c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

})
if err != nil {
        return err
}

msg := map[string]interface{}{
        "time": time.Now(),
        "idx":  i,
        "msg":  rand.Int(),
}
if err := producer.Send(msg); err != nil {
        return err
}
```

## Consumer

```golang
consumer, err := kafka.NewConsumer(brokers, topic, kafka.AllPartitions, sarama.OffsetNewest)
if err != nil {
        return err
}

messages, err := consumer.Messages()
if err != nil {
        return err
}

for m := range messages {
        var v interface{}
        d := json.NewDecoder(bytes.NewReader(m.Value))
        d.UseNumber()

        if err := d.Decode(&v); err != nil{
                return err
        }

        fmt.Printf("[Receive] Offset:%d\tPartition:%d\tValue:%v\n", m.Offset, m.Partition, v)
}
```

## Consumer Group

```golang
consumer, err := kafka.NewConsumerGroup(groupId, brokers, topic)
if err != nil {
        return err
}

messages, err := consumer.Messages()
if err != nil {
        return err
}

for m := range messages {
        var v interface{}
        d := json.NewDecoder(bytes.NewReader(m.Value))
        d.UseNumber()

        if err := d.Decode(&v); err != nil{
                return err
        }

        fmt.Printf("[Receive] Offset:%d\tPartition:%d\tValue:%v\n", m.Offset, m.Partition, v)
}
```