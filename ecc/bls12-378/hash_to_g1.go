// Copyright 2020 ConsenSys Software Inc.
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

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bls12378

//Note: This only works for simple extensions

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-378/fp"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
)

func g1IsogenyXNumerator(dst *fp.Element, x *fp.Element) {
	g1EvalPolynomial(dst,
		false,
		[]fp.Element{
			{6470543203100547353, 14665988170059032885, 1939515187828768580, 3603560708219821999, 1095026321559208988, 183854362073986052},
			{14854739450145697573, 10725610731781000990, 8384146146676813741, 3792524792234600218, 10980132992245118504, 127811114313784269},
			{7743341424938057712, 1621446876320497654, 9161388112343345155, 13809747965729688982, 14072110284352226775, 77964074263176954},
		},
		x)
}

func g1IsogenyXDenominator(dst *fp.Element, x *fp.Element) {
	g1EvalPolynomial(dst,
		true,
		[]fp.Element{
			{11480213446153845907, 9569059723295472762, 4133212223950692662, 5656914875334883651, 989021686692303102, 227886835744873894},
		},
		x)
}

func g1IsogenyYNumerator(dst *fp.Element, x *fp.Element, y *fp.Element) {
	var _dst fp.Element
	g1EvalPolynomial(&_dst,
		false,
		[]fp.Element{
			{4780610586872381554, 18423150125994475588, 13123772012819260137, 14591853195646969647, 5804992647820306014, 20966723178287685},
			{2213410472918458692, 9774999437004642556, 13319410716190829459, 14774058963459238058, 12357489533274646735, 248758148590692830},
			{16759481071713625783, 8645096532612011693, 7097905075491715268, 932195041550141716, 4227816384078368107, 50037860715544812},
			{3871670712469028856, 10034095475015024635, 4580694056171672577, 16128246019719620299, 7036055142176113387, 38982037131588477},
		},
		x)

	dst.Mul(&_dst, y)
}

func g1IsogenyYDenominator(dst *fp.Element, x *fp.Element) {
	g1EvalPolynomial(dst,
		true,
		[]fp.Element{
			{17641076928456688137, 8306475833976684855, 8359419817241119003, 12641605213272639883, 9863039736160487870, 55368217170706106},
			{17969127027499495035, 2884196005202455728, 15703691879418613809, 10094567750702434230, 12004334191193297464, 175264043429118194},
			{12350127924441855415, 17380644983358010734, 8933124167467608227, 16391120112507168124, 9337764864048325557, 116945264214095313},
		},
		x)
}

func g1Isogeny(p *G1Affine) {

	den := make([]fp.Element, 2)

	g1IsogenyYDenominator(&den[1], &p.X)
	g1IsogenyXDenominator(&den[0], &p.X)

	g1IsogenyYNumerator(&p.Y, &p.X, &p.Y)
	g1IsogenyXNumerator(&p.X, &p.X)

	den = fp.BatchInvert(den)

	p.X.Mul(&p.X, &den[0])
	p.Y.Mul(&p.Y, &den[1])
}

