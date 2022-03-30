# データ圧縮 について
## 言語に依存しない知識としての勉強
Compress: 圧縮  
CoinMarketCap で `Accept-Encoding: deflate, gzip` について書かれていたから調べる  
通信データを圧縮するなら クライアントとサーバが同じ圧縮アルゴリズムを使わないといけない  
クライアントが対応している圧縮方法を伝えるのが リクエストヘッダー の `Accept-Encoding`  
サーバがどの圧縮方法を使ったか レスポンスヘッダー の `Content-Encoding` に含まれる  
サーバが圧縮しないでレスポンスする可能性もある  

```
// 複数のアルゴリズムを quality value で重み付けする構文:
Accept-Encoding: deflate, gzip;q=1.0, *;q=0.5
```

学術的な名称 と Accept-Encoding で表記される名称には少しズレがある?  

### Accept-Encoding の deflate(デフレート)
zlib 構造体と deflate 圧縮アルゴリズムを用いた圧縮形式 と書かれた文献があった  
調べた感じ deflate (RFC 1951) + RFC 1950 のこと言ってるっぽい  
だから実質 zlib と同義  

### Accept-Encoding の gzip
Lempel-Ziv coding (LZ77) と32ビット CRC を用いた圧縮形式 と書かれた文献があった  
調べた感じ deflate (RFC 1951) + RFC 1952 または LZ77 + RFC 1952 っぽい  
なぜ または って言っているかというと  
deflate (RFC 1951) が LZ77 + ハフマン符号化 だから deflate (RFC 1951) = LZ77 とは言えないため  
まぁ LZ77 だけを使っているとは思えないから おそらく deflate (RFC 1951) + RFC 1952 だろう  
つまり Accept-Encoding の deflate との違いは データ形式のみで アルゴリズムは同一となる  

### 学術的な deflate  
LZ77 と ハフマン符号化 を組み合わせた 可逆データ圧縮アルゴリズム  
1996年5月に RFC 1951 としてドキュメント化  
ヘッダーやフッターをつけた zlib (RFC 1950) 形式や gzip (RFC 1952) 形式とともに使われる事が多い  

### LZ77 とは  
開発されたデータ圧縮アルゴリズム であり アルゴリズム名は Lempel-Ziv アルゴリズム  
厳密には ほとんどのケースで LZ77 の改良である LZSS が使われている  

### ハフマン符号 とは  
文字列をはじめとするデータの可逆圧縮などに使用される  
- 静的ハフマン符号 (Static Huffman coding)  
ファイルなど固定長のデータに対し 1度全部読み込んで データのすべての記号を調べて符号木を構築しておき もう1度頭から読み直して符号化を行う方法  
- ダイナミックハフマン符号 (Dynamic Huffman coding)  
データの全部ではなく ブロック毎に符号を作る deflate で利用されている  
- 適応型ハフマン符号 (Adaptive Huffman coding)  
最初の状態では頻度情報を持たず 記号を1個読み込むごとに符号木を作り直す方法  

### 学術的な zlib  
データの圧縮および伸張を行うためのフリーのライブラリ  
可逆圧縮アルゴリズムである deflate (RFC 1951) を実装  
ヘッダーやフッターなどのデータ形式 は RFC 1950  
2バイト以上のヘッダーと4バイトのフッター  

### 学術的な gzip  
データ圧縮プログラムのひとつ および その圧縮データのフォーマット  
多くの UNIX に標準搭載  
Lempel-Ziv アルゴリズム (LZ77) と ハフマン符号 を使用  
ヘッダーやフッターなどのデータ形式 は RFC 1952  
10バイト以上のヘッダーと8バイトのフッター  

### Go 言語での データ圧縮
deflate は 標準ライブラリの [compress/flate](https://pkg.go.dev/compress/flate@go1.17.8) が使える  
注意: http リクエストヘッダー は deflate と言いつつ zlib なので 標準ライブラリの [compress/zlib](https://pkg.go.dev/compress/zlib@go1.17.8) がよい場合がある  
学術的な名前の deflate と Accept-Encoding の名前の deflate が名前に歪みが生じている弊害だろうな  
だから 混乱を生まないためにも gzip を使うことが多いらしい  

gzip は 標準ライブラリの [compress/gzip](https://pkg.go.dev/compress/gzip@go1.17.8) が使える  

どのライブラリにも 以下が定義されている
```go
const (
	NoCompression      = 0
	BestSpeed          = 1
	BestCompression    = 9
	DefaultCompression = -1
	HuffmanOnly        = -2
)
```

サンプルでは 圧縮するときに `NewWriter()` を使うが 中身は `NewWriterLevel()` だったりする  


## 参考文献
[Accept-Encoding](https://developer.mozilla.org/ja/docs/Web/HTTP/Headers/Accept-Encoding)  
[deflate と zlib と gzip の整理](https://qiita.com/ryskiwt/items/5ca10826252390a15d10)  
