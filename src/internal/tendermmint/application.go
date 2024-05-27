package tendermmint

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	abciserver "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
)

type Application struct {
	types.BaseApplication
	storage sync.Map
}

func (app *Application) StartABCI(address string) {
	server := abciserver.NewSocketServer(address, app)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
	defer server.Stop()

	log.Printf("ABCI server started at %s\n", address)
}

func (app *Application) DeliverTx(tx types.RequestDeliverTx) types.ResponseDeliverTx {
	var payload map[string]string
	if err := json.Unmarshal(tx.GetTx(), &payload); err != nil {
		return types.ResponseDeliverTx{Code: 1, Log: fmt.Sprintf("Invalid JSON: %v", err)}
	}

	if key, ok := payload["key"]; ok {
		if value, ok := payload["value"]; ok {
			app.storage.Store(key, value)
			return types.ResponseDeliverTx{Code: 0, Log: "Success"}
		}
	}
	return types.ResponseDeliverTx{Code: 1, Log: "Invalid Tx"}
}

func (app *Application) Query(req types.RequestQuery) types.ResponseQuery {
	if value, ok := app.storage.Load(string(req.Data)); ok {
		return types.ResponseQuery{Code: 0, Value: []byte(value.(string))}
	}
	return types.ResponseQuery{Code: 1, Log: "Not found"}
}
