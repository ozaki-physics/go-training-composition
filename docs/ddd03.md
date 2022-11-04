# DDD をやってみる
## DDD を試して思った見解
### 各層の責務
- domain 層:  
  - ドメインオブジェクト(entity と 値オブジェクト)  
    - ドメイン知識(ルール/制約)の表現  
  - ドメインサービス  
  - リポジトリ(インタフェース)の仕様定義  
- usecase 層:  
  - ユースケースの実現  
    - ドメインオブジェクトの 生成, 使用, 永続化依頼  
    - ドメインオブジェクト を presen 層に渡す値への変換  
- infra 層:  
  - リポジトリ(実装)の処理定義  
    - ドメインオブジェクト の 永続化/検索 の実装  
- presen 層:  
  - エンドポイントの定義  
  - http Request で渡された値 と usecase 層に渡す値 の マッピング  
  - 入力値のバリエーション(一部)  

### 1パッケージ内(domain 層, usecase 層など)にファイル名が増えていくが いいのか?  
Go の思想として 1パッケージ内に結構いろいろ存在してもいいことになっている  
### 本当は error を 各層で独自の型 にした方がいいと思う  
すべて catch するのは presen 層でもいい  
<!-- domain 層 寄りの話 -->
### 値オブジェクト と プリミティブ型 の使い分けは?  
値オブジェクト にするということは 型安全は得られるが それだけコスト(コード作成, 管理, 保守)がかかる  
ドメインを鑑みて 値オブジェクト にするだけの価値があるなら 値オブジェクト にする  
### entity の 構造体のフィールドは private にして ゲッターメソッドを定義すると セッターが存在せずカプセル化できてると思う  
entity の生成, 値変更には コンストラクタ や リコンストラクタ を使うしないという制約を設けて 不整合な entity 生成を防げる  
### 値オブジェクト を プリミティブ型のエイリアス を定義するだけでいいのか?  
厳密にしたいなら 値オブジェクト に private フィールド を用意した方がいい  
### コンストラクタ と ファクトリーメソッド(リコンストラクタ など) はどのように使い分ける?  
Go の言語仕様では コンストラクタ が無いから 使い分けが必要な発想になる?  
1つの案として 引数に受け取るもの を コンストラクタ は プリミティブ型, ファクトリーメソッド は 値オブジェクト で使い分けるのは有りかも?  
### コンストラクタにプリミティブ型の引数が多くなったら もう少し詳細に概念分けをした方がいいか検討する機会?  
その都度検討するしかない  
### 構造体 に メソッドを作る時  
主語が何になるか考えて `対象.動詞(対象)` という形になるように 意識すると良いかも  
"ユーザー に 利用権限 を 付与する" なら `ユーザー.付与(利用権限)` など  
### domain 層の リポジトリ(repo) には全件取得メソッド(fetch) を定義した方がいいと思う  
<!-- usecase 層 寄りの話 -->
### usecase に CRUD を書かない  
CRUD は リポジトリ(インタフェース)に定義すること  
usecase は 少し抽象的にやりたいことを書く  
### 状態の変化は usecase から domain のメソッドを呼んで domain の状態を変化させる  
usecase で domain から値を取り出して 処理して domain に格納し直すのは違う  
保存するときに `dao.save(dto)` と 保存したい対象を引数に渡すのではなく `domain.save()` が美しいと思う  
### domain 層 に infra 層 が依存する インタフェース が存在するように usecase 層 に presen 層が依存する インタフェース を用意しなくていいのか?  
presen 層 に変化が無いのに usecase 層だけ変わるってことがありうるか? 発生頻度低そうだし 無駄が多そう  
テストの観点からは? presen 層 だけのテストができるようになる 現時点 不要と思うなら作らない 必要になったときに増やせばいい  
usecase メソッド一覧が欲しいなら `go doc -u -all ./usecase` ってやればいい  
### dto は usecase 層 から presen 層 にデータを運ぶ構造体  
### dto にゲッターを作る方が良いのか?  
JSON にするためには 構造体のフィールドが public にする必要がある かつ dto まで ゲッター 作るのは冗長と思われる  
<!-- infra 層 寄りの話 -->
### JSON で保存するときに infra 層で保持してる状態(entity の map)から JSON 用の構造体に詰め替えるのが面倒  
JSON にするためには 構造体のフィールドを public にする必要がある  
entity のフィールドは private かつ その面倒を担うのが infra 層 の役割だから仕方ない?  
JSON 用の 構造体は プリミティブ型じゃなくてもいいなら せめて 値オブジェクト とか使いたい  
<!-- presen 層 寄りの話 -->
### ブラウザで閲覧する or CSV でダウンロードする の切り替えは プレゼン層の責務と思う  
### dto から JSON 用の 構造体 に 変換するメソッドがあると楽かも?  
データ永続層(JSON)だから ORM が存在してて欲しいという主張に聞こえる  


