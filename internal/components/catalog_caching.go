package components

var catalogCaching = Section{
	Name: "Caching & Messaging",
	Items: []InstallItem{
		{
			Name:    "Redis",
			Desc:    "In-memory key-value store and cache for performant applications",
			Formula: "redis",
			Color:   rgb(255, 80, 70),
			Binary:  "redis-cli",
		},
		{
			Name:    "Memcached",
			Desc:    "Distributed memory object caching for high-scale web workloads",
			Formula: "memcached",
			Color:   rgb(160, 200, 100),
			Binary:  "memcached",
		},
		{
			Name:    "NATS",
			Desc:    "Cloud-native messaging system for microservices and edge computing",
			Formula: "nats-server",
			Color:   rgb(150, 220, 255),
			Binary:  "nats",
		},
		{
			Name:    "Kafka",
			Desc:    "Distributed event streaming for high-throughput data pipelines",
			Formula: "kafka",
			Color:   rgb(200, 200, 200),
			Binary:  "kafka-topics",
		},
		{
			Name:    "RabbitMQ",
			Desc:    "Robust message broker supporting multiple messaging protocols",
			Formula: "rabbitmq",
			Color:   rgb(255, 140, 50),
			Binary:  "rabbitmqctl",
		},
		{
			Name:    "ActiveMQ",
			Desc:    "Multi-protocol broker for enterprise-grade message queuing",
			Formula: "activemq",
			Color:   rgb(255, 100, 150),
			Binary:  "activemq",
		},
		{
			Name:    "ZeroMQ",
			Desc:    "Asynchronous messaging library for distributed system development",
			Formula: "zeromq",
			Color:   rgb(150, 150, 150),
			Binary:  "zmq",
		},
		{
			Name:    "kcat",
			Desc:    "Generic tool for producer and consumer operations on Kafka",
			Formula: "kcat",
			Color:   rgb(60, 180, 255),
			Binary:  "kcat",
		},
	},
}
