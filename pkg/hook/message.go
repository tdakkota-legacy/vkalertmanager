package hook

import (
	"time"
)

// Message defines the JSON object send to webhook endpoints.
type Message struct {
	Data

	// The protocol version.
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts uint64 `json:"truncatedAlerts"`
}

// Data is the data passed to notification templates and webhook pushes.
//
// End-users should not be exposed to Go's type system, as this will confuse them and prevent
// simple things like simple equality checks to fail. Map everything to float64/string.
type Data struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   Alerts `json:"alerts"`

	GroupLabels       KV `json:"groupLabels"`
	CommonLabels      KV `json:"commonLabels"`
	CommonAnnotations KV `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

// KV is a set of key/value string pairs.
type KV map[string]string

// Alerts is a list of Alert objects.
type Alerts []Alert

// Alert holds one alert for notification templates.
type Alert struct {
	Status       string    `json:"status"`
	Labels       KV        `json:"labels"`
	Annotations  KV        `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Fingerprint  string    `json:"fingerprint"`
}