// g1SqrtRatio computes the square root of u/v and returns 0 iff u/v was indeed a quadratic residue
// if not, we get sqrt(Z * u / v). Recall that Z is non-residue
// The main idea is that since the computation of the square root involves taking large powers of u/v, the inversion of v can be avoided
func g1SqrtRatio(z *fp.Element, u *fp.Element, v *fp.Element) uint64 {

	// Taken from https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/13/ F.2.1.1. for any field

	tv1 := fp.Element{3422016347327078217, 15952935974507985473, 10210560017327941857, 6195437588884472512, 1531492004832937820, 17090488542823369} //tv1 = c6

	var tv2, tv3, tv4, tv5 fp.Element
	var exp big.Int
	// c4 = 2199023255551 = 2^41 - 1
	// q is odd so c1 is at least 1.
	exp.SetBytes([]byte{1, 255, 255, 255, 255, 255})

	tv2.Exp(*v, &exp)
	tv3.Mul(&tv2, &tv2)
	tv3.Mul(&tv3, v)

	// line 5
	tv5.Mul(u, &tv3)

	// c3 = 137617509170765099891752579783724504691201148437113468788429769127729045045134922651478473733013131816
	exp.SetBytes([]byte{251, 172, 16, 89, 161, 52, 100, 20, 242, 215, 73, 3, 180, 65, 232, 161, 1, 103, 173, 145, 196, 8, 201, 166, 3, 112, 216, 52, 41, 39, 95, 243, 165, 253, 218, 160, 139, 0, 0, 38, 82, 40})
	tv5.Exp(tv5, &exp)
	tv5.Mul(&tv5, &tv2)
	tv2.Mul(&tv5, v)
	tv3.Mul(&tv5, u)

	// line 10
	tv4.Mul(&tv3, &tv2)

	// c5 = 1099511627776
	exp.SetBytes([]byte{1, 0, 0, 0, 0, 0})
	tv5.Exp(tv4, &exp)

	isQNr := g1NotOne(&tv5)

	tv2.Mul(&tv3, &fp.Element{17614810958234635860, 11393801269165528284, 8781501035240632779, 8106712880529013806, 4971838157288047198, 122121039825317715})
	tv5.Mul(&tv4, &tv1)

	// line 15

	tv3.Select(int(isQNr), &tv3, &tv2)
	tv4.Select(int(isQNr), &tv4, &tv5)

	exp.Lsh(big.NewInt(1), 41-2)

	for i := 41; i >= 2; i-- {
		//line 20
		tv5.Exp(tv4, &exp)
		nE1 := g1NotOne(&tv5)

		tv2.Mul(&tv3, &tv1)
		tv1.Mul(&tv1, &tv1)
		tv5.Mul(&tv4, &tv1)

		tv3.Select(int(nE1), &tv3, &tv2)
		tv4.Select(int(nE1), &tv4, &tv5)

		exp.Rsh(&exp, 1)
	}

	*z = tv3
	return isQNr
}

func g1NotOne(x *fp.Element) uint64 {

	var one fp.Element
	return one.SetOne().NotEqual(x)

}

/*
// g1SetZ sets z to [11].
func g1SetZ(z *fp.Element) {
    z.Set( &fp.Element  { 5249763402351377716, 3384457438931451475, 13367120442609335946, 13855353052415766542, 11761008755492169078, 30127809456627797 } )
}*/

// g1MulByZ multiplies x by [11] and stores the result in z
func g1MulByZ(z *fp.Element, x *fp.Element) {

	res := *x

	res.Double(&res)
	res.Double(&res)
	res.Add(&res, x)
	res.Double(&res)
	res.Add(&res, x)

	*z = res
}

// From https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/13/ Pg 80
// sswuMapG1 implements the SSWU map
// No cofactor clearing
func sswuMapG1(u *fp.Element) G1Affine {

	var tv1 fp.Element
	tv1.Square(u)

	//mul tv1 by Z
	g1MulByZ(&tv1, &tv1)

	var tv2 fp.Element
	tv2.Square(&tv1)
	tv2.Add(&tv2, &tv1)

	var tv3 fp.Element
	//Standard doc line 5
	var tv4 fp.Element
	tv4.SetOne()
	tv3.Add(&tv2, &tv4)
	tv3.Mul(&tv3, &fp.Element{10499526804702755432, 6768914877862902950, 8287496811509120276, 9263962031121981469, 5075273437274786541, 60255618913255595})

	tv2NZero := g1NotZero(&tv2)

	// tv4 = Z
	tv4 = fp.Element{5249763402351377716, 3384457438931451475, 13367120442609335946, 13855353052415766542, 11761008755492169078, 30127809456627797}

	tv2.Neg(&tv2)
	tv4.Select(int(tv2NZero), &tv4, &tv2)
	tv2 = fp.Element{15314533651602404840, 3999629397495592995, 17991228730268553058, 13253234862282888158, 4784493033884022421, 276795783356562829}
	tv4.Mul(&tv4, &tv2)

	tv2.Square(&tv3)

	var tv6 fp.Element
	//Standard doc line 10
	tv6.Square(&tv4)

	var tv5 fp.Element
	tv5.Mul(&tv6, &fp.Element{15314533651602404840, 3999629397495592995, 17991228730268553058, 13253234862282888158, 4784493033884022421, 276795783356562829})

	tv2.Add(&tv2, &tv5)
	tv2.Mul(&tv2, &tv3)
	tv6.Mul(&tv6, &tv4)

	//Standards doc line 15
	tv5.Mul(&tv6, &fp.Element{10499526804702755432, 6768914877862902950, 8287496811509120276, 9263962031121981469, 5075273437274786541, 60255618913255595})
	tv2.Add(&tv2, &tv5)

	var x fp.Element
	x.Mul(&tv1, &tv3)

	var y1 fp.Element
	gx1NSquare := g1SqrtRatio(&y1, &tv2, &tv6)

	var y fp.Element
	y.Mul(&tv1, u)

	//Standards doc line 20
	y.Mul(&y, &y1)

	x.Select(int(gx1NSquare), &tv3, &x)
	y.Select(int(gx1NSquare), &y1, &y)

	y1.Neg(&y)
	y.Select(int(g1Sgn0(u)^g1Sgn0(&y)), &y, &y1)

	//Standards doc line 25
	x.Div(&x, &tv4)

	return G1Affine{x, y}
}

