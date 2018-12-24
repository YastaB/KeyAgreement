package main

import (
	"fmt"
	"math/big"
)

// participants perform key agreement protocol
// initialization parameters are defined before

/**
p and q are large prime numbers
g = (Zp)
x = (Zq)
y = g^x mod p
k,v = (Zq) (are randomly selected)
r = (Zq) (randomly selected)
**/
var q, _ = new(big.Int).SetString("70830437809992052122705382232092710520096640888343042864124381269620926369090824665254595835832924535603173757631672116331232968683454751265258387352641382374279467305216459307052164504547904309385274902216434059305334668453580540584548701719844965116691799027770171599194704612669102357388906362282730675783", 10)
var p, _ = new(big.Int).SetString("141660875619984104245410764464185421040193281776686085728248762539241852738181649330509191671665849071206347515263344232662465937366909502530516774705282764748558934610432918614104329009095808618770549804432868118610669336907161081169097403439689930233383598055540343198389409225338204714777812724565461351567", 10)
var g, _ = getGenerator(p)

type participant struct {
	// inited values
	privateKey PkSigSchnorrPrivate
	publicKey  PkSigSchorrPublic
	t          *big.Int
	v          *big.Int
	tInv       *big.Int
	vInv       *big.Int
	// First section values
	omega *big.Int
	A     *big.Int
	B     PkSigSchorrSign
	// Second section values
	a   *big.Int
	ckI *big.Int
	cI  *big.Int
	kI  *big.Int
	// Last Step = Calculated group key
	groupKey *big.Int

	// signature struct
	sig PkSigSchnorr
}

var participants [3]participant

const participantCount = 3

func printParticipant(node participant) {
	fmt.Println("Private Key :", node.privateKey, " Public Key:", node.publicKey, " t:", node.t,
		" v:", node.v, " tInv:", node.tInv, " vInv:", node.vInv)
}

func initParticipant() participant {
	// initial variables
	sig := newPkSigSchnorr(p, q)
	privateKey, publicKey := sig.keyGen()
	t, _ := generateRandom(q)
	v, _ := generateRandom(q)
	tInv := new(big.Int).ModInverse(t, q)
	vInv := new(big.Int).ModInverse(v, q)
	return participant{privateKey: privateKey, publicKey: publicKey, t: t, v: v, tInv: tInv, vInv: vInv}
}

func findNeigbour(index, max int) (int, int) {
	var before int
	var after int
	if index-1 < 0 {
		before = max - 1
	} else {
		before = index - 1
	}
	after = (index + 1) % max
	return before, after
}

func calcTempPubParams(node participant) participant {
	node.omega = new(big.Int).Exp(g, node.v, p)
	node.A = new(big.Int).Exp(g, node.v, p)
	node.sig = newPkSigSchnorr(p, q)
	node.privateKey, node.publicKey = node.sig.keyGen()
	node.B = node.sig.sign(node.publicKey, node.privateKey, node.omega)
	fmt.Println("omega: ", node.omega, " A: ", node.A)
	return node
}

func calcTempSecretKeys(node participant, index int) participant {
	_, after := findNeigbour(index, participantCount)

	node.a, _ = generateRandom(q)
	node.ckI = new(big.Int).Exp(participants[after].omega, node.v, p)
	/*
		Signing will be implemented
	*/
	node.cI = new(big.Int).Exp(g, node.a, p)
	exp := new(big.Int).Exp(node.ckI, node.a, p)
	node.kI = exp.Mod(exp, q)
	return node
}

func calculateKey(node participant, index int) *big.Int {
	resultingKey := new(big.Int).Mod(node.ckI, p)
	for i := 0; i < len(participants); i++ {
		if i != index {
			fmt.Println("Calculating key in of ", i)
			mulled := new(big.Int).Mul(resultingKey, participants[i].ckI)
			resultingKey = mulled.Mod(mulled, p)
		}
	}
	return resultingKey
}

func verifyTempSecretKeys(nodeWi, index int) bool {

	/*
	 */
	return true
}

func verifyPubVariables(nodeWillVerify participant, nodePublic PkSigSchorrPublic, sign PkSigSchorrSign, M *big.Int) bool {
	return nodeWillVerify.sig.verify(nodePublic, sign, M)
}

func main() {

	node1 := initParticipant()
	node2 := initParticipant()
	node3 := initParticipant()

	participants[0] = node1
	participants[1] = node2
	participants[2] = node3

	// printout public variables
	fmt.Println("Public Variables")
	fmt.Println("q:", q, " p:", p, " g:", g)

	for i := 0; i < len(participants); i++ {
		fmt.Println("test", i)
		participants[i] = calcTempPubParams(participants[i])
	}

	for i := 0; i < participantCount; i++ {
		for j := 0; j < participantCount; j++ {
			if i != j {
				fmt.Println("")
				verifyRes := verifyPubVariables(participants[i], participants[j].publicKey, participants[j].B, participants[j].omega)
				fmt.Println("Node ", i, " verifies Node ", j, " Res: ", verifyRes)
				participants[i] = calcTempPubParams(participants[i])
			}
		}
	}

	//verifyRes1 := verifyPubVariables(float64(node1.w), float64(node1.A),
	//	float64(node1.B), float64(node1.y))

	for i := 0; i < len(participants); i++ {
		participants[i] = calcTempSecretKeys(participants[i], i)
		//printParticipant(participants[i])
		fmt.Println(i, " ckI  ", participants[i].ckI)
	}

	for i := 0; i < len(participants); i++ {
		fmt.Println("Node ", i, " Key: ", calculateKey(participants[i], i))
	}

	/*
		//sigYFalse, _ := new(big.Int).SetString("1231123", 10)
		sig := newPkSigSchnorr(p, q)
		privateKey, publicKey := sig.keyGen()
		sigOmega, _ := new(big.Int).SetString("123123", 10)

		sign := sig.sign(publicKey, privateKey, sigOmega)
		res := sig.verify(publicKey, sign, sigOmega)
		fmt.Println("Res: ", res)
	*/
}