## DDD を試して思った命名規則
### 現時点では 良さそうな命名規則  
- domain 層:  
  - foo.go: ドメインオブジェクト(entity), ドメインモデル  
  - bar_repo.go: リポジトリ(インタフェース)で infra 層の定義  
  - baz.go: ドメインオブジェクト(値オブジェクト), プリミティブ型のラッパー的な  
- usecase 層:  
  - foo_usecase.go: ユースケース, Logic, 処理そのもの  
  - bar_dto.go: 層を横断するときの構造体  
- infra 層:  
  - bar_json.go: データ永続層 とのやり取り, 戻り値が domain 層のリポジトリで実装を強制する  
  - bar_rdb.go: データ永続層 とのやり取り, 戻り値が domain 層のリポジトリで実装を強制する  
  - share: infra 層のいろんな所で使う処理をまとめたもの  
    - json_detail.go: JSON からの読み込み, JSON への保存  
- presen 層:  
  - foo_endpoint.go: GET, POST などの処理  
  - foo_json.go: 返す JSON の構造体  
  - api_handler.go: リクエスト(GET, POST)に応じて go の起動するメソッドを振り分ける  
  - csv_detail.go: CSV への保存(書き出し)  
- /configs: 環境変数的な設定ディレクトリ  
  - json.go: データ永続層(JSON) のディレクトリ path 一覧  
  - csv.go: csv を保存する path?  

### ディレクトリ名 が長くなるときはどのように命名したらいいのか?  
kubernetes や terraform が ケバブケース(foo-bar 小文字のみ)を使っていることがある  
### パッケージ名はこれでいいのか?  
go の思想的に パッケージ名は 1単語(domain, usecase, infra, presen, share)になるからいいと思う  
パッケージ名が長くなってもに アンダースコア を使うことはできるが ハイフンは使えない  
テストコードだけ `_test` サフィックス をつける  
### ディレクトリ名と同一にした方が無難  
一般的には ディレクトリ名 は ハイフン, パッケージ名 では アンダースコア になる  
一般的ではないが 僕の好みで ディレクトリ名も パッケージ名と同一にするために アンダースコアとする  
なぜなら なるべく一致させた方が検索とかしやすそうだから  

### ファイル名にアンダースコア(スネークケース)を使うのはいいのか?  
kubernetes や terraform が スネークケース(foo_bar 小文字のみ)を使っていることがあるので 許容  
### ファイル名, 構造体名 に どの層かの情報 を付与しない(構造体名に JSON, CSV というプレフィックスをつけることはあるかも)  
どの層かの情報 を付与するかは 判断が2転3転している  

過去 DDD で作ったときに ファイル名 と 構造体名 が乖離して どこに何が書いてあるか分からなくなった  
検索を楽にするという観点では 1クラス1ファイルのように 1構造体1ファイル にすれば 楽になる?  
VSCode で検索するときは ファイル名で検索するから ファイル名 と 構造体名 は同期させた方がいい?  
だから ファイル名, 構造体名 は どの層かの情報 を付与した方がいいと思っていた  
`foo_infra.go`, `type fooInfra struct{}`  

でも VSCode は検索でディレクトリ名も検索対象にできるから ファイル名 にわざわざ どの層かの情報 を付与しなくてもいいかも?  
ただ 層を横断して構造体を使っているときは 構造体名 から どの層の構造体なのか判断できるようにしたくて どの層かの情報 は付与した方がいいかも?  
よって ファイル名には どの層かの情報 を付与しない, 構造体には どの層かの情報 を付与する と決めた  
ファイル名 と 構造体名 が同期しないが諦めようと思った  
`infra/foo.go`, `type fooInfra struct{}`  

