package lambdahndlr

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
)

type Event struct {
	ID    float64 `json:"id"`
	Value string  `json:"value"`
}

type Response struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Handler(event Event) (Response, error) {
	log.Debug("Got event: %s", event.Value)

	// TODO: Call to another function for actual logic

	return Response{
		Message: fmt.Sprintf("Processed request ID %f", event.ID),
		Ok:      true,
	}, nil
}
