package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestGitInfoRefs(t *testing.T) {
	example := strings.Join([]string{
		"001e# service=git-upload-pack",
		"000000b3c7d3d3371baa35587fb66d8a79c6d999a4dafd8e HEAD\000multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done agent=git/1.8.4",
		"003cd58c4a91450924a963d2cc7407dfa3e38866cb06 refs/heads/0.2",
		"004ca8fde08d941efc61322ae0302bf8bafc13e2275c refs/heads/add-version-to-join",
		"003fc7d3d3371baa35587fb66d8a79c6d999a4dafd8e refs/heads/master",
		"004448da4910b78e24d8d3a831839cc751700ddc6e10 refs/heads/update-docs",
		"003ee2f04208620f3bc9d01cc3fb216b92fa4e4a5767 refs/pull/1/head",
		"003f9b36b682ebbd7bd224b621fb90864821726b11b3 refs/pull/1/merge",
		"003f1eb0be10fe9ebf6e99a6c16abd3e583a68533dbd refs/pull/10/head",
		"0000",
	}, "\n")

	in := strings.NewReader(example)
	gup, err := parseGitUploadPack(ioutil.NopCloser(in))

	if err != nil {
		t.Fatalf("Failed parding git-upload-pack: %v", err)
	}

	// Stringification spits something similar back out...
	if len(example) != len(gup.String()) {
		t.Errorf(
			"Stringified doc should be length %v, got %v",
			len(example), len(gup.String()),
		)
		t.Errorf(
			"Expected String() to return \n%v\nGot:\n%v",
			example, gup.String(),
		)
	}
	/*
		if gup.String() != "001e# service=git-upload-pack\n0000" {
			t.Errorf("Didn't get the right output.")
		}
	*/

	if gup.capabilities != "multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done agent=git/1.8.4" {
		t.Errorf(
			"Expected capabilities to be \n\t%v\nGot:\n\t%v",
			"multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done agent=git/1.8.4",
			gup.capabilities,
		)
	}

	// HEAD
	if gup.refs["HEAD"] != "c7d3d3371baa35587fb66d8a79c6d999a4dafd8e" {
		t.Errorf(
			"Expected refs.HEAD to be \n\t%v\nGot:\n\t%v",
			"c7d3d3371baa35587fb66d8a79c6d999a4dafd8e",
			gup.refs["HEAD"],
		)
	}

	// Can find the tag "update-docs"
	// TODO: Table of stuff we should test here
	var tests = []struct {
		q string
		c string
	}{
		{"update-docs", "48da4910b78e24d8d3a831839cc751700ddc6e10"},
		{"9b36b682", "9b36b682ebbd7bd224b621fb90864821726b11b3"},
		{"1111111111111111111111111111111111111111", "1111111111111111111111111111111111111111"},
	}

	for _, tt := range tests {
		err, commit := gup.findCommitish(tt.q)
		if err != nil || commit != tt.c {
			t.Errorf("Could not find commitish '%v'; got %v and commit %v.", tt.q, err, commit)
		}
	}

	// Set commit at someting bogus blows up
	err = gup.SetMaster("does-not-exist")
	if err == nil {
		t.Errorf("Expected SetHead(does-not-exist) to fail. It didn't.")
	}

	err = gup.SetMaster("update-docs")
	if err != nil {
		t.Errorf("Expected SetMaster(update-docs) to work, failed with %v", err)
	} else if gup.refs["refs/heads/master"] != "48da4910b78e24d8d3a831839cc751700ddc6e10" {
		t.Errorf(
			"Expected master to be at %v, but got %v",
			"004448da4910b78e24d8d3a831839cc751700ddc6e10",
			gup.refs["refs/heads/master"],
		)
	}
}

