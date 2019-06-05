package variableserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type VariableSocket struct {
	data map[string]interface{}
	mux  sync.Mutex
}

func New() VariableSocket {
	return VariableSocket{data: make(map[string]interface{})}
}

func (vs *VariableSocket) Set(name string, v interface{}) {
	vs.mux.Lock()
	defer vs.mux.Unlock()

	vs.data[name] = v
}

func (vs *VariableSocket) Delete(name string) {
	vs.mux.Lock()
	defer vs.mux.Unlock()

	delete(vs.data, name)
}

func (vs *VariableSocket) Start(addr, path string) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, vs)
	})
	fmt.Println(http.ListenAndServe(addr, nil))
}

func checkOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{CheckOrigin: checkOrigin}

func fn(w http.ResponseWriter, r *http.Request, vs *VariableSocket) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		vs.mux.Lock()
		b, err := json.Marshal(vs.data)
		if err != nil {
			return
		}
		vs.mux.Unlock()

		err = c.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			return
		}
	}

}
