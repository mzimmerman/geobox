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
	"fmt"
	"math"
	"strings"
)

type Location struct {
	Lat        float64
	Long       float64
	Slice      int
	Resolution int
	Geocell    string
}

func roundSliceDown(coord, slice float64) float64 {
	remainder := math.Mod(coord, slice)
	if math.Signbit(coord) {
		return coord - remainder
	}
	return coord - remainder + slice
}

func (l Location) ComputeSet(tier int) (geocells []string) {
	geocells = make([]string, (1+(2*tier))*(1+(2*tier)))
	format := fmt.Sprintf("%%0.%df", l.Resolution)
	fresolution := 1.0 / math.Pow10(l.Resolution)
	fslice := float64(l.Slice) * fresolution
	top := roundSliceDown(l.Lat, fslice)
	right := roundSliceDown(l.Long, fslice)
	left := right - fslice
	bottom := top - fslice
	walker := 0
	for x := -tier; x <= tier; x++ {
		lat_diff := fslice * float64(x)
		for y := -tier; y <= tier; y++ {
			long_diff := fslice * float64(y)
			geocell := []string{fmt.Sprintf(format, top+lat_diff), fmt.Sprintf(format, left+long_diff), fmt.Sprintf(format, bottom+lat_diff), fmt.Sprintf(format, right+long_diff)}
			geocells[walker] = strings.Join(geocell, "|")
			walker++
		}
	}
	return geocells
}

func Compute(lat, long float64, resolution, slice int) *Location {
	format := fmt.Sprintf("%%0.%df", resolution)
	location := new(Location)
	location.Lat = lat
	location.Long = long
	location.Slice = slice
	location.Resolution = resolution
	fresolution := 1.0 / math.Pow10(resolution)
	fslice := float64(slice) * fresolution
	adj_lat := roundSliceDown(lat, fslice)
	adj_long := roundSliceDown(long, fslice)
	geocell := []string{fmt.Sprintf(format, adj_lat), fmt.Sprintf(format, adj_long-fslice), fmt.Sprintf(format, adj_lat-fslice), fmt.Sprintf(format, adj_long)}
	location.Geocell = strings.Join(geocell, "|")
	return location
}
