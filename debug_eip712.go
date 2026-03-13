package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"golang.org/x/crypto/sha3"
)

func main() {
	// 配置 - 根据你的实际情况
	privateKeyHex := "d712af69c090e5a7f855c9d42561a59aa6461bfae03fd3ac4ff4d9872f2bb461"
	chainID := int64(11155111)
	contractAddress := common.HexToAddress("0xc94df36b1c2e262A8f18088443f1a5a96BD02eeD")

	// 计算私钥对应的地址
	privateKey, _ := crypto.HexToECDSA(privateKeyHex)
	publicKeyECDSA, _ := privateKey.Public().(*ecdsa.PublicKey)
	signerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Println("=== 配置信息 ===")
	fmt.Println("Signer地址:", signerAddress.Hex())
	fmt.Println("合约地址:", contractAddress.Hex())
	fmt.Println("Chain ID:", chainID)

	// 领取参数 - 与你的交易参数匹配
	user := common.HexToAddress("0xF321240E4E2880910f24d67B8C9B9299767988fe")
	amount := uint64(1)
	nonce := uint64(0)
	deadline := uint64(0x69b3873d) // 1773905981

	fmt.Println("\n=== 交易参数 ===")
	fmt.Println("User:", user.Hex())
	fmt.Println("Amount:", amount)
	fmt.Println("Nonce:", nonce)
	fmt.Printf("Deadline: %d (0x%x)\n", deadline, deadline)

	// 手动构建 EIP-712 哈希
	// 1. 构建 domain separator
	fmt.Println("\n=== EIP-712 哈希构建过程 ===")

	// DOMAIN_TYPEHASH = keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)")
	domainTypehash := crypto.Keccak256([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))
	fmt.Printf("DOMAIN_TYPEHASH: 0x%x\n", domainTypehash)

	// keccak256("Airdrop")
	nameHash := crypto.Keccak256([]byte("Airdrop"))
	// keccak256("1")
	versionHash := crypto.Keccak256([]byte("1"))

	// domainSeparator = keccak256(abi.encode(DOMAIN_TYPEHASH, nameHash, versionHash, chainId, verifyingContract))
	domainSeparator := crypto.Keccak256(padBytes(append(append(append(append(
		domainTypehash,
		nameHash...),
		versionHash...),
		padUint256(big.NewInt(chainID))...),
		padAddress(contractAddress)...)))
	fmt.Printf("Domain Separator: 0x%x\n", domainSeparator)

	// 2. 构建 struct hash
	// CLAIM_TYPEHASH = keccak256("ClaimRequest(address user,uint256 amount,uint256 nonce,uint256 deadline)")
	claimTypehash := crypto.Keccak256([]byte("ClaimRequest(address user,uint256 amount,uint256 nonce,uint256 deadline)"))
	fmt.Printf("CLAIM_TYPEHASH: 0x%x\n", claimTypehash)

	// structHash = keccak256(abi.encode(CLAIM_TYPEHASH, user, amount, nonce, deadline))
	structHash := crypto.Keccak256(padBytes(append(append(append(append(
		claimTypehash,
		padAddress(user)...),
		padUint256(new(big.Int).SetUint64(amount))...),
		padUint256(new(big.Int).SetUint64(nonce))...),
		padUint256(new(big.Int).SetUint64(deadline))...)))
	fmt.Printf("Struct Hash: 0x%x\n", structHash)

	// 3. 构建最终哈希
	// hash = "\x19\x01" + domainSeparator + structHash
	finalHash := crypto.Keccak256(append(append([]byte("\x19\x01"), domainSeparator...), structHash...))
	fmt.Printf("Final Hash: 0x%x\n", finalHash)

	// 4. 签名
	sig, _ := crypto.Sign(finalHash, privateKey)
	sig[64] += 27 // 转换为以太坊签名格式
	sigHex := "0x" + common.Bytes2Hex(sig)
	fmt.Printf("\n=== 签名结果 ===\nSignature: %s\n", sigHex)

	// 5. 验证签名
	recoveredAddr, _ := recoverAddress(finalHash, sig)
	fmt.Printf("恢复的地址: %s\n", recoveredAddr.Hex())
	fmt.Printf("签名验证: %v\n", recoveredAddr == signerAddress)

	// 使用 apitypes.TypedData 验证
	fmt.Println("\n=== 使用 apitypes.TypedData 验证 ===")
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"ClaimRequest": {
				{Name: "user", Type: "address"},
				{Name: "amount", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		},
		PrimaryType: "ClaimRequest",
		Domain: apitypes.TypedDataDomain{
			Name:              "Airdrop",
			Version:           "1",
			ChainId:           math.NewHexOrDecimal256(chainID),
			VerifyingContract: contractAddress.Hex(),
		},
		Message: apitypes.TypedDataMessage{
			"user":     user.Hex(),
			"amount":   new(big.Int).SetUint64(amount),
			"nonce":    new(big.Int).SetUint64(nonce),
			"deadline": new(big.Int).SetUint64(deadline),
		},
	}

	tdHash, _, _ := apitypes.TypedDataAndHash(typedData)
	fmt.Printf("TypedData Hash: 0x%x\n", tdHash)
	fmt.Printf("Hashes match: %v\n", hex.EncodeToString(finalHash) == hex.EncodeToString(tdHash))

	// 解析你的签名
	fmt.Println("\n=== 你的签名解析 ===")
	yourSig, _ := hex.DecodeString("2647652258cf66c7518617575b38349f9188a202a96e2622a1bfb074838a441e72fc1eb047dbd660767d46702c83485793d1848acde22f203c38b1d387ada0551c")
	if len(yourSig) == 65 {
		r := new(big.Int).SetBytes(yourSig[:32])
		s := new(big.Int).SetBytes(yourSig[32:64])
		v := int(yourSig[64])
		fmt.Printf("r: 0x%x\n", r)
		fmt.Printf("s: 0x%x\n", s)
		fmt.Printf("v: %d\n", v)

		// 验证你的签名
		recAddr, err := recoverAddress(finalHash, yourSig)
		if err != nil {
			fmt.Println("签名恢复失败:", err)
		} else {
			fmt.Printf("你的签名恢复的地址: %s\n", recAddr.Hex())
			fmt.Printf("与 signer 匹配: %v\n", recAddr == signerAddress)
		}
	}

	fmt.Println("\n=== 重要检查项 ===")
	fmt.Println("1. 合约中的 signer 地址必须设置为:", signerAddress.Hex())
	fmt.Println("2. 调用 claim 的 msg.sender 必须是:", user.Hex())
	fmt.Println("3. 合约地址必须是:", contractAddress.Hex())
}

func padBytes(b []byte) []byte {
	// 确保字节切片长度为32的倍数
	for len(b)%32 != 0 {
		b = append(b, 0)
	}
	return b
}

func padUint256(n *big.Int) []byte {
	b := n.Bytes()
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}

func padAddress(addr common.Address) []byte {
	padded := make([]byte, 32)
	copy(padded[12:], addr.Bytes())
	return padded
}

func recoverAddress(hash []byte, sig []byte) (common.Address, error) {
	sigCopy := make([]byte, 65)
	copy(sigCopy, sig)
	if sigCopy[64] >= 27 {
		sigCopy[64] -= 27
	}
	pubKey, err := crypto.SigToPub(hash, sigCopy)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubKey), nil
}

// Keccak256 计算哈希
func Keccak256(data ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, d := range data {
		h.Write(d)
	}
	return h.Sum(nil)
}
