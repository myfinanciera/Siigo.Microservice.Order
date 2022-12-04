package consumer

import (
	slim "dev.azure.com/SiigoDevOps/Siigo/_git/go-slim.git/abstractions"
	log "github.com/sirupsen/logrus"
)

// TestConsumer Consumer Order.Created
func TestConsumer(message slim.Message) {
	log.Info("[Consumer] Message on ", message.Topic, " ", string(message.Value))
	// Commit message
	message.CommitMessage(message)
}
