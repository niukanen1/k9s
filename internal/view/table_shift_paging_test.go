// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package view

import (
	"testing"

	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/model1"
	"github.com/derailed/tcell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Shift+Up/Down must page like PageUp/PageDown (vi-style), not move one row.
func TestTableShiftUpDownPaging(t *testing.T) {
	v := NewTable(client.NewGVR("test"))
	require.NoError(t, v.Init(makeContext(t)))

	evts := make([]model1.RowEvent, 0, 30)
	for i := range 30 {
		evts = append(evts, model1.RowEvent{
			Row: model1.Row{Fields: model1.Fields{"ns1", string(rune('a' + i)), "10", "3m"}},
		})
	}
	data := model1.NewTableDataWithRows(
		client.NewGVR("test"),
		model1.Header{
			model1.HeaderColumn{Name: "NAMESPACE"},
			model1.HeaderColumn{Name: "NAME"},
			model1.HeaderColumn{Name: "FRED"},
			model1.HeaderColumn{Name: "AGE", Attrs: model1.Attrs{Time: true}},
		},
		model1.NewRowEventsWithEvts(evts...),
	)
	cdata := v.Update(data, false)
	v.UpdateUI(cdata, data)
	require.Equal(t, 31, v.GetRowCount())

	uu := map[string]struct {
		key tcell.Key
		e   tcell.Key
	}{
		"shift_down": {key: tcell.KeyDown, e: tcell.KeyPgDn},
		"shift_up":   {key: tcell.KeyUp, e: tcell.KeyPgUp},
	}

	for k, u := range uu {
		t.Run(k, func(t *testing.T) {
			v.Select(1, 0)
			res := v.keyboard(tcell.NewEventKey(u.key, 0, tcell.ModShift))
			require.NotNil(t, res)
			assert.Equal(t, u.e, res.Key())
		})
	}
}
