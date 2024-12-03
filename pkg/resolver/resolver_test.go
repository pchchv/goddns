package resolver

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	servers := []string{"8.8.8.8", "8.8.4.4"}
	expectedServers := []string{"8.8.8.8:53", "8.8.4.4:53"}
	resolver := New(servers)
	if !reflect.DeepEqual(resolver.Servers, expectedServers) {
		t.Error("resolver.Servers: ", resolver.Servers, "should be equal to", expectedServers)
	}
}
