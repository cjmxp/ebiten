// Copyright 2014 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ebiten

import (
	"math"
)

// ColorMDim is a dimension of a ColorM.
const ColorMDim = 5

var colorMI = ColorM{
	initialized: true,
	es: [ColorMDim - 1][ColorMDim]float64{
		{1, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0},
	},
}

// A ColorM represents a matrix to transform coloring when rendering an image.
//
// A ColorM is applied to the source alpha color
// while an Image's pixels' format is alpha premultiplied.
// Before applying a matrix, a color is un-multiplied, and after applying the matrix,
// the color is multiplied again.
//
// The initial value is identity.
type ColorM struct {
	initialized bool
	es          [ColorMDim - 1][ColorMDim]float64
}

func (c ColorM) dim() int {
	return ColorMDim
}

// Element returns a value of a matrix at (i, j).
func (c ColorM) Element(i, j int) float64 {
	if !c.initialized {
		if i == j {
			return 1
		}
		return 0
	}
	return c.es[i][j]
}

// Concat multiplies a color matrix with the other color matrix.
// This returns c.
func (c *ColorM) Concat(other ColorM) ColorM {
	if !c.initialized {
		*c = colorMI
	}
	result := ColorM{}
	mul(&other, c, &result)
	*c = result
	return *c
}

// Add adds a color matrix with the other color matrix.
// This returns c.
func (c *ColorM) Add(other ColorM) ColorM {
	if !c.initialized {
		*c = colorMI
	}
	result := ColorM{}
	add(&other, c, &result)
	*c = result
	return *c
}

// SetElement sets an element at (i, j).
func (c *ColorM) SetElement(i, j int, element float64) {
	if !c.initialized {
		*c = colorMI
	}
	c.es[i][j] = element
}

// Monochrome returns a color matrix to make an image monochrome.
func Monochrome() ColorM {
	const r = 6968.0 / 32768.0
	const g = 23434.0 / 32768.0
	const b = 2366.0 / 32768.0
	return ColorM{
		initialized: true,
		es: [ColorMDim - 1][ColorMDim]float64{
			{r, g, b, 0, 0},
			{r, g, b, 0, 0},
			{r, g, b, 0, 0},
			{0, 0, 0, 1, 0},
		},
	}
}

// ScaleColor returns a color matrix that scales a color matrix by the given color (r, g, b, a).
func ScaleColor(r, g, b, a float64) ColorM {
	return ColorM{
		initialized: true,
		es: [ColorMDim - 1][ColorMDim]float64{
			{r, 0, 0, 0, 0},
			{0, g, 0, 0, 0},
			{0, 0, b, 0, 0},
			{0, 0, 0, a, 0},
		},
	}
}

// TranslateColor returns a color matrix that translates a color matrix by the given color (r, g, b, a).
func TranslateColor(r, g, b, a float64) ColorM {
	return ColorM{
		initialized: true,
		es: [ColorMDim - 1][ColorMDim]float64{
			{1, 0, 0, 0, r},
			{0, 1, 0, 0, g},
			{0, 0, 1, 0, b},
			{0, 0, 0, 1, a},
		},
	}
}

// RotateHue returns a color matrix to rotate the hue
func RotateHue(theta float64) ColorM {
	sin, cos := math.Sincos(theta)
	v1 := cos + (1.0-cos)/3.0
	v2 := (1.0/3.0)*(1.0-cos) - math.Sqrt(1.0/3.0)*sin
	v3 := (1.0/3.0)*(1.0-cos) + math.Sqrt(1.0/3.0)*sin
	return ColorM{
		initialized: true,
		es: [ColorMDim - 1][ColorMDim]float64{
			{v1, v2, v3, 0, 0},
			{v3, v1, v2, 0, 0},
			{v2, v3, v1, 0, 0},
			{0, 0, 0, 1, 0},
		},
	}
}