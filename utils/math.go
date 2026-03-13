package utils

import (
	"fmt"
	"math/big"
	"strings"
)

func StringToDivideBy1e18(str string) (string, error) {
	str = strings.ReplaceAll(str, ",", "")

	num := new(big.Float)
	_, success := num.SetString(str)
	if !success {
		return "", fmt.Errorf("invalid number string: %s", str)
	}

	divisor := new(big.Float).SetFloat64(1e18)
	result := new(big.Float).Quo(num, divisor)

	// 先格式化为固定小数位的字符串
	formatted := result.Text('f', 18)

	// 去除末尾的零和小数点
	trimmed := strings.TrimRight(formatted, "0")
	trimmed = strings.TrimRight(trimmed, ".")

	return trimmed, nil
}

// BigIntDivideBy1e18ToUint64 将 *big.Int 除以 1e18 后转换为 uint64
// 注意：除法结果的小数部分会被截断
func BigIntDivideBy1e18ToUint64(num *big.Int) (uint64, error) {
	if num == nil {
		return 0, fmt.Errorf("input big.Int is nil")
	}

	// 检查是否为负数
	if num.Sign() < 0 {
		return 0, fmt.Errorf("negative number cannot be converted to uint64")
	}

	// 转换为 big.Float 进行除法运算
	floatNum := new(big.Float).SetInt(num)
	divisor := new(big.Float).SetFloat64(1e18)
	result := new(big.Float).Quo(floatNum, divisor)

	// 转换为 uint64
	uint64Result, accuracy := result.Uint64()
	if accuracy == big.Above {
		return 0, fmt.Errorf("number exceeds uint64 max value")
	}

	return uint64Result, nil
}

// BigIntDivideBy1e18ToUint64WithRound 将 *big.Int 除以 1e18 后转换为 uint64
// 使用四舍五入处理小数部分
func BigIntDivideBy1e18ToUint64WithRound(num *big.Int) (uint64, error) {
	if num == nil {
		return 0, fmt.Errorf("input big.Int is nil")
	}

	// 检查是否为负数
	if num.Sign() < 0 {
		return 0, fmt.Errorf("negative number cannot be converted to uint64")
	}

	// 转换为 big.Float 进行除法运算
	floatNum := new(big.Float).SetInt(num)
	divisor := new(big.Float).SetFloat64(1e18)
	result := new(big.Float).Quo(floatNum, divisor)

	// 加 0.5 后取整实现四舍五入
	half := new(big.Float).SetFloat64(0.5)
	rounded := new(big.Float).Add(result, half)

	// 转换为 uint64
	uint64Result, accuracy := rounded.Uint64()
	if accuracy == big.Above {
		return 0, fmt.Errorf("number exceeds uint64 max value")
	}

	return uint64Result, nil
}

// BigIntToUint64 直接将 *big.Int 转换为 uint64（不进行除法操作）
func BigIntToUint64(num *big.Int) (uint64, error) {
	if num == nil {
		return 0, fmt.Errorf("input big.Int is nil")
	}

	// 检查是否为负数
	if num.Sign() < 0 {
		return 0, fmt.Errorf("negative number cannot be converted to uint64")
	}

	// 检查是否超出 uint64 范围
	maxUint64 := new(big.Int).SetUint64(^uint64(0))
	if num.Cmp(maxUint64) > 0 {
		return 0, fmt.Errorf("number exceeds uint64 max value")
	}

	return num.Uint64(), nil
}

func StringToDivideBy1e10(str string) (string, error) {
	str = strings.ReplaceAll(str, ",", "")

	num := new(big.Float)
	_, success := num.SetString(str)
	if !success {
		return "", fmt.Errorf("invalid number string: %s", str)
	}

	divisor := new(big.Float).SetFloat64(1e10)
	result := new(big.Float).Quo(num, divisor)

	// 先格式化为固定小数位的字符串
	formatted := result.Text('f', 10)

	// 去除末尾的零和小数点
	trimmed := strings.TrimRight(formatted, "0")
	trimmed = strings.TrimRight(trimmed, ".")

	return trimmed, nil
}
