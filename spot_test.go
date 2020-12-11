package main

import (
	"fmt"
	"testing"
)

type testCase struct {
	name string
	s    Spot
	want string
}

func TestClusterFormat(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{name: "Long Call sign", s: Spot{SpotID: 123, Activator: "SX2500S", Frequency: "3640", SpotTime: "2020-11-28T20:09:11", Spotter: "SV2DSJ-1",
			Comments: "", Source: "", Invalid: nil, Name: "Something"},
			want: "DX de SV2DSJ-1:   3640.0  SX2500S                                     2009Z"},
		//			{want: "DX de LU4DPL:    28443.0  CE3TVR       5/7 GF05QL                     2007Z"},
		{name: "Long frequency with comment", s: Spot{SpotID: 123, Activator: "PU2TTN", Frequency: "2400000", SpotTime: "2020-11-28T20:07:11", Spotter: "PU2RND",
			Comments: "CQ/CQ 10M", Source: "", Invalid: nil, Name: "Something"},
			want: "DX de PU2RND:  2400000.0  PU2TTN       CQ/CQ 10M                      2007Z"},
		{name: "No comments", s: Spot{SpotID: 123, Activator: "HI8JSG", Frequency: "28463", SpotTime: "2020-11-28T20:07:11", Spotter: "PR2E",
			Comments: "", Source: "", Invalid: nil, Name: "Something"},
			want: "DX de PR2E:      28463.0  HI8JSG                                      2007Z"},
	}
	for _, tc := range testCases {
		got := tc.s.ToClusterFormat()
		if fmt.Sprintf("%s\a\a\x0c", tc.want) != got {
			t.Errorf("\nwant: %s|||| (Test: %s)\ngot : %s||||", tc.want, tc.name, got)
		}
	}
}
