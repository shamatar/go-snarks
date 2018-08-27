package verifier

import (
	"math/big"
	"testing"
)

// TestVerification tests basic snark verification from ZoKrates
func TestVerification(t *testing.T) {
	vk := &VerifyingKey{}
	// grab some example from ZoKrates
	vk.A, _ = NewG2FromStrings(
		[2]string{"0x95aacb3c23ba931cea76b29a19617c347943ecf1002fc0fcf798e0f8d74846c",
			"0x6d039a9508f5e58dd3646aec40b650cf2622990ef101628046382d9d0a11c2b"},
		[2]string{"0x1396a7eda2554a0ddef4978de4de39b8618c0d170eaa3b4e3528c1e5e44c7050",
			"0x35186be372d0ef1ca7821d2eeb6575accda23f59fe655f4d087031b75ca9d48"},
		16)
	vk.B, _ = NewG1FromStrings("0x23b16b75c253c67a07b2dfb4642f916ce483f52d7d957c43f1cafa9e8efbd4b6",
		"0x10b9cf66c4b02c8eedc287dbfabde1d7e1be9c817eced0610e4232063b473347",
		16)
	vk.C, _ = NewG2FromStrings(
		[2]string{"0x8d841a59c62049761286964354fbb9547aacb249038adcb0b8b393a0a71f872",
			"0x23bd90ef5307660e965fd3144418d2b5a8e6c22a87900541e12745ea62f6375f"},
		[2]string{"0x10935f302a61d826b52731e5951219552832fe1be64ef191364fb8a264036d67",
			"0x608317e583295bdc60d86315f7923e2677143c6eeb381f8f4b828bf1a055fcd"}, 16)
	vk.gamma, _ = NewG2FromStrings(
		[2]string{"0x838d4dcd872e6bcc73ef8d9f2d61c1b1223672a4ce4483d14aea08acc895106", "0xc700c249203de561efcda3fd4ec7658d528701da3c1fca098a2d42f011fde37"},
		[2]string{"0x2979bf58a5a9ef40265dea31187fd807257683846c5832c7b87951a78aa68841", "0x19e61fbd9e5b5e25dd7d110cd4b0071123975620e37c3cd72a101b63eff4f185"},
		16)
	vk.gammaBeta1, _ = NewG1FromStrings("0x22dc5c5252437b5f3feb464de52cb4306a0b45bbad6a4a19b36b1f991eabbe47",
		"0x177a001e3312096da4d75bfb2b4cc07bec195dccb48ecc6808f81c5f232e02e", 16)
	vk.gammaBeta2, _ = NewG2FromStrings(
		[2]string{"0x723f2cb7ff1065bbb44c112f68e1b632de6b6be20920beae4f28d9e0057cc47", "0xdd6e0ed155dbbc30988ad7bfd9cbcab76850d09bb6b668c1b831354fe2c1180"},
		[2]string{"0x2fbb278596ffa9342d987295ada995129d5ed9b484a3337cbd095b38f1a834bd", "0x19bb949e51477e54f9ae892093beda448c33b21b43366502cd4d8a02897017a"},
		16)
	vk.Z, _ = NewG2FromStrings(
		[2]string{"0x15425ab686d62e846de5ca98c92843cad80f06401610fcce27cf6f94e07f7f8b", "0xe21bc3f9c196bb742aa2404a6550635fd632b86cb0fbae13962df909543c7c3"},
		[2]string{"0x6538ba629c4809a327453895dade6110275faa25e14e4ce6390c3b095f6519c", "0x5aaf52260e7758c6cd487c1e88ef9791b518fbf09e0ec5402825687b75f3ab9"},
		16)
	vk.IC = make([]*G1, 7)
	vk.IC[0], _ = NewG1FromStrings("0x406d95f7ae8feb816b5a4b4cf4a949730f31de0b3ff1e8b800db9969514a256", "0xe0818f24fef31a82b9e7e587c9934f812e5a196a5e7aa7d495a8494f9557f4c", 16)
	vk.IC[1], _ = NewG1FromStrings("0xf394dc0e39f79bee4a1e0577299de796c7a03ad6e3949884d2d2fbd73f0d76c", "0x16a514b14810b2217478d803fb68372e34dde4fb624b9f4e156470364cb19402", 16)
	vk.IC[2], _ = NewG1FromStrings("0x134d26c1eec281bedabb8d4ee3fb61d2e9113668b3fa45f110d8547fd8f5b94d", "0x59f3d50872f8ee5662e25fc1f13e08acc01f5a7ac049b3661b87cd2b78bdd11", 16)
	vk.IC[3], _ = NewG1FromStrings("0xcd6823439d212cdf695ae58c3a9cb50bc31f285e986a1ca00aed91669401419", "0xbea35dbd15dd6e212050e923a6f748a981ea0f1bca31305e87a487c5c07f506", 16)
	vk.IC[4], _ = NewG1FromStrings("0x139af07e752cc39bc1a2c6b0468415afa072bddd5599f38e7cc4031031a5ec02", "0x27233e6a491cbdb09fd5b3b94917658ab20b2207f3d7e6aa556ca68bff5af037", 16)
	vk.IC[5], _ = NewG1FromStrings("0x20a7eda84420d312d03674aff46ae1202db426e9ae4ea1d400e385973ea31497", "0x61e54c3741544a6a99149a282cf1c9c08e1dd6061a593260bebba7d5c5e0b55", 16)
	vk.IC[6], _ = NewG1FromStrings("0x8ed7634c368516b852917dac2d165ac153204b550c5ba0456512066daef108b", "0x14c048afac3cdb256d375825101691093f73aeaf6802df2fba3c0e83de9feb86", 16)

	// Proof:
	proof := &Proof{}
	proof.A, _ = NewG1FromStrings("0x29a6ef0f8e73e5c389221e262b7f1695cd71d064667240bbb8a8aa143a5e3a3a", "0x2ce302d2d95bef56c48d6eba6c661e7190606e465e1bda2595f3782ab76901ed", 16)
	proof.Ap, _ = NewG1FromStrings("0x521375a85479309045805ca35252fa42e8eab986658d942c029caf7da2b53a", "0x85348b0d8402170a9f8d0358eaf8e0cdf5bff2d1bcad770e909045c434d87a2", 16)
	proof.B, _ = NewG2FromStrings(
		[2]string{"0x261f0899cd9ac24f6ea6b06a4d7db6fcebf39a632f4ffe99f85b85c2a9058127", "0x1ef64852f0192760310dda76a07da1bb0124869cd8f6fada287ce612e2c1987b"},
		[2]string{"0x611c4c383c5e7d261ed96a1f5c145c940b1bcbdec28b74d58e3ee4f51658733", "0x2df7f862dfc698b0b944c7c8fab618c0d82a22c666632b060b4c3f92ecd80e0b"}, 16)
	proof.Bp, _ = NewG1FromStrings("0x173250553b786a5125ee76b1e539d36aa6209c86711485ccb28d66149869756b", "0x2d06f0def23ee2f9b0f0f32093d7f034f300c7abb762065dc2784052e8b635a7", 16)
	proof.C, _ = NewG1FromStrings("0x279279bc0ce40f286e300f5c339620d02c83a48a9a2a82cfc2dc386941508969", "0x3040f8131a73d41701b66257a3f23b0be0f01dcc1ee36ab65cc8e4c947331262", 16)
	proof.Cp, _ = NewG1FromStrings("0x1765ef75ba208e9c3ac184254e0b5386c5257c15150d8fbe22d22922343be9ab", "0x7ba5ae37e01dad31e8002be5839f475a35001c4dc6c02af2bee529a10ee8a7d", 16)
	proof.H, _ = NewG1FromStrings("0x2892ea5a02c1c3304f48e6dea85028916ff8796d9e7ac380b7ff1aac11f05326", "0x2425d5f1649f7d52ebdb68a2bdc70af0abe2d0146eb8d1ec28fd8f0c2f818c8a", 16)
	proof.K, _ = NewG1FromStrings("0x2f768eb3ffd67d561d0dd8cf09648fa3585080d32fff605cfc1fb15b83c6c2c6", "0x8d7bbf5a99e154824ecf6e2099cb50f5b16d053601b6c6c2824380e7bb5ae20", 16)

	// Witness:
	witness := make([]*big.Int, 6)
	witness[0] = big.NewInt(0) // nonce 0
	witness[1] = big.NewInt(1) // fib 1
	witness[2] = big.NewInt(1) // fib 1
	witness[3] = big.NewInt(1) // nonce 1
	witness[4] = big.NewInt(1) // fib 1
	witness[5] = big.NewInt(2) // fib 2

	err := naiveSplitVerification(witness, proof, vk)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTrivialPairingCheck(t *testing.T) {
	scalar := big.NewInt(2)
	G1Point := GetG1Base()
	G2Point := GetG2Base()

	TwoG1 := new(G1).ScalarMult(G1Point, scalar)
	TwoG2 := new(G2).ScalarMult(G2Point, scalar)

	temp := new(G1)

	success := PairingCheck([]*G1{G1Point, temp.Neg(TwoG1)}, []*G2{TwoG2, G2Point})
	if !success {
		t.Fatal("Failed trivial pairing check")
	}
}
