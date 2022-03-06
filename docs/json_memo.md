# Go で JSON を扱う
java では プロパティファイル(example.properties)を 設定ファイルとして読み込むことができる  
同じように Go でも設定ファイルを扱えないか?  
JSON なら読み込むのは簡単だが JSON はコメントを書けない(書けるように仕様変更されたらしい?)  
コメントが書ける JSON 的な立ち位置が YAML  
でも YAML は Go の標準ライブラリで対応してない  
外部ライブラリだと [go-yaml/yaml](https://github.com/go-yaml/yaml) や [gopkg.in/yaml.v2](https://pkg.go.dev/gopkg.in/yaml.v2) が有名らしい
でも 両方ドキュメントは [go-yaml/yaml](https://github.com/go-yaml/yaml) を見ろと書いてある??  
__設定ファイルはコメント書けないけど 標準ライブラリで終わらしたいから JSON にしよう__  

## struct に付与するメタタグ
struct に メタタグをつけると struct の要素名と json の key 名が異なっていても紐付けてくれる  
>The encoding of each struct field can be customized by the format string stored under the "json" key in the struct field's tag.  
>The format string gives the name of the field, possibly followed by a comma-separated list of options.  
>The name may be empty in order to specify options without overriding the default field name.  
>>各構造体フィールドのエンコーディングは、構造体フィールドのタグの「json」キーの下に格納されているフォーマット文字列によってカスタマイズすることが可能です。  
>>フォーマット文字列は、フィールドの名前と、その後に続くカンマ区切りのオプションのリストです。  
>>デフォルトのフィールド名を上書きせずにオプションを指定するために、名前を空にすることができます。  

```go
  Name string `json:"name"`
``` 
ただ 別パッケージで定義された struct に JSON の値を格納したい場合は struct の public にしか格納できない  
メタタグにはオプションが付けられる  
- JSON を読み込む = デコード  
- JSON を出力する = エンコード  
ハイフンは デコードにもエンコードにも作用し JSON から値を取ってこない かつ JSON にしたとき存在しない  
omitempty は JSON にしたとき(エンコード)でフィールドを省略する  

## UnmarshalJSON について
[GoでYAMLを扱うすべての人を幸せにするべく、ライブラリをスクラッチから書いた話](https://qiita.com/goccy/items/86abe72b472993b5516a)  
たしかに 標準ライブラリに 以下のようなメソッドがあるなら JSON を YAML に変えたメソッドを作りたくもなる  
- MarshalJSON() ([]byte, error)  
- UnmarshalJSON([]byte) error  

JSON を読み込むときに  
- `func Unmarshal(data []byte, v interface{}) error`
- `func (m *RawMessage) UnmarshalJSON(data []byte) error`
の違いはなんだろう  
Unmarshal() は直接 JSON を取り込む感じがして UnmarshalJSON() は構造体に入れて 構造体を経由して操作する感じがある  
[type RawMessage](https://pkg.go.dev/encoding/json@go1.17.8#RawMessage) を読んでもよく分からなかったし 結局 Unmarshal() を使っている感じだった  
