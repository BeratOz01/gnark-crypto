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

package bw6756

//Note: This only works for simple extensions

import (
	"github.com/consensys/gnark-crypto/ecc/bw6-756/fp"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
)

func g1IsogenyXNumerator(dst *fp.Element, x *fp.Element) {
	g1EvalPolynomial(dst,
		false,
		[]fp.Element{
			{15061764643505403874, 13688922241935764410, 13945096449919150208, 8436067060318510061, 4337050741185236015, 1993211805393091759, 13630591835847640353, 16993584682990519360, 8971204976640521592, 11479362202198843551, 7505613342680806802, 2623277435220162},
			{2256652953469429710, 7986151259564652590, 7173717765514014760, 5325211358772542105, 11720722258571317741, 1827380411336239380, 2945628459421726487, 15684476260997112721, 16275963394068506199, 3647627339057733797, 6840443943923060230, 2788437309956548},
			{18446744073709550557, 14704964401664098303, 11309681711681649385, 14698000069680550214, 13454959400917254747, 10801466177653926084, 15239716878514852502, 16581223703780321416, 14202840701532685566, 17948097119401515826, 5922093086188465715, 2268644466718328},
		},
		x)
}

func g1IsogenyXDenominator(dst *fp.Element, x *fp.Element) {
	g1EvalPolynomial(dst,
		true,
		[]fp.Element{
			{9026611813877718838, 6035453010921316536, 17062904658006647585, 11643627583491349159, 14230214847810026090, 5797047097534633924, 14807366920250713864, 14779814887555727504, 10347853169840149009, 1246596934729393539, 13141619286505907834, 2448522848965627},
		},
		x)
}

func g1IsogenyYNumerator(dst *fp.Element, x *fp.Element, y *fp.Element) {
	var _dst fp.Element
	g1EvalPolynomial(&_dst,
		false,
		[]fp.Element{
			{7417, 4223785464154554368, 17883137121165483676, 8662155583981545100, 8014120713948809746, 1201662146881371066, 16399484201235277962, 16741102533822300073, 10092115902138952610, 11731609449449782212, 5432405448762957777, 1529941514692833},
			{17685820612365145165, 10133590496207434412, 4439341389815863358, 16237166311366871531, 12586675296657379027, 5752817071182370007, 18238899548945746049, 5474128448956424977, 4657871854268244383, 3092250989713984389, 13902867206193867696, 2147032802810159},
			{3384979430204144565, 11979226889346978885, 1537204611416246332, 17211189075013588966, 17581083387856976611, 11964442653859134878, 13641814725987365538, 14303342354640893273, 15190573054247983491, 5471441008586600696, 10260665915884590345, 4182655964934822},
			{18446744073709551087, 18441456226093760511, 13174518475280565460, 9763304497739979922, 14890648247077438592, 5778851725779543941, 6863645168616474272, 15668448372570953649, 6955362397092648018, 12310026665076143326, 11127771683818204033, 3310628831074305},
		},
		x)

	dst.Mul(&_dst, y)
}

