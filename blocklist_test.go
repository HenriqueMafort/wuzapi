package main

import (
	"testing"

	"go.mau.fi/whatsmeow/types"
)

// TestFormatBlocklist covers the response shaping for GET /user/blocklist:
// a nil blocklist yields an empty list (never null) and an empty dhash, and a
// populated blocklist stringifies every JID and passes the dhash through.
func TestFormatBlocklist(t *testing.T) {
	t.Run("nil yields empty list and dhash", func(t *testing.T) {
		got := formatBlocklist(nil)
		jids, ok := got["Blocklist"].([]string)
		if !ok {
			t.Fatalf("Blocklist is not []string: %T", got["Blocklist"])
		}
		if len(jids) != 0 {
			t.Errorf("len(Blocklist) = %d; want 0", len(jids))
		}
		if got["DHash"] != "" {
			t.Errorf("DHash = %v; want empty", got["DHash"])
		}
	})

	t.Run("stringifies JIDs and passes dhash through", func(t *testing.T) {
		j1, ok1 := parseJID("5491155554445")
		j2, ok2 := parseJID("5491155553935")
		if !ok1 || !ok2 {
			t.Fatalf("parseJID failed: ok1=%v ok2=%v", ok1, ok2)
		}
		bl := &types.Blocklist{DHash: "1234567890", JIDs: []types.JID{j1, j2}}

		got := formatBlocklist(bl)
		jids := got["Blocklist"].([]string)
		if len(jids) != len(bl.JIDs) {
			t.Fatalf("len(Blocklist) = %d; want %d", len(jids), len(bl.JIDs))
		}
		for i := range bl.JIDs {
			if jids[i] != bl.JIDs[i].String() {
				t.Errorf("Blocklist[%d] = %q; want %q", i, jids[i], bl.JIDs[i].String())
			}
		}
		if got["DHash"] != "1234567890" {
			t.Errorf("DHash = %v; want %q", got["DHash"], "1234567890")
		}
	})
}