func TestLitmus(t *testing.T) {
	in := strings.Join([]string{
		`001e# service=git-upload-pack`,
		"000000b33be69f09ac6b5155cdf487e5d94a165098b33705 HEAD\x00multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done agent=git/1.8.4",
		`003cbd0d2e84945601fd2d76ebcfd5081d78936026a1 refs/heads/0.2
004ca8fde08d941efc61322ae0302bf8bafc13e2275c refs/heads/add-version-to-join
003f3be69f09ac6b5155cdf487e5d94a165098b33705 refs/heads/master
004448da4910b78e24d8d3a831839cc751700ddc6e10 refs/heads/update-docs
003ee2f04208620f3bc9d01cc3fb216b92fa4e4a5767 refs/pull/1/head
003f9b36b682ebbd7bd224b621fb90864821726b11b3 refs/pull/1/merge
003f1eb0be10fe9ebf6e99a6c16abd3e583a68533dbd refs/pull/10/head
0040333d2ab7d90ed61353aac2870662752487fcfc38 refs/pull/10/merge
0040928781aaa37362a43188ba8caddb792756612e6e refs/pull/100/head
004160d2424bff23a0282bee48b804f9ebb9a8c5b07d refs/pull/100/merge
0040d88bfc084b63145b14333492305c95ecd14fdb63 refs/pull/102/head
0041076b922cb851a6922374eb01dfda96b782990ea0 refs/pull/102/merge
004082fe001c65c26082faf2028af5402e94ece9a571 refs/pull/108/head
0041c017aa70c225c1f9acdbac231c226762efc09e57 refs/pull/108/merge
0040915266d5f50169e306e8d0113828e02dd0253de6 refs/pull/109/head
0041ad6a77a94d28fb5c5468aaeedef58d8d22d94349 refs/pull/109/merge
003f931afb93949b75209ad891102239fc973df98c4b refs/pull/11/head
0040cff225ffd2b8f3b73418b87dacf284430fe817ec refs/pull/11/merge
00400c42a03267c29bd8aa1a3bd5d5d210e89f05952b refs/pull/111/head
0041f5da44647712ca78bbd698b5c05d56f1c09de34b refs/pull/111/merge
0040f5d90994e814a2da91721a9a5a076ce35510a1bc refs/pull/112/head
004159ad708c3c831b903da664558f21ccad87134044 refs/pull/112/merge
00409ede78d75fc13e82707f5fa25299d49687d23528 refs/pull/114/head
00413be470d468b03c1a8536501b6d6c3dd31a709a1b refs/pull/114/merge
0040ac9801f5702045badcf9b4dbff00cf469dd148c0 refs/pull/117/head
00414358250fa059a47f3fbd20a32b5adbb5585f9051 refs/pull/117/merge
00403e59badf1ae36b7819fe022f3849f60761b27f73 refs/pull/118/head
0041ff18583eda39e19533ea67634ae6f8ec895ac9a2 refs/pull/118/merge
003f45af72c941be41f259b82b2762118b733d88014c refs/pull/12/head
0040de5cdd6e4e1efa30ff325139c291b06d0fe7a668 refs/pull/12/merge
0040bd65b4570435f9efc2c74ddd088dc9b63befb06a refs/pull/120/head
0041757633aae7aca86759f9472ba9af5d846d33dc09 refs/pull/120/merge
004070f25901279789b48d8fd3858ae7e395fdda4fd0 refs/pull/121/head
0041c0041f073e4b62e87b699873560ed868d2962fd9 refs/pull/121/merge
00402c09cd7d1a927fa30b0b8c6b6c1c939d7db92c52 refs/pull/122/head
0041ff930c92e00fdef400e2695ffc0d620f6ac1c9d3 refs/pull/122/merge
0040037a9c9957196ac8b68bf053d00e5b436dde92f0 refs/pull/123/head
00419987be78c8eddd7501234717ba26ab40f68b22ad refs/pull/123/merge
0040e50871cc36249f1b87b65f365f519bd5e4b37088 refs/pull/124/head
00410ac9c9bbb5e13247d566f3d3cb75b5d4d53ae27a refs/pull/124/merge
0040bcc77db8a93eb5a9f31bcb82d2968128203eeb09 refs/pull/126/head
0041765566ded3a611c11bd7fa5cbc2dd789c3a3cb03 refs/pull/126/merge
00401527b7008ce63f047fa6bc7d56a827ae874679cd refs/pull/127/head
004169de793cd6d6bc7416b06742642eb2dc69f830f3 refs/pull/127/merge
0040f813017f1bcbb7be21707f30cbea4b4af9f66264 refs/pull/128/head
004159230c7a4f287493ba22243b5edbc9f7ada1ee99 refs/pull/128/merge
0040bfc68e8e37fe8bedb2c9b73a5d521ec18a76efec refs/pull/129/head
003fc2a80df3f940f8c2e1271bb28f44a0e82275c783 refs/pull/13/head
00403aca49e460524e9f49ffd5a751da1fe8db295a13 refs/pull/13/merge
0040b430a07e1bb5931eecf7ceb2235c4f48d8b18c2d refs/pull/130/head
004129603c1bba48524d7992d77c26758d731f74cbae refs/pull/130/merge
00402991bf58e1b5fc60dd795e713c6dd0b113bf86a8 refs/pull/131/head
00410659f4d290bb28fe6b8a1b5deaca74c68ddccb24 refs/pull/131/merge
0040e848659db6dad7cd95ec2be1de93bc4ccc16cabb refs/pull/132/head
0041c584685cd06d2f957aeee81b75b850087fc03534 refs/pull/132/merge
004057ef6e9f5a2a4fa50a7b1ca1e1d7170045d4cd5f refs/pull/133/head
0041f18bc625a4d448ff974871e5d4386a114079420d refs/pull/133/merge
0040dd2f856d63fc1c9d0850c3f18e04fb9b7443d143 refs/pull/134/head
0041ecc27be5d401fe8a712e3538c8397f32408875ae refs/pull/134/merge
004086e03d22982ab382d59e726d86f3ea3b14eea29a refs/pull/135/head
004117c5a4c960aefb16c672643015a6f1cc826de358 refs/pull/135/merge
00406add32f1c8e3341b21183bbb11a2becd32a3560f refs/pull/136/head
0041d28a9de86c892eaf9a594bc67c72b0138ee83507 refs/pull/136/merge
00407563a13621ad61a33bfe7e158e47b1565ee8c180 refs/pull/137/head
004130215516ca06d09f0337312c56b6be94387bfe96 refs/pull/137/merge
0040800c4718c1691ddec877c5828c69dc3c112b9088 refs/pull/138/head
00412ed2efcb845c877604c34825276cb1faa37f5389 refs/pull/138/merge
003f89bacedf3e88810388ba657e95f86f323691fa68 refs/pull/14/head
004072c945cca229c37a9d056253bb863414ae1aa39f refs/pull/14/merge
0040e856acf05e828fcfb31687232153ac294a55811f refs/pull/140/head
0041bf98444cc1b23c5592b9f7b980d64b03b912130a refs/pull/140/merge
00404f436ae70ae21429baea18ec5bb62c88a2ac9b0e refs/pull/141/head
0041f6f20a3f324fb3a42cfdb427c783e94d0d5b4141 refs/pull/141/merge
0040a543d644b401b3aaa36d7219700de595853bb874 refs/pull/142/head
00419a9450a509f623017f67b9a723af5a7b7a7c9780 refs/pull/142/merge
00406108f8536f932dad27947b77ec4a9c6cf8aabbec refs/pull/143/head
0041e6df29189f64eb76c1721069e2de490ded48712d refs/pull/143/merge
0040351e84aece21dade1818e07e0ce56a23380a33ef refs/pull/145/head
00414262a0584e177da94f2be51cbdd3e2d066767085 refs/pull/145/merge
0040b8d85e627e86dffec240cbec293af93c921ab838 refs/pull/147/head
0041d891eb50e971e29d8e5d0bf5885f6b42cf0e4f35 refs/pull/147/merge
003f3e436473036f4c757403d86dfd009017fc554184 refs/pull/15/head
0040ede395b0e8b18e838987d5d2f2853a87564fab22 refs/pull/15/merge
0040a22bd2b8b27b095b3d19f7a62f4fdcd93cf95fb2 refs/pull/150/head
00411f7a3a1ee5e1e5e81da0165e882dc564a1162f08 refs/pull/150/merge
0040de0a8c60ac52073267df2ced3d0250bee885db41 refs/pull/151/head
0041ddc46e90bed5757fc76ca81b224783ccc5182027 refs/pull/151/merge
004090d7ebec47967e76c4986c71f2a5b80e78b7b3c1 refs/pull/152/head
00411ea349d3a3b29a11b62db61ab5cfe64b27c996a6 refs/pull/152/merge
00402f5015552e8db84d04cb75b0adb3f7b9d3be5c56 refs/pull/153/head
0041564994e657d31b1ce737bba7026add165666f43f refs/pull/153/merge
00402022c4bce6ad89cbac58854c717bd54014fc4199 refs/pull/154/head
0041e7d5a3552bedbe18925d342300c82d429ac59149 refs/pull/154/merge
0040b366f1044688e3ee8cacc09c4fd1b9ffeb399167 refs/pull/159/head
0041a235359297e274e2c41541c4123fa22edf6077cd refs/pull/159/merge
003f6e669dab8ea1ea53b0524488492e3890fad1ebea refs/pull/16/head
004083e330524f66addece65a49faa9a3d46d3a436f9 refs/pull/16/merge
0040a623effaf14bb3a035f749df0ec180e0eb6ebb4e refs/pull/160/head
00416606e166f82995575f98eb7411fe96fc774800ce refs/pull/160/merge
0040380326b5d1ba26a9c5943a8c09072ab45e1b9d1e refs/pull/161/head
004144f02da86ac8fbc657a03fd90ab93be0502da34d refs/pull/161/merge
00401427a73f9b9b75b648174a8952dc37ea85639a8e refs/pull/167/head
00417ca2754d787329e9edd4c76387b83eac1c194a05 refs/pull/167/merge
00408d245b546f93761be35b465e189f53df78be723e refs/pull/168/head
00411e30eae8323ee22a01fb3575fd0651ea15515e46 refs/pull/168/merge
003f346da34a0002ecb0f0edbb52e6b354a5c355b456 refs/pull/17/head
0040f5d7a7b8c2c68b42ef6fb2fbf537d39a1e91fc8d refs/pull/17/merge
0040f9235481829a53fbefb5788971afa807468cba1e refs/pull/171/head
004154a01c147b2d6decc2c1a47270bb066a3b134fa7 refs/pull/171/merge
00402055fa8048c4cc39da3d6bd5e2e36e0c3a1b1297 refs/pull/172/head
004196956849f428e781307e26060a6e636c787c6317 refs/pull/172/merge
004062b8b7a6a88b50548864b8553ca188ad9cb1db30 refs/pull/173/head
0041fa4d66164925f5fb004ec2625d9f65a30345d5f3 refs/pull/173/merge
0040f03481f7335868c8780c1be939f3807a67cdceeb refs/pull/175/head
0041e4646204ad41b9eaa4111f3054dda278b39e6763 refs/pull/175/merge
0040ca84d11f8cd117564cf8d0bbf4a8d12312e54080 refs/pull/179/head
00414a81e5085be4bc1d97a47edd8820abc9cd09e029 refs/pull/179/merge
003f0999cc1115c1b15541cf0450bd39666eee4369a0 refs/pull/18/head
00401220188ba21f7db9278c1c39208c6cf14da7562f refs/pull/18/merge
00407d341d37d91c4801dc36982927440b1599e03f57 refs/pull/181/head
0041b6c519788a73f9d09e95ca6bfe4cbf0465547cb7 refs/pull/181/merge
0040d3fbf6d997b0cefc1c847f39d781057c2a9e1701 refs/pull/182/head
00416a997f16d1437c1e26c6ecf06199bfbfa685ad26 refs/pull/182/merge
004024b34d0a1e6d2c4b78b8e0567a62964f634e90f7 refs/pull/183/head
00412374186358e0d899e28f983f5c5bca2fe633c07a refs/pull/183/merge
0040cc722a413f836f80a35629070e9010278caf7b7f refs/pull/184/head
00412a9ea89ad1fc747adb8033fee75821f76fe08d2e refs/pull/184/merge
0040744ac68d144f31cd23004a930012700d628bce97 refs/pull/186/head
0041e56af9381b28b41b3c5cb4880de6601e92178a8c refs/pull/186/merge
00409825976e06bf3066d3dcaf391ae929ce2377c012 refs/pull/187/head
00410a874afd6694b2c0a58b9c5f4afd0096a564f1e0 refs/pull/187/merge
0040da01fe602774341353b8eebeff0f1c845dc80eb3 refs/pull/189/head
0041bae8c5fa2b7edf56a223a2481c9bb40bbf1f1f82 refs/pull/189/merge
003ff1a1f9cac7ebf5c55ad75b403038eecff633ed64 refs/pull/19/head
0040cecf34933b15cb6dcc7bb8a16ad3a1c0d265b1d9 refs/pull/19/merge
00401ff7777e01a950d9c75398812441f735d038f922 refs/pull/190/head
0041f14ae570b6e47e3884006b9c0ad511369c6f8138 refs/pull/190/merge
00406fb1d8a3775371ba3d32142a7521eef538331f15 refs/pull/191/head
0041f5938c3eafdcffb2351dd7ae847c4aa841fe0dc1 refs/pull/191/merge
0040266519c8d253f9bfb25f6296fbcdd411ba856ba0 refs/pull/192/head
0041b49f0227025bb8a10dd348a1fa83d75e98255aed refs/pull/192/merge
00404bf57537b5664c3bfb85146f4f04e1ff6996cb84 refs/pull/193/head
0041bf0a7533bdc50943dbb31b3b10ea8b5effcb7d54 refs/pull/193/merge
0040248992e380976885e9d49f1f45a9651601cbac44 refs/pull/194/head
0041ba927c86927601e55768b461c611e44c47cc6b4a refs/pull/194/merge
004098eba608fca99a064b24bc74b6ad250e0eeb8d92 refs/pull/197/head
004102f064e7836e8a8481c97c32fea219c4d0978404 refs/pull/197/merge
0040fdeeeeccebf14dc59d9990e9f0ce91315160d63a refs/pull/198/head
0041fe8188b00714a50d9f099c66f46e6cf99d1b9140 refs/pull/198/merge
00405a7ba24790b1ce61f53053b6ce47a4ee0649ca62 refs/pull/199/head
0041fa0c2f28d2f31bd54677fac339b667418386632f refs/pull/199/merge
003e67a06ecca5ae13790e68f2de7ce018f8d0101938 refs/pull/2/head
003f47151c242c1397f811f3f9d55f3d7461ff150ba4 refs/pull/2/merge
003ffd1d908e79ffc4ea9a1e45fd90406f091d65b2ca refs/pull/20/head
004062e5ce95e6a4be69b9c3e2599e89bba0526a8e07 refs/pull/20/merge
0040d64cf5f64c459d839e00482c4d4ada26e1186e34 refs/pull/201/head
0041ac6317949e052e80b71587f0ec50aef4b9f54588 refs/pull/201/merge
00403be4a751350ef9b46cf491a2151af3255b86bb8a refs/pull/204/head
004110e362eb6fe08798a43507c9908befda62ddd67b refs/pull/204/merge
0040682410c66109aaf204eeae38291cf038b49039ee refs/pull/205/head
0041aba58db50a9f1ae0eaa84abcf21f0ea10f04dfb1 refs/pull/205/merge
004041b9051686a946fc1f6255249479cc0e6d5f7940 refs/pull/206/head
00414c7232c20f781b88baa0f0adaf7b7ddce855e317 refs/pull/206/merge
004005202c9ce92106310c5ca385116a94ec068e9dfd refs/pull/207/head
0041faa329c91d17cd5bdac5c1292189c20ffb129ac0 refs/pull/207/merge
004052d5773cc637e6f9524c0ae688a3ad16c3e067a0 refs/pull/208/head
00413da0a63cee556995c5b9d805d0d30775d131af0f refs/pull/208/merge
003fdcfd6f1e07ba91295507cf986d8747a0b90d0c2f refs/pull/21/head
0040c71ab362ce9bec2cb0227335b9d3464e03dd970a refs/pull/21/merge
0040af30cb87251d2a2b7d50a12eabdb7ede20b74b05 refs/pull/210/head
0041722b682d4a38b51d9f28473a180274f023968bd7 refs/pull/210/merge
0040255e14a5c416b06a84f6a4dd76f189e867ba01ea refs/pull/212/head
0041dcaadd5d1f6e4eb59a26e0558b89a84ea090e843 refs/pull/212/merge
0040cb323e95ff01f34fc0a453544265406e23ada8af refs/pull/213/head
00417a4968951d9c3405cf90dcb91cd1806a5576e576 refs/pull/213/merge
004040c520ca1b5b0e7b3ad942d1baf770e8098e06dd refs/pull/214/head
00418bd68f27c9f21e5eddbd1602e1829c0fc712005e refs/pull/214/merge
0040e1186bbead85e1798c97872b9c44ce2abe42fbad refs/pull/215/head
00417c4ef1f4f71e8ac9ce1470dddb5fbf493c2b13cd refs/pull/215/merge
0040beac6d85892ddbc70c39eddd72453c87f18f08fd refs/pull/217/head
0041090435a50042634611f2764e49215670fa2d70e7 refs/pull/217/merge
0040dc59bd8d77610070cd45552c33b1b8de1fc40a53 refs/pull/218/head
00418dfe766e08f81a599ddc803fcfd74ed681501f01 refs/pull/218/merge
003fbac9bf59cf0e14b366867857c7e7119fbab4106b refs/pull/22/head
004009993281f5918483d8f69245e2361ca71fb39d33 refs/pull/22/merge
0040250ab37ce1d3d5665b1b1ab2811f96771fa6ce17 refs/pull/221/head
0041b72574cd740393b9da674689c35004bc473fc3cf refs/pull/221/merge
0040b0793e2dd9de1ea4b99e4cf01d371a8859af5a2c refs/pull/223/head
004177fc5fb65568a8214f90c022404a8c7996d226a4 refs/pull/223/merge
0040bd893986b26729960752e2551c3e526951d936ac refs/pull/224/head
004122f79b0ea04310c0c013cb676762aef081338006 refs/pull/224/merge
0040d46a956a336298846ffe3ab761febc129e9f4253 refs/pull/226/head
0041f35554189be859691db82d8fc9ffbbc01861a28e refs/pull/226/merge
00407565313290edff9fc4b870b44bb7398df3e0f006 refs/pull/227/head
0041721e30804a23ca432f028192d6d30edf820f26ea refs/pull/227/merge
0040ef74464aea88a7a6574595d2511f037f6dcbecb7 refs/pull/228/head
0041999637d8a1cde3b836abc5369463543aecd6a7a2 refs/pull/228/merge
0040a0007dbeab58983a8b692e243b11de40a33f8650 refs/pull/229/head
0041b441e79da93f9a60fa66d43736f2057a7dc1a13d refs/pull/229/merge
003f9da955ca755fd2c970659a3837de964a968b119b refs/pull/23/head
00406e50c384f90fa285d1c2105a6eeda54e943f0445 refs/pull/23/merge
0040c7e89e041061596fb5221aa76e1df7672adc9da2 refs/pull/230/head
0040a6a32a592d4fd8dae18f431f28c3b3914f2e6cf4 refs/pull/231/head
00418d28214eb54c7cdf9303b9b9582668b1e9677228 refs/pull/231/merge
004012636ce89cd275ee865b7dfb0cf69639eb9b535b refs/pull/232/head
0041559994db53d3f7e2635934f590974a2abe33089e refs/pull/232/merge
00401321c63f3b70368ff75b358d01c649172a72f368 refs/pull/233/head
00416a47e0bbbd1993bdfe37167079d360867b3b5aee refs/pull/233/merge
0040e7598075ac85512c419bb0c27e055799f6bbb2e9 refs/pull/234/head
0041d7da1259b9cf45bc860209a4b4deed5e8d839519 refs/pull/234/merge
00402b6c62834222a1b1c870047d999300505d7582ce refs/pull/235/head
0041df34de31ff943087fccf94ba115e0bd30017eb56 refs/pull/235/merge
0040fb7a91739ad9d7764d4959ecd797293a810d1731 refs/pull/236/head
004185e744a171fb25b75a8dad517a94d36e520f729f refs/pull/236/merge
004089a2c61de203546d9b7e50f94b56f580cceeb3aa refs/pull/237/head
0041d5081105c8d51f1d5234ab64b298fec24bcc8054 refs/pull/237/merge
0040e680f28c2f02ce819f2d6e7f1a4f499aa14d8f86 refs/pull/238/head
00416c36657dfc89660f6b8aa2f64d33cf1b2b013a77 refs/pull/238/merge
0040cf8d99689dde86d68eba986e7cd2f50cfb30342a refs/pull/239/head
0041b545cfc4ceeec19a1b4ffe38da0bcad18f2f65bc refs/pull/239/merge
003ff6b46fb0f21f724e0c7a4a076e3a51e1452f95be refs/pull/24/head
00407ac0ca8f09e9d9d3b743a4ccf45aa7b50acf9978 refs/pull/24/merge
004053e9e1966f3eae79ecb9e9069bc47f6ab283032b refs/pull/240/head
00416f8290f1153a53c2890cee15255538aa012a0e57 refs/pull/240/merge
00404c3509ac09c00e5c511c7bc1cc21ab9051a265ea refs/pull/241/head
00418882cef1f805fab40edac7a212ac1223e42383a2 refs/pull/241/merge
00407670c85d70f87d1749d6a49a202c5067a55d2748 refs/pull/242/head
0040b4d311d6a1a04c70803c04f228c99946a02ec12b refs/pull/243/head
004189a07d417a923512659e1eb905fa2f886a58efdc refs/pull/243/merge
00404774466a75c31ae7677b1800a20beb80ae210742 refs/pull/244/head
0040c39f0963ae88a5f8be74a7aee645549386419ed3 refs/pull/245/head
00407a54d4cf6fd5ce26fee7f0b932b6d00ee0692e54 refs/pull/246/head
0041828f058055539d5f3aa8672cb7e7106ccefa66c8 refs/pull/246/merge
004087b2a0d9c5e2b95d2cf2d409338aeb1f91953ccf refs/pull/247/head
00404fb3a01f25480030d3dc413993f41fc8dbd5c394 refs/pull/248/head
0041bbdec4869464e420de78e6521e475df16e6767b1 refs/pull/248/merge
00404135f425d30ce06503c98d02fec3e141f20104ea refs/pull/249/head
003f5f51c728c695c89739c14261d183fc77f3fddf00 refs/pull/25/head
0040d5393f83551a270e08786c158be85b49832c1963 refs/pull/25/merge
0040a6f706926438bbbfaf539cee323d89d5d7d66fe0 refs/pull/250/head
004062b890839a4f3d21685ae0d4b216c8185787f5cf refs/pull/251/head
004055f904556133406af186328530471a92de39902b refs/pull/253/head
00406684186033661d6cdf62bb1e801fab5354e25887 refs/pull/255/head
0040b852a746cebdefc4ad96629fb7de2ce9aa243040 refs/pull/256/head
00401119a86cc337e686514ea87d9a0c34beeb1dcabe refs/pull/258/head
00402cafdcd657bbe974f8c583b385509b32a9229c69 refs/pull/259/head
003fa2c52a7bcb4099fc1b9819ca6574cc3e5c40b28b refs/pull/26/head
004071f78486a038d48f89f73056848646d648386f72 refs/pull/26/merge
004023119d1672e2d65a761b7e43517db433d91e8bea refs/pull/260/head
00408fb94e5aaf8965e96a4e2d6b415ee748dd32519e refs/pull/263/head
00419ec50194d4f530728cd3014cca32132bc7690ab5 refs/pull/263/merge
00406803d077b99b270488932d7df79eeb516f0c031a refs/pull/264/head
004091cad59adabceeed1ac5ce8f6bf3492d6849dd21 refs/pull/265/head
004041e7266ff7b9b5c5ca8f676c2b22d32732f0098d refs/pull/266/head
004063456b5c4be1160e6e59840bc45fcbb937f89999 refs/pull/267/head
004053b638558b7c505a292ae8f60b2231a4a8bb7db4 refs/pull/268/head
00415247d760f379b63712457ae3e1b8df9a99b628fc refs/pull/268/merge
003fd2382c232e60689af4bb4c632d809f8dec5ba7fd refs/pull/27/head
004047ac845e59aa5dc894af07540bcaaf0aabf4bc75 refs/pull/27/merge
004032c2f89fe64167ee2b7cfb7e4c5c0a0729caef21 refs/pull/270/head
0040e254b9cec654bb066b730a5f47a15522eb203002 refs/pull/272/head
00401fb3799118ba2669e478b9a312b50539e512f60d refs/pull/275/head
00407b60f8bdc378d99d80aecfdb84e0e1877cab4a2f refs/pull/276/head
0041dc9819c2ada8ee355e4aeb45766c60d7b8c1d135 refs/pull/276/merge
0040aa9ae329989c8a018cc5f2dcc967ba44419dd694 refs/pull/277/head
003f88d042b7a6bb9b9599757e8e9c88225b4cd59ff7 refs/pull/28/head
004067940c5c5f43972c98753ec91ec19cefa7491b9f refs/pull/28/merge
003fd17fec82341fc83e0d33928740007e604ee1667e refs/pull/29/head
00409865abf26a67c8d0bf81c845c91b58137c37815a refs/pull/29/merge
003e55146c3ece3da1947a66f54b1ee81f7a81de42ed refs/pull/3/head
003fb63e9c90d3f163b2118da6efe1dc4fdf6c90060a refs/pull/3/merge
003fdccc8b6193e79923d7dd76466b48ec48642fc56d refs/pull/30/head
0040bdac2824ec6d4fd7ae6e3792ff4d4c717ca939b2 refs/pull/30/merge
003f6c6f94fedd6dc299fa826482e8d525a5815d4d3b refs/pull/31/head
00407a82460b73d6415d9917f9c53c37c085184a2e59 refs/pull/31/merge
003fa6119d94990ccffa228f6f83b773985ddbe77467 refs/pull/32/head
0040c6bf80e9244fb02d5b1e82876666f8ddc46ae39a refs/pull/32/merge
003fd1158f5a11c9d6f4790062c4f7065ee75dc19cd8 refs/pull/33/head
0040003bf6559c2a97c80e070b658154103c5b62b12b refs/pull/33/merge
003f7d876fc4147ff2c103f75f31e1baef7bdb69ada1 refs/pull/34/head
00400742863235189a6657849007075fcf2736c80fc0 refs/pull/34/merge
003ff26416d6dd617773b461830b319c24d5e4b50c65 refs/pull/35/head
0040994ddcc2f08d362fc70d521eb0ccab74481deb51 refs/pull/35/merge
003f8ad7a23648b9908bc5f2d2634fa49758a3728646 refs/pull/36/head
0040553b6a856ee8e65cf7aeb7271f543445b538e572 refs/pull/36/merge
003f066cf2b811790ccc0034d5c4a305f90c3ec3756e refs/pull/37/head
0040a57c1ecfbee22013776f60044d04d6f2edc01e12 refs/pull/37/merge
003f69f8b8693dd7d9213f13e5ee54987234f5c91181 refs/pull/38/head
00406a1c0fb627661ae17a2a01398698fc5d2790f9fb refs/pull/38/merge
003f89d13993049aaa9b5d79c967f027ef954978148c refs/pull/39/head
0040d16f2287c0bb5d409c58f9b2e60287853d9315cc refs/pull/39/merge
003e2589a76d664e943d2dbfa11af950c1346c7de1b4 refs/pull/4/head
003f9cf603d559455761ef60524165c2e540fd210397 refs/pull/4/merge
003f3381177341c11aa9d89b6172699d64ddc2146a11 refs/pull/40/head
00405e08b717c210ae995abafad3f223c6a18c5d8a4c refs/pull/40/merge
003faa3f0b3a2aeb1e6af36fd0d73dca991a69cdb339 refs/pull/41/head
0040b94395eec1ba2735f32628609ab7dd50c6619141 refs/pull/41/merge
003f38b9758a5b8f2375744b36f91a32e0d0fd4ee54b refs/pull/42/head
004002b10430bd0dae6f9807393700f6d76654be3f6d refs/pull/42/merge
003f72c1a6135d191b5b72a5e882c5cf2f3d51a6d33d refs/pull/43/head
00406681bf7196ad794b880612272174dba43182ddfc refs/pull/43/merge
003f93750397422effd15fdef8f378ceb6ca18f434ec refs/pull/44/head
004025cd8b21dd58495a104d0fc002cd30c3efbde711 refs/pull/44/merge
003f13af54fdd06ddc5ad2d7e2ddccb199c31e271fd1 refs/pull/45/head
00406fbe724e8d46c53b8978d8536d044111e854ac79 refs/pull/45/merge
003f36839921833938603bb6cafe3162ced63ad60b7b refs/pull/46/head
00407ee3ac2bc07b65d535ced4664c53689624e308ea refs/pull/46/merge
003f8f551e3dc164c39d085b5d2c19c6623b873fc77e refs/pull/47/head
0040ec93c65f2b6f2b0d7069f43ffb6c620a701733a9 refs/pull/47/merge
003e1e9307e88aa468fe853f319151572bc31b4aae7e refs/pull/5/head
003ff71d93d6df1ca72fe51110ebc443a8ac630d982b refs/pull/5/merge
003f0d0b59b254c844141d539cb7ae9a10b2c91d8a1d refs/pull/51/head
0040f2088a0f62ce7129ac8e511a7218c3205eea2610 refs/pull/51/merge
003f60ae2e8c892bdbaf4f7e2d6ae379e9bbbb8ca093 refs/pull/53/head
0040802b4b3ad3b03521b1fb9a8912352e4c641f2095 refs/pull/53/merge
003f026f93cd837e70a74e020233630e7669ca890bab refs/pull/54/head
0040e0e8450735442d2993a9a2d2f6bd6c22bc55cb29 refs/pull/54/merge
003fe5aa6256f06c8cd5dcd82fe66dc3e49f2eda9bd1 refs/pull/55/head
004091a314e6dc96999f89a9741a6786a20625cd4c12 refs/pull/55/merge
003f5f6d331e6ab653da376af4743af562f36fd49a30 refs/pull/56/head
0040968468d1d7781034344a393e203c9903e944e66b refs/pull/56/merge
003fe2a1662318483fb225149f534425c92daec51b94 refs/pull/57/head
004065219c983b1af58b11211bb89ce321bbc9e88f62 refs/pull/57/merge
003f253765d81eefe79d312bea176a72744e8d577b67 refs/pull/58/head
00406f324f84bf95788bc4ce58f28bc73467248fb7bb refs/pull/58/merge
003fa7deba0f905ff0a89421c83958046cadc3404225 refs/pull/59/head
0040a3bbe6400670d4380efce095852fb429a4720a9e refs/pull/59/merge
003e5e06b2b9fa58969febe5d1a20049cd63d6d06376 refs/pull/6/head
003f6271a7cca60b788a33f69e9291ba0b3481378f4e refs/pull/6/merge
003fab9c0448c984becbc75cfdbb34b616222307c2b3 refs/pull/60/head
0040284284cd3b9e80319ecb91afaa405edd0e8252bb refs/pull/60/merge
003fb6378dae511dd21263e629ca9559b10ccb9e59b8 refs/pull/61/head
00400197e3136a66f0ce037791475e97f4a651d0697a refs/pull/61/merge
003f51ae55edffa0e6f652f8b8cdf65bc895599fed23 refs/pull/62/head
0040615b026e922093230fdc667fcabe9935826f6ac5 refs/pull/62/merge
003fba7ddbc015810cfb821e788c13a7f6c5c84708fc refs/pull/63/head
00406d5b2d09842061febaf63617ad45eb3e1a52c138 refs/pull/63/merge
003fa15b4bb6874e110761301cc3bce5754cae53e6e4 refs/pull/64/head
0040883bfce8df9d0cfa2a1bd13f3cf995c3f0d2c85d refs/pull/64/merge
003fa5b2daaf1bdf58a9a0ed9580656fe0b58b9034da refs/pull/65/head
0040d53b4508fc255390dd5fc1e0ed1a29200fe7d326 refs/pull/65/merge
003fc39f7712f740aee0f4af905e73699a17b5eab378 refs/pull/66/head
004032b09d0bd01f936327f4787db0a6502de131c561 refs/pull/66/merge
003f7247d9b3eed526f8af226a0046877a5d32348476 refs/pull/68/head
0040df8d49a26cede559d756163d9e20a485f111b552 refs/pull/68/merge
003ee0ae3e245a421d02e3519db3776c7166d03ba799 refs/pull/7/head
003f1886d2cc2f79e0f23946ea8769c6a7b7e4adb0da refs/pull/7/merge
003f9d488f55296710347d0b2081e914c196b4908e3c refs/pull/70/head
0040e3b55172abc556c5863495fe99ac1a0f3595d155 refs/pull/70/merge
003f5000e1de4a9fc4dc865231e3a26d414b85e4f4bd refs/pull/71/head
0040ee1189786a3a7abf7cf43333aec4595926d64471 refs/pull/71/merge
003ffd58cf8975c28e940efd4642fb02de5fca9a0be4 refs/pull/72/head
004054152cfa33f8e8d98169a208d0943e32648282fd refs/pull/72/merge
003fe0b87719503dfaf94d87162a79d919776240a9f5 refs/pull/73/head
00405e13384eed2f086b23411201f7d6aedd9f06bddb refs/pull/73/merge
003f384d79d67184cbab86bad714936c65578cfc4540 refs/pull/74/head
00403d88cc65793e9c39412ab61c5b0604247d884426 refs/pull/74/merge
003f22ef9161c2ab3728b4b63238dbd68097b2a98bb4 refs/pull/75/head
0040436084d3304d39aa18204b60f08b2e2d30657996 refs/pull/75/merge
003f428695f732d07fb95ce4807949112eb5f0285423 refs/pull/77/head
0040a8133b8ebce27b5f39b1ccffcb5f3f9902886f45 refs/pull/77/merge
003fce3c55ba3f66ece7af529057dae7141747f7c186 refs/pull/78/head
0040616ab475bfae337bb8d163484e181e09af6ec5bb refs/pull/78/merge
003efad972304cb9c6a8051df6d7b065e5c051963fce refs/pull/8/head
003f61262ebac2c4fd336e35118705a24c7366257d0e refs/pull/8/merge
003f60c71a98a4988ab811ee506c4f80b043192604ef refs/pull/80/head
00409d81a0492337146fe7a7b92e6e77947795179a5f refs/pull/80/merge
003f6610fc39cc1a1916bfa58b63e3836f51ca4c6816 refs/pull/81/head
0040926a0eb61bd6bb5c8840431ba6245c712cf6e81c refs/pull/81/merge
003f466229cf4d7f80f4a3ca5afbf676ee033636910b refs/pull/82/head
004001765f88af5086f8fca37d4cbd15c1436cbd2673 refs/pull/82/merge
003f03a2d608c29929591445edb5f25c365280352a54 refs/pull/84/head
004035e6f0bb5e50612f27d1fe43d4607eaab0a03b06 refs/pull/84/merge
003fa86e3aa7d9ce13002255fd53d3f960aca8dda6da refs/pull/85/head
00401cfb67355d8dd2d902f96a6b4c9a8fe00206df7f refs/pull/85/merge
003f58e9e0c5578a1919bf1f055aab561657ac857496 refs/pull/86/head
0040a191e27375b2229e4e21ba8cf3a939421e47145a refs/pull/86/merge
003fbe85442e5e9b834d1ce9f2599266c0ff830ad9f2 refs/pull/87/head
0040f5dee680ec9056cfbf2015af3105c616151b7668 refs/pull/87/merge
003fd8a05f1766c4efb770afa101b0493533b18c4ef4 refs/pull/88/head
00406d6f0a181d9d4384a95da824e99ce37aa00140db refs/pull/88/merge
003f08abb4bb6a6a7b3d0416d2ca8ed2b576675f9219 refs/pull/89/head
0040db4373208bfbe5a0abc7c1e7bec4abb0ab279637 refs/pull/89/merge
003e7c1c56deb3de9572e9cf3b5c78f17fb8ccc360be refs/pull/9/head
003fff08b8247e91028214b2755d361ca6f1f6c12760 refs/pull/9/merge
003f5a88da1d3742e52c1cbc97b3a4615a06b7fb2c02 refs/pull/92/head
004089e3f742fd5f0cbacca0215e64fdc1959401ba74 refs/pull/92/merge
003fb3654e68d986a05674d7cd7c43c7f0435b96670d refs/pull/93/head
00403eb50f2132763277cc484798c6dcd09413249ae6 refs/pull/93/merge
003fe8a284d2954afca787a8a57e68d112a23d741dd0 refs/pull/94/head
0040b6229718f37ec77647c2d249142e9d38aafab991 refs/pull/94/merge
003f1bf4e656a876bdeb1b6a0190a9ce48557a47c8dd refs/pull/98/head
00408db7fad7e6763f320bbdafe12bb666acccf8c319 refs/pull/98/merge
003920ca21a3f7122cf7caa91cb0e9b9c69be9279950 refs/tags/0
003e2403fe79c1372329d63a51e54e37cec677acca27 refs/tags/v0.1.0
003e7b289043c7beced434be4334fb909ba0b16b57b1 refs/tags/v0.1.1
003e5589b6faabc822255c87b096c57afaef9fa47d6f refs/tags/v0.1.2
0041e77b9aa020a2041a0f88459a4d82f236517dff09 refs/tags/v0.1.2^{}
0042e2e035cac84c1a971df68bda66ae15ac84f367bb refs/tags/v0.2.0-rc0
0045088a01f19cda4818716c5a2b6a216752ca8825e4 refs/tags/v0.2.0-rc0^{}
0000`,
	}, "\n")

	g, err := parseGitUploadPack(ioutil.NopCloser(strings.NewReader(in)))

	if err != nil {
		t.Fatalf("Got unexpected error: %v", err)
	}

	if g.capabilities != "multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done agent=git/1.8.4" {
		t.Errorf(
			"Expected capabilities to be \n\t%v\nGot:\n\t%v",
			"multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done agent=git/1.8.4",
			g.capabilities,
		)
	}

	if len(g.refs) != 401 {
		t.Errorf("Expected 401 refs, got %v", len(g.refs))
	}

}
