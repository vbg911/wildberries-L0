package message

import (
	"Wildberries-L0/internal/common"
	"encoding/json"
	"fmt"
	"os"
)

func Generate(msgName string, id string) ([]byte, error) {
	msg, err := os.ReadFile(msgName)
	if err != nil {
		return nil, fmt.Errorf("error while reading file with msg: %w", err)
	}
	var order common.Order
	err = json.Unmarshal(msg, &order)
	if err != nil {
		return nil, fmt.Errorf("error while parsing msg from file: %w", err)
	}
	order.OrderUid = id
	order.Payment.Transaction = id
	orderInfo, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("error while marshal order: %w", err)
	}
	return orderInfo, nil
}