func g1IsogenyYDenominator(dst *fp.Element, x *fp.Element) {
	g1EvalPolynomial(dst,
		true,
		[]fp.Element{
			{33905, 17146720447904350208, 4439688729395715465, 17578719130422492408, 13009726419073394047, 17697253371943596573, 3126642749599797879, 15395188361529500510, 1711728968787230262, 252749604653387985, 14374699731745910598, 1397801387328302},
			{17266182735847260650, 2756724634261996626, 1058985484926618486, 5542736661445693400, 1636838704864690045, 12564472000181253073, 6593938171842885778, 4786584001734410359, 1113513889954172149, 12176975388293072748, 6351201184167411439, 1553692990780230},
			{8633091367923604897, 14375155055950078505, 17702614661430909603, 11655529750965086232, 7917563376482904218, 16634904018698739975, 9041039154614942318, 11137027547596045016, 12888931343158284942, 15514578667146961408, 4644663504360229534, 2992955351466600},
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

	tv1 := fp.Element{17302715199413996045, 15077845457253267709, 8842885729139027579, 12189878420705505575, 12380986790262239346, 585111498723936856, 4947215576903759546, 1186632482028566920, 14543050817583235372, 5644943604719368358, 9440830989708189862, 1039766423535362} //tv1 = c6

	var tv2, tv3, tv4, tv5 fp.Element
	var exp big.Int
	// c4 = 4835703278458516698824703 = 2^82 - 1
	// q is odd so c1 is at least 1.
	exp.SetBytes([]byte{3, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255})

	tv2.Exp(*v, &exp)
	tv3.Mul(&tv2, &tv2)
	tv3.Mul(&tv3, v)

	// line 5
	tv5.Mul(u, &tv3)

	// c3 = 37877157660731232732990269576663233239936484746509109593426423261538632780449313352717366389444912082695314931794809746268936574949192324351273838279701014606648452884726586254167471840902479876056412368
	exp.SetBytes([]byte{1, 238, 213, 183, 107, 119, 49, 92, 85, 130, 79, 195, 198, 173, 25, 235, 146, 241, 154, 95, 88, 89, 209, 63, 126, 70, 68, 40, 170, 44, 116, 217, 152, 213, 206, 120, 133, 72, 219, 61, 96, 89, 2, 93, 64, 159, 85, 65, 79, 214, 57, 103, 160, 220, 200, 220, 82, 89, 162, 189, 182, 200, 212, 168, 96, 85, 71, 132, 177, 188, 251, 218, 22, 208, 189, 13, 10, 73, 216, 6, 120, 252, 199, 240, 208})
	tv5.Exp(tv5, &exp)
	tv5.Mul(&tv5, &tv2)
	tv2.Mul(&tv5, v)
	tv3.Mul(&tv5, u)

	// line 10
	tv4.Mul(&tv3, &tv2)

	// c5 = 2417851639229258349412352
	exp.SetBytes([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	tv5.Exp(tv4, &exp)

	isQNr := g1NotOne(&tv5)

	tv2.Mul(&tv3, &fp.Element{13990906742184113945, 15879050380504523621, 13768460034940508157, 12337541071329853620, 6296858130192020747, 9289986178217863086, 18403114759403589657, 4546259071787184045, 5504643400205978814, 13830311104669138548, 96107744534255859, 1024735223965534})
	tv5.Mul(&tv4, &tv1)

	// line 15

	tv3.Select(int(isQNr), &tv3, &tv2)
	tv4.Select(int(isQNr), &tv4, &tv5)

	exp.Lsh(big.NewInt(1), 82-2)

	for i := 82; i >= 2; i-- {
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
    z.Set( &fp.Element  { 18446744073709504998, 11529623972028612607, 739483395258014634, 5527028560780200701, 11477868704616895891, 15905434021829949368, 2844651761892435780, 17567410508478669002, 4162242322955979641, 15743938111024983262, 11916654042695069468, 4062866236140222 } )
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
	tv3.Mul(&tv3, &fp.Element{18446744073709458379, 881299893533802495, 4886355625346099349, 6225448195760991771, 6629400315996169345, 12607886696045185322, 7201730065066775519, 1932403901886200506, 8616600553259348813, 6369175937589644082, 7499857803942196586, 3773119276850162})

	tv2NZero := g1NotZero(&tv2)

	// tv4 = Z
	tv4 = fp.Element{18446744073709504998, 11529623972028612607, 739483395258014634, 5527028560780200701, 11477868704616895891, 15905434021829949368, 2844651761892435780, 17567410508478669002, 4162242322955979641, 15743938111024983262, 11916654042695069468, 4062866236140222}

	tv2.Neg(&tv2)
	tv4.Select(int(tv2NZero), &tv4, &tv2)
	tv2 = fp.Element{6087387690755251612, 7643068232434215576, 6195945763281467660, 97569654519975969, 1505434147110560758, 12342644747290341982, 14059794106692380317, 15229664573794943703, 16908793757593141664, 1949816925291208189, 9451095697369482684, 234190359239853}
	tv4.Mul(&tv4, &tv2)

	tv2.Square(&tv3)

	var tv6 fp.Element
	//Standard doc line 10
	tv6.Square(&tv4)

	var tv5 fp.Element
	tv5.Mul(&tv6, &fp.Element{6087387690755251612, 7643068232434215576, 6195945763281467660, 97569654519975969, 1505434147110560758, 12342644747290341982, 14059794106692380317, 15229664573794943703, 16908793757593141664, 1949816925291208189, 9451095697369482684, 234190359239853})

	tv2.Add(&tv2, &tv5)
	tv2.Mul(&tv2, &tv3)
	tv6.Mul(&tv6, &tv4)

	//Standards doc line 15
	tv5.Mul(&tv6, &fp.Element{18446744073709458379, 881299893533802495, 4886355625346099349, 6225448195760991771, 6629400315996169345, 12607886696045185322, 7201730065066775519, 1932403901886200506, 8616600553259348813, 6369175937589644082, 7499857803942196586, 3773119276850162})
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

	return x[0] | x[1] | x[2] | x[3] | x[4] | x[5] | x[6] | x[7] | x[8] | x[9] | x[10] | x[11]

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
