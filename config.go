/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package producer

import (
	"github.com/yanzay/cfg"
	"strconv"
	"strings"
	"time"
)

type ProducerConfig struct {
	Partitioner       Partitioner
	MetadataExpire    time.Duration
	CompressionType   string
	BatchSize         int
	Linger            time.Duration
	Retries           int
	RetryBackoff      time.Duration
	BlockOnBufferFull bool

	ClientID        string
	MaxRequests     int
	SendRoutines    int
	ReceiveRoutines int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	RequiredAcks    int
	AckTimeoutMs    int32
	BrokerList      []string
}

func NewProducerConfig() *ProducerConfig {
	return &ProducerConfig{
		Partitioner:     NewHashPartitioner(),
		MetadataExpire:  time.Minute,
		BatchSize:       1000,
		ClientID:        "siesta",
		MaxRequests:     10,
		SendRoutines:    10,
		ReceiveRoutines: 10,
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		RequiredAcks:    1,
		AckTimeoutMs:    1000,
		Linger:          1 * time.Second,
	}
}

func ProducerConfigFromFile(filename string) (*ProducerConfig, error) {
	c, err := cfg.LoadNewMap(filename)
	if err != nil {
		return nil, err
	}

	producerConfig := NewProducerConfig()
	if err := setDurationConfig(&producerConfig.MetadataExpire, c["metadata.max.age"]); err != nil {
		return nil, err
	}
	if err := setIntConfig(&producerConfig.BatchSize, c["batch.size"]); err != nil {
		return nil, err
	}
	if err := setIntConfig(&producerConfig.RequiredAcks, c["acks"]); err != nil {
		return nil, err
	}
	if err := setInt32Config(&producerConfig.AckTimeoutMs, c["timeout.ms"]); err != nil {
		return nil, err
	}
	if err := setDurationConfig(&producerConfig.Linger, c["linger"]); err != nil {
		return nil, err
	}
	setStringConfig(&producerConfig.ClientID, c["client.id"])
	if err := setIntConfig(&producerConfig.SendRoutines, c["send.routines"]); err != nil {
		return nil, err
	}
	if err := setIntConfig(&producerConfig.ReceiveRoutines, c["receive.routines"]); err != nil {
		return nil, err
	}
	setBoolConfig(&producerConfig.BlockOnBufferFull, c["block.on.buffer.full"])
	if err := setIntConfig(&producerConfig.Retries, c["retries"]); err != nil {
		return nil, err
	}
	if err := setDurationConfig(&producerConfig.RetryBackoff, c["retry.backoff"]); err != nil {
		return nil, err
	}
	setStringConfig(&producerConfig.CompressionType, c["compression.type"])
	if err := setIntConfig(&producerConfig.MaxRequests, c["max.requests"]); err != nil {
		return nil, err
	}

	setStringsConfig(&producerConfig.BrokerList, c["bootstrap.servers"])
	if len(producerConfig.BrokerList) == 0 {
		setStringsConfig(&producerConfig.BrokerList, c["metadata.broker.list"])
	}

	return producerConfig, nil
}

func setStringConfig(where *string, what string) {
	if what != "" {
		*where = what
	}
}

func setBoolConfig(where *bool, what string) {
	if what != "" {
		*where = what == "true"
	}
}

func setDurationConfig(where *time.Duration, what string) error {
	if what != "" {
		value, err := time.ParseDuration(what)
		if err == nil {
			*where = value
		}
		return err
	}
	return nil
}

func setIntConfig(where *int, what string) error {
	if what != "" {
		value, err := strconv.Atoi(what)
		if err == nil {
			*where = value
		}
		return err
	}
	return nil
}

func setInt32Config(where *int32, what string) error {
	if what != "" {
		value, err := strconv.Atoi(what)
		if err == nil {
			*where = int32(value)
		}
		return err
	}
	return nil
}

func setStringsConfig(where *[]string, what string) {
	*where = strings.Split(what, ",")
}
