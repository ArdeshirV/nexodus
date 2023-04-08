/*
Nexodus API

This is the Nexodus API Server.

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public

// ModelsEndpoint struct for ModelsEndpoint
type ModelsEndpoint struct {
	// IP address and port of the endpoint.
	Address string `json:"address,omitempty"`
	// Distance in milliseconds from the node to the ip address
	Distance int32 `json:"distance,omitempty"`
	// How the endpoint was discovered
	Source string `json:"source,omitempty"`
}