// MapToG1 invokes the SSWU map, and guarantees that the result is in g1
func MapToG1(u fp.Element) G1Affine {
	res := sswuMapG1(&u)
	//this is in an isogenous curve
	g1Isogeny(&res)
	res.ClearCofactor(&res)
	return res
}

// EncodeToG1 hashes a message to a point on the G1 curve using the Simplified Shallue and van de Woestijne Ulas map.
// It is faster than HashToG1, but the result is not uniformly distributed. Unsuitable as a random oracle.
// dst stands for "domain separation tag", a string unique to the construction using the hash function
//https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/13/#section-6.6.3
func EncodeToG1(msg, dst []byte) (G1Affine, error) {

	var res G1Affine
	u, err := hashToFp(msg, dst, 1)
	if err != nil {
		return res, err
	}

	res = sswuMapG1(&u[0])

	//this is in an isogenous curve
	g1Isogeny(&res)
	res.ClearCofactor(&res)
	return res, nil
}

// HashToG1 hashes a message to a point on the G1 curve using the Simplified Shallue and van de Woestijne Ulas map.
// Slower than EncodeToG1, but usable as a random oracle.
// dst stands for "domain separation tag", a string unique to the construction using the hash function
// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#section-3
func HashToG1(msg, dst []byte) (G1Affine, error) {
	u, err := hashToFp(msg, dst, 2*1)
	if err != nil {
		return G1Affine{}, err
	}

	Q0 := sswuMapG1(&u[0])
	Q1 := sswuMapG1(&u[1])

	//TODO: Add in E' first, then apply isogeny
	g1Isogeny(&Q0)
	g1Isogeny(&Q1)

	var _Q0, _Q1 G1Jac
	_Q0.FromAffine(&Q0)
	_Q1.FromAffine(&Q1).AddAssign(&_Q0)

	_Q1.ClearCofactor(&_Q1)

	Q1.FromJacobian(&_Q1)
	return Q1, nil
}

// g1Sgn0 is an algebraic substitute for the notion of sign in ordered fields
// Namely, every non-zero quadratic residue in a finite field of characteristic =/= 2 has exactly two square roots, one of each sign
// Taken from https://datatracker.ietf.org/doc/draft-irtf-cfrg-hash-to-curve/ section 4.1
// The sign of an element is not obviously related to that of its Montgomery form
func g1Sgn0(z *fp.Element) uint64 {

	nonMont := *z
	nonMont.FromMont()

	return nonMont[0] % 2

}

func g1EvalPolynomial(z *fp.Element, monic bool, coefficients []fp.Element, x *fp.Element) {
	dst := coefficients[len(coefficients)-1]

	if monic {
		dst.Add(&dst, x)
	}

	for i := len(coefficients) - 2; i >= 0; i-- {
		dst.Mul(&dst, x)
		dst.Add(&dst, &coefficients[i])
	}

	z.Set(&dst)
}

func g1NotZero(x *fp.Element) uint64 {

	return x[0] | x[1] | x[2] | x[3] | x[4] | x[5]

}

// hashToFp hashes msg to count prime field elements.
// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#section-5.2
func hashToFp(msg, dst []byte, count int) ([]fp.Element, error) {

	// 128 bits of security
	// L = ceil((ceil(log2(p)) + k) / 8), where k is the security parameter = 128
	const L = 16 + fp.Bytes

	lenInBytes := count * L
	pseudoRandomBytes, err := ecc.ExpandMsgXmd(msg, dst, lenInBytes)
	if err != nil {
		return nil, err
	}

	res := make([]fp.Element, count)
	for i := 0; i < count; i++ {
		res[i].SetBytes(pseudoRandomBytes[i*L : (i+1)*L])
	}
	return res, nil
}
