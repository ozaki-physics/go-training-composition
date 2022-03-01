// ファイルの中身を暗号化してファイルに出力するツール
// Use Case
// 5ファイル
//     - 暗号化と復号に使う key ファイル(sample: fileCrypto/key/AES-CTR.key)
//     - 暗号化したい内容が書かれたファイル(sample: fileCrypto/plain/plainText.md)
//     - 暗号化した結果を書き出すファイル(sample: fileCrypto/cipher/cipherText.md)
//     - 復号した結果を書き出すファイル(sample: fileCrypto/decode/decodeText.md)
//     - 上記の4ファイルのパスを書いた JSON ファイル(sample: fileCrypto/json/filepath.json)
//
// RunFileEnCrypto() 暗号化
// RunFileDeCrypto() 復号
//
// 実行するときに コマンド引数で JSON ファイルのパスを渡すか
// 実行したあとに 標準入力で JSON ファイルのパスを渡す
//
// sample
// go run main.go fileCrypto/json/filepath_product.json
package fileCrypto

/*
	計算量問題があるから ここからここまでしか暗号化実施しないとかできないかな
	key の大きさを公表しないようにしたい
*/
