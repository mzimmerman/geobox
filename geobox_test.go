/*
# A Go implementation of Brett Slatkin's mutiny - http://code.google.com/p/mutiny/source/browse/trunk/geobox.py
#
# Copyright 2012 Matt Zimmerman
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

package geobox

import (
	"sort"
	"strconv"
	"testing"
)

func TestLocations(t *testing.T) {
	tests := []interface{}{
		"37.78452", "-122.395320", 6, 10, "37.784530|-122.395330|37.784520|-122.395320",
		"37.78452", "-122.395320", 6, 25, "37.784525|-122.395325|37.784500|-122.395300",
		"37.78452", "-122.3953200", 7, 25, "37.7845225|-122.3953225|37.7845200|-122.3953200",
		"37.78452", "-122.3953200", 7, 1, "37.7845201|-122.3953201|37.7845200|-122.3953200",
		"37.78452", "-122.39532", 5, 15, "37.78455|-122.39535|37.78440|-122.39520",
		"37.78452", "-122.39532", 4, 17, "37.7859|-122.3966|37.7842|-122.3949",
		"37.78452", "-122.39531", 4, 17, "37.7859|-122.3966|37.7842|-122.3949",
		"37.78452", "-122.39667", 4, 17, "37.7859|-122.3983|37.7842|-122.3966",
		"37.78452", "-122.39532", 5, 2, "37.78452|-122.39532|37.78450|-122.39530",
		"37.78452", "-122.39532", 5, 10, "37.78460|-122.39540|37.78450|-122.39530",
		"37.78452", "-122.39532", 5, 25, "37.78475|-122.39550|37.78450|-122.39525",
		"38.96367", "-76.501164", 1, 5, "39.0|-77.0|38.5|-76.5",
	}
	for row := 0; row < len(tests); row += 5 {
		lat, _ := strconv.ParseFloat(tests[row].(string), 64)
		long, _ := strconv.ParseFloat(tests[row+1].(string), 64)
		resolution := tests[row+2].(int)
		slice := tests[row+3].(int)
		geocode := tests[row+4].(string)
		location := Compute(lat, long, resolution, slice)
		if location.Geocell != geocode {
			t.Errorf("#%d - For %f, %f - %d - %d", row/5, lat, long, resolution, slice)
			t.Errorf("Result - %s", location.Geocell)
			t.Errorf("Answer - %s", geocode)
		}
	}
}

type ComputeTest struct {
	Lat        float64
	Long       float64
	Slice      int
	Resolution int
	Geocells   []string
	Tier       int
}

func TestComputeSet(t *testing.T) {
	tests := []ComputeTest{
		ComputeTest{
			Lat:        37.78452,
			Long:       -122.395320,
			Slice:      5,
			Resolution: 1,
			Tier:       1,
			Geocells: []string{
				"37.5|-122.0|37.0|-121.5",
				"38.0|-122.0|37.5|-121.5",
				"38.5|-122.0|38.0|-121.5",
				"37.5|-122.5|37.0|-122.0",
				"38.0|-122.5|37.5|-122.0",
				"38.5|-122.5|38.0|-122.0",
				"37.5|-123.0|37.0|-122.5",
				"38.0|-123.0|37.5|-122.5",
				"38.5|-123.0|38.0|-122.5",
			},
		},
		ComputeTest{
			Lat:        37.78452,
			Long:       -122.395320,
			Slice:      5,
			Resolution: 1,
			Tier:       2,
			Geocells: []string{
				"37.0|-121.5|36.5|-121.0",
				"37.5|-121.5|37.0|-121.0",
				"38.0|-121.5|37.5|-121.0",
				"38.5|-121.5|38.0|-121.0",
				"39.0|-121.5|38.5|-121.0",
				"37.0|-122.0|36.5|-121.5",
				"37.5|-122.0|37.0|-121.5",
				"38.0|-122.0|37.5|-121.5",
				"38.5|-122.0|38.0|-121.5",
				"39.0|-122.0|38.5|-121.5",
				"37.0|-122.5|36.5|-122.0",
				"37.5|-122.5|37.0|-122.0",
				"38.0|-122.5|37.5|-122.0",
				"38.5|-122.5|38.0|-122.0",
				"39.0|-122.5|38.5|-122.0",
				"37.0|-123.0|36.5|-122.5",
				"37.5|-123.0|37.0|-122.5",
				"38.0|-123.0|37.5|-122.5",
				"38.5|-123.0|38.0|-122.5",
				"39.0|-123.0|38.5|-122.5",
				"37.0|-123.5|36.5|-123.0",
				"37.5|-123.5|37.0|-123.0",
				"38.0|-123.5|37.5|-123.0",
				"38.5|-123.5|38.0|-123.0",
				"39.0|-123.5|38.5|-123.0",
			},
		},
	}
	for x := range tests {
		location := Compute(tests[x].Lat, tests[x].Long, tests[x].Resolution, tests[x].Slice)
		geocells := location.ComputeSet(tests[x].Tier)
		sort.Strings(geocells)
		sort.Strings(tests[x].Geocells)
		for y := range geocells {
			if geocells[y] != tests[x].Geocells[y] {
				t.Errorf("Failing as generated %s != test %s", geocells[y], tests[x].Geocells[y])
			}
		}
	}
}
