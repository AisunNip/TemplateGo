package sseapi

import (
	"encoding/json"
	"github.com/r3labs/sse/v2"
	"net/http"
)

type SseNotifyMsg struct {
	StreamID string `json:"streamID,omitempty"`
	Data     string `json:"data,omitempty"`
}

type SseResponse struct {
	Code string `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func clientDisconnected(r *http.Request) {
	// Received Browser Disconnection
	<-r.Context().Done()
	println("The client is disconnected. stream=" + r.URL.Query().Get("stream"))
	return
}

func StartSSEApi(addr string, streamIDList []string) {
	server := sse.New()
	server.AutoReplay = false

	if len(streamIDList) > 0 {
		for _, streamID := range streamIDList {
			server.CreateStream(streamID)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Auto create a stream
		streamID := r.URL.Query().Get("stream")

		if streamID != "" {
			if !server.StreamExists(streamID) {
				println("Create a new stream: " + streamID)
				server.CreateStream(streamID)
			}
		}

		go clientDisconnected(r)

		server.ServeHTTP(w, r)
	})

	mux.HandleFunc("/notifyMsg", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var msg SseNotifyMsg
		err := decoder.Decode(&msg)
		var sseResp SseResponse

		if err != nil {
			http.Error(w, "Decode http body error!", http.StatusInternalServerError)
			println("Decode http body error " + err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if !server.StreamExists(msg.StreamID) {
			sseResp.Code = "1"
			sseResp.Msg = "StreamID " + msg.StreamID + " not found"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(sseResp)
			return
		}

		server.Publish(msg.StreamID, &sse.Event{
			Data: []byte(msg.Data),
		})

		sseResp.Code = "0"
		sseResp.Msg = "Success"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sseResp)
	})

	http.ListenAndServe(addr, mux)
}
