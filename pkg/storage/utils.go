package Storage

import (
	"context"
	"strings"

	"github.com/atomix/atomix-go-client/pkg/atomix"
)

func getRawDataFromStore(urn string) ([]byte, error) {
	ctx := context.Background()

	// Create a slice of URN elements
	urnElems := strings.SplitN(urn, ".", 2)

	// Getting Map (store)
	store, err := atomix.GetMap(ctx, urnElems[0])
	if err != nil {
		log.Errorf("Error getting store \"%s\": %v", urnElems[0], err)
		return nil, err
	}

	// Request value from Map (store)
	data, err := store.Get(ctx, urnElems[1])
	if err != nil {
		log.Errorf("Error getting entry \"%s\" from %s: %v", urnElems[1], urnElems[1], err)
		return nil, err
	}

	return data.Value, nil
}
