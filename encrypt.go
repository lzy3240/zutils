package zutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/tjfoc/gmsm/sm4"
)

const (
	PWD string = "default@66668888"
)

func EncryptBySM4(src, pwd string) string {
	return base64.StdEncoding.EncodeToString(encryptSM4([]byte(src), []byte(pwd)))
}

func DecryptBySM4(src, pwd string) string {
	tmpRune, _ := base64.StdEncoding.DecodeString(src)
	return string(decryptSM4(tmpRune, []byte(pwd)))
}

func EncryptByAES(src, pwd string) string {
	return base64.StdEncoding.EncodeToString(encryptAES([]byte(src), []byte(pwd)))
}

func DecryptByAES(src, pwd string) string {
	tmpRune, _ := base64.StdEncoding.DecodeString(src)
	return string(decryptAES(tmpRune, []byte(pwd)))
}

//SM4算法
//sm4加密
func encryptSM4(src, pwd []byte) []byte {
	//创建加密块
	block, err := sm4.NewCipher(pwd)
	if nil != err {
		panic(err)
	}
	//填充数据
	src = paddingText(src, block.BlockSize())
	//初始化向量
	iv := []byte("1234567887654321")
	//设置加密模式
	blockmode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(src))
	blockmode.CryptBlocks(dst, src)
	return dst
}

//sm4解密
func decryptSM4(src, pwd []byte) []byte {
	//创建解密块
	block, err := sm4.NewCipher(pwd)
	if nil != err {
		panic(err)
	}
	//初始化向量
	iv := []byte("1234567887654321")
	//创建解密模式
	blockmode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(src))
	//解密
	blockmode.CryptBlocks(dst, src)

	//去除填充
	dst = unPaddingText(dst)
	return dst
}

//给最后一组数据填充至64字节
func paddingText(src []byte, blockSize int) []byte {
	//求出最后一个分组需要填充的字节数
	padding := blockSize - len(src)%blockSize
	//创建新的切片，切片字节数为padding
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	//将新创建的切片和带填充的数据进行拼接
	nextText := append(src, padText...)
	return nextText

}

//取出数据尾部填充的赘余字符
func unPaddingText(src []byte) []byte {
	//获取待处理数据长度
	length := len(src)
	//取出最后一个字符
	num := int(src[length-1])
	newText := src[:length-num]
	return newText
}

//AES算法
//aes加密
func encryptAES(data, key []byte) []byte {
	//key := []byte(sKey)
	//创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//判断加密块的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted
}

//aes解密
func decryptAES(data, key []byte) []byte {
	//key := []byte(sKey)
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		panic(err)
	}
	return crypted
}

//pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

//pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}
