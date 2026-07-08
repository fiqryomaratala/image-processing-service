package queue

import "fmt"

func Health() error {
	conn := GetConnection()
	ch := GetChannel()

	if conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}

	if ch.IsClosed() {
		return fmt.Errorf("rabbitmq channel is closed")
	}

	return nil
}
