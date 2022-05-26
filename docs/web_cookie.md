# Cookie の勉強
## Cookie とは
### 概要
正式名は HTTP Cookie らしい  
Cookie は サーバ が ブラウザ に渡す小さなデータで ブラウザに保存される  
Cookie を送ってきた サーバ にリクエストするときに 一緒に その Cookie も送信する  
http は ステートレス(状態を保持しない) だったが 状況によって異なる内容を表示する(動的ページ)のニーズがあり生まれた  
用途は Session 管理(ログイン情報など), パーソナライズ(ユーザー固有の設定), トラッキング(ページ閲覧行動記録) など  
RFC 6265 などで 仕様が定義されている  

Cookie は リクエストヘッダー や レスポンスヘッダー に含まれて やり取りされる  

仕様は毎回少しづつ変わってきていて 新たな属性(SameSite 属性)や プレフィックス(__Secure- と __Host-) が増えたりしている  

Cookie のサンプル: `Set-Cookie: <cookie-name>=<cookie-value>`  

1個の Cookie には 4096(2の12乗) バイト = 4KB まで  
1文字2バイトだから 2048文字まで  
1台のサーバが同じブラウザに発行できる Cookie は20個まで  
ブラウザに保存できる Cookie は 300個まで  
これを超えると古いものや有効期限が過ぎたものが破棄される  
### 接続時間の定義  
- Session Cookie は セッションが終わると終了  
- Expires 属性, Max-Age 属性 で 時間を設定する  

ブラウザによっては Session を復元するって機能があるから 実質無制限になったりする  

推奨は Max-Age 属性 だけど IE が対応してないとかで 実質両方を用意する必要がある  
両方用意してあったら Max-Age 属性 が優先される  

Expires 属性 は 指定された時刻で書く  
`Set-Cookie: id=a3fWa; Expires=Thu, 31 Oct 2021 07:28:00 GMT;`  
Max-Age 属性 は 指定された時間が経過したら削除  
`Set-Cookie: hello=world; Max-Age=30`  

Cookie は常に上書きしかできない  
削除するためには Expires を過去の日付にしたり Max-Age をマイナスや0にする  
### Cookie のアクセス制限  
以下2個の俗世は 十分なアクセス制限にはならず 無いよりマシ程度  
- Secure 属性  
  http では cookie を送信しない  
  https のみにできる  
  でもブラウザで確認できてしてしまう  
- HttpOnly 属性  
  js から cookie にアクセスさせない  
  クロスサイトスクリプティング(XSS)攻撃を軽減できる程度  

サンプル: `Set-Cookie: id=a3fWa; Expires=Thu, 21 Oct 2021 07:28:00 GMT; Secure; HttpOnly`  
### Cookie の送信先の定義  
- Domain 属性  
  Cookie を受信することができるホストを指定  
  Domain 属性を書かないと サブドメイン除外  
  Domain 属性を書くとサブドメイン含むため サブドメイン間でユーザー情報を共有する場合は必要  
- Path 属性  
  リクエストされた URL の中に含む必要がある URL のパス  
  例えば `Path=/docs` を設定すると `/docs/`, `/docs/Web/` は リクエストパスに一致するが `/`, `/fr/docs` は一致しない  
- SameSite 属性  
  同一サイトという権限レベルみたいものを設定する  
  値は Strict, Lax, None のいずれか  
  SameSite 属性 を設定しないと Lax になる  
  サンプル: `Set-Cookie: mykey=myvalue; SameSite=Strict`  
  - Strict: Cookie を生成したサイトのみ Cookie を送信  
  - Lax: Cookie のオリジンのサイトに移動したときにも Cookie を送信
    例えば BサイトからAサイトのリンクをクリックしたときブラウザはAサイトの cookie を送信など  
  - None: Cookie を生成したサイト と サイト間のリクエストの両方に Cookie を送信, 同時に Secure 属性 も必須   

origin オリジン とは URL のスキーム (プロトコル), ホスト (ドメイン), ポート の3セットによって定義される  
つまり 以下のパターン違うオリジン扱い  
- http と https  
- example.com と www.example.com  
- ポートなし と :8080 (ただし :80 だけは例外で同じオリジン扱いになる)  

Cookie は origin で管理しているわけではなく Domain で管理しているため プロトコル(http, https)やポートが違っても送信される  
Domain 属性をつけると サブドメインにまで Cookie が送信されるから注意  
### Cookie の接頭辞  
`__Host-`: Cookie 名の接頭辞についていたら Secure 属性と `Path=/` が必須で https かつ ドメイン固定になる(Domain 属性は指定できないためサブドメインもダメ)  
`__Secure-`: Cookie 名の接頭辞についていたら Secure 属性 が必須で https になる(`__Host-` より弱い)  
```
// どちらも安全な (HTTPS の) オリジンから受け入れられます
Set-Cookie: __Secure-ID=123; Secure; Domain=example.com
Set-Cookie: __Host-ID=123; Secure; Path=/

// 拒否
Set-Cookie: __Secure-id=1 // Secure ディレクティブが無いため
Set-Cookie: __Host-id=1; Secure // Path=/ ディレクティブが無いため
Set-Cookie: __Host-id=1; Secure; Path=/; Domain=example.com // Domain を設定したため
```
### セキュリティ  
>情報を Cookie に保存するときは、すべての Cookie の値がエンドユーザーから見え、変更できることを理解しておいてください。  
>アプリケーションによっては、サーバー側で検索される不透明な識別子を使用するか、 JSON ウェブトークンのような代替の認証/機密性メカニズムを調べたほうが良いかもしれません。  

Cookie への攻撃を緩和するには  
- HttpOnly 属性  
- 接続時間を短く  
- SameSite 属性を Strict, Lax に  
- ユーザー認証するたびに Session Cookie をすべて再生成して Session 固定攻撃を防ぐ  
### サードパーティ Cookie
ファーストパーティ Cookie とは  
Cookie ドメインとスキームが現在のページと一致しているとき その Cookie は同じサイトからの Cookie とみなされる(Domain 属性を書いたらサブドメインも)  
サードパーティ Cookie とは  
Cookie ドメインとスキームが現在のページと一致していないとき その Cookie は異なるサイトからの Cookie とみなされる  
サーバは Cookie を設定するとき SameSite 属性を指定すべきらしい  
### Cookie に関する法規制
- EU の 一般データ保護規則 (GDPR)
- EU の ePrivacy 指令
- カリフォルニア州消費者プライバシー法

規制の要件  
- サイトが Cookie を使用することをユーザーに通知する
- ユーザーが一部またはすべての Cookie をオプトアウト(停止)できるようにする
- ユーザーが Cookie を受け取らなくても サービスのほとんどを利用できるようにする

特に EU はどのようなサイトに対しても適用されるため注意が必要  
### 参考文献
[HTTP Cookie の使用](https://developer.mozilla.org/ja/docs/Web/HTTP/Cookies)  
[Origin (オリジン)](https://developer.mozilla.org/ja/docs/Glossary/Origin)  
[同一オリジンポリシー](https://developer.mozilla.org/ja/docs/Web/Security/Same-origin_policy)  
[Set-Cookie](https://developer.mozilla.org/ja/docs/Web/HTTP/Headers/Set-Cookie)  

[Cookie 概説](https://numb86-tech.hatenablog.com/entry/2020/01/19/004420)  
[SameSite Cookie の説明](https://web.dev/i18n/ja/samesite-cookies-explained/)  
[HTTP クッキーをより安全にする SameSite 属性について (Same-site Cookies)](https://laboradian.com/same-site-cookies/)  