しかし 考え直してみると  
構造体名 に どの層かの情報 を付与してると 層を横断して使うと `domain.fooDomain` と domain と2回書かれるのは冗長になるのではって思った  
また 構造体 を private で定義したら 同一パッケージ内(同一層)でしか使えないため 層の特定ができると思った  
よって ファイル名 にも 構造体名 にも どの層かの情報 を付与しなくても大丈夫だと思われる  
`infra/foo.go`, `type foo struct{}`  
ただ 同一パッケージ内(同一層) で バリエーションが生まれるとき(presen 層, infra 層 で JSON, CSV 的なもの)は ファイル名, 構造体名 に情報を付与していいと思う  
`infra/foo_json.go`, `type fooJson struct{}`  
`infra/foo_csv.go`, `type fooCsv struct{}`  
### 過去に良いと思っていた命名規則  
- domain 層:  
  - foo_entity.go: ドメインモデルなイメージ  
  - foo_repo.go: リポジトリ, インタフェースで infra 層の実装元  
  - foo_value.go: 値オブジェクト, プリミティブ型のラッパー的な  
- usecase 層:  
  - foo_usecase.go: ユースケース, Logic, 処理 そのもの  
  - foo_dto.go: 層を横断するときの構造体, presen 層が求める構造体(foo_by_bar_usecase_dto.go など)  
- infra 層:  
  - foo_infra.go: データ永続層 とのやり取り, 戻り値が domain 層のリポジトリで実装を強制する(foo_json_infra.go など)  
  - conf.go: データ永続層(JSON) のディレクトリ path 一覧  
  - share: infra 層のいろんな所で使う処理をまとめたもの  
    - json_detail.go: JSON からの読み込み, JSON への保存  
- presen 層:  
  - foo_end_point.go: GET, POST などの処理  
  - foo_json.go: 返す JSON の構造体  
  - api_handler.go: リクエスト(GET, POST)に応じて go の起動するメソッドを振り分ける  
  - csv_conf.go: csv を保存する path?  
  - csv_detail.go: CSV への保存(書き出し)  
#### 過去に良いと思っていた命名規則(具体例)  
- domain 層:  
  - account_domain.go  
  - account_repo.go  
  - explain_value.go  
- usecase 層:  
  - role_usecase.go  
  - user_dto.go  
  - user_by_account_dto.go  
- infra 層:  
  - share  
    - json_detail.go  
  - conf.go  
  - account_json_infra.go  
  - account_rdb_infra.go  
- presen 層:  
  - api_handler.go  
  - csv_detail.go  
  - csv_conf.go  
  - account_by_service_end_point.go  
  - account_by_service_json.go  



## DDD を試して思った Go の言語仕様
### Go 言語の仕様として スライスには削除メソッドがないが map には削除メソッドがある
### Go の思想として インタフェース名  
メソッドが 1 つの場合: `メソッド名 + er/or` (Write メソッドを持つインターフェースなら Writer)  
目的語がある場合は `目的語 + 動詞 + er/or` のパターンもあるらしい `ObjectPrinter`  
メソッドが 2 つ以上の場合: `メソッドを連ねて er` (ReadCloser)  
### go の思想として cmd ディレクトリに全部書かなくていいのか?  
気にしていたら前に進まなくなるから 必要に応じて ディレクトリ作ろう([メルカリの GitHub](https://github.com/mercari) を覗いても すべてに cmd ディレクトリがあるわけではなかった)  
### メソッド名, 構造体名は?  
Go の思想として キャメルケース(fooBar)  
### 変数名や定数名は?  
変数名は短いもの, 定数名は アンダースコアは許容  
### 変数名を大文字にする例外は?  
http を HTTP とかはしてる, 略語は大文字にするか迷う  
もし coin market cap id で 省略しつつキャメルケースにすると CmcId になる  
略語は大文字にしつつだと CMCId になる  

[Initialisms](https://github.com/golang/go/wiki/CodeReviewComments#initialisms)  
>Words in names that are initialisms or acronyms (e.g. "URL" or "NATO") have a consistent case.  
>For example, "URL" should appear as "URL" or "url" (as in "urlPony", or "URLPony"), never as "Url".  
>As an example: ServeHTTP not ServeHttp. For identifiers with multiple initialized "words", use for example "xmlHTTPRequest" or "XMLHTTPRequest".  
>>頭文字や頭字語である名前の単語（たとえば「URL」や「NATO」）は、大文字と小文字を統一してください。  
>>例えば、「URL」は「URL」または「url」（「urlPony」、「URLPony」のように）と表示されるべきで、決して「Url」と表示してはいけません。  
>>例としてServeHTTPは、ServeHttpではありません。複数の単語が初期化されている識別子の場合、例えば「xmlHTTPRequest」または「XMLHTTPRequest」を使用します。  

一般的に普及した略称じゃないなら 省略しない方がいいと思われる  
省略したくなるほど 何度も出てくるなら 変数名が無駄に長い or カプセル化が上手にできてない と設計を見直した方がいいと思われる  
