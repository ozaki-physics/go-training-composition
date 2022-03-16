# CoinMarketCap の API について
[API ドキュメント トップ(外部リンク)](https://coinmarketcap.com/api/)  
[API ドキュメント(外部リンク)](https://coinmarketcap.com/api/documentation/v1#section/Standards-and-Conventions)  
## 概要
API を使うには アカウント登録が必要  
必要に応じて 有料アカウントにアップグレードでき 基本は無料の Basic プランで良い  
[API Plan Feature Comparison(API プランの特徴比較)](https://pro.coinmarketcap.com/api/features)  

- Sandbox 環境
  - エンドポイント: `https://sandbox-api.coinmarketcap.com`  
  - API Key: `b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c`  
- Live 環境
  - エンドポイント: `https://pro-api.coinmarketcap.com`  
  - API Key: `ログインして Developer Portal の API Key をコピー`  


API Key を渡してリクエストする方法は2種類ある  
1. カスタムヘッダーに key を格納してリクエストする  
  `"X-CMC_PRO_API_KEY: b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c"`  
2. URL にクエリ文字列パラメータとして 渡す  
  `?CMC_PRO_API_KEY=b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c`  

1の __カスタムヘッダーに格納するやり方が推奨__  
なぜなら クエリ文字列パラメータ だと key が公開されてしまい セキュリティ的にアウトだから  
言われてみれば 確かに PayPal も リクエストヘッダーに格納してたなぁ  

### limit rate
hard cap と呼ばれる 呼び出し回数の上限が存在する  
ステータスコード 200 の successful に対して 1回カウントされる  
Basic プランでは 1日333回, 1か月10000回まで  
API の limit rate になったら 429 "Too Many Requests" になる  

## API の概要
CoinMarketCap(略称 CMC) の API は 8種類の top-level カテゴリー がある  
- `/cryptocurrency/*`  
  暗号通貨 に関する情報を返す 暗号通貨のリスト, 価格, volume(流通量?)など  
- `/exchange/*`  
  暗号通貨取引所 に関する情報を返す 取引所のリスト, マーケットのペアなど  
- `/global-metrics/*`  
  マーケットの集約データ に関する情報を返す BTC の優位性など  
- `/tools/*`  
  便利 Util に関する情報を返す fiat currency(フィアット カレンシー 法定通貨 USDなど) への換算など  
- `/blockchain/*`  
  ブロックチェーンの block explorer に関する情報を返す  
- `/fiat/*`  
  CMC ID へのマッピングを含む フィアット通貨 に関する情報を返す  
- `/partners/*`  
  3rd party な 暗号データ に関する情報を返す  
- `/key/*`  
  API Key 管理 に関する情報を返す  

また endpoint の末は以下のようなルールに従っているらしい
- `*/latest` 最新データ  
- `*/historical` チャートライブラリで使用するデータ  
- `*/info` メタデータ  
- `*/map` CoinMarketCap の ID マップ

cryptocurrency と exchange では 目的に応じて2種類の方法でアクセスできる  
- listing endpoints(リスティングエンドポイント) `*/listings/* `  
  ソートやフィルタリングできる  
- item endpoints(アイテムエンドポイント) `*/quotes/*` や `*/market-pairs/*`  
  特定の暗号通貨セットの市場相場を1回の呼び出しで取得できる  

http リクエストには ヘッダーに `Accept: application/json` が必要  
また 高速かつ効率的に アクセスするために `Accept-Encoding: deflate, gzip` も必要  

### endpoint の format  
成功コールでも失敗コールでも なるべく status オブジェクトは返ってくる  
status オブジェクトには  
- timestamp  
- credit_count この呼び出しが利用した API call credit 数  
- elapsed リクエスト処理に使った ミリ秒  
- error_code 0 なら Error なし  
- error_message null なら Error なし  
など含まれる  
```json
{
  "data" : {
    ...
  },
  "status": {
    "timestamp": "2018-06-06T07:52:27.273Z",
    "error_code": 400,
    "error_message": "Invalid value for \"id\"",
    "elapsed": 0,
    "credit_count": 0
  }
}
```

### 暗号通貨(cryptocurrency), 取引所(exchange), fiat currency の識別子  
暗号通貨なら 例えば  
eg. id=1 for Bitcoin  
eg. symbol=BTC for Bitcoin  
取引所なら 例えば  
eg. id=270 for Binance  
eg. slug=binance for Binance  
fiat currency なら 例えば  
eg. USD for the US Dollar  

貴金属換算もできる(ただ単位が Troy Ounce)  
ヤード・ポンド法の質量の単位であり、1トロイオンス = 31.103 4768グラム  
Gold(金), Silver(銀), Platinum(プラチナ), Palladium(パラジウム) の 4種  

__CoinMarketCap ID を使用することを常に推奨__  
- 暗号通貨のシンボルが一意でない  
- 暗号通貨のリブランドで変更される  

の可能性があるから  
/map で素早く CoinMarketCap ID を見つけよう  

### 法定通貨
convert パラメータで fiat currency を設定できる  

### バンドリングされた API Calls  
endpoint に カンマ区切りな値を渡して 複数の項目を1回で取得したり 変換できる  
バンドリング対応した endpoint は データを配列ではなく オブジェクトで返す  
オブジェクトで返すときの key は こちらが渡した値になっている  
`/v1/cryptocurrency/quotes/latest` に対して  
`symbol=BTC,ETH` なら  
```json
"data" : {
    "BTC" : {
      ...
    },
    "ETH" : {
      ...
    }
}
```
`id=1,1027` なら  
```json
"data" : {
    "1" : {
      ...
    },
    "1027" : {
      ...
    }
}
```

### Time format
ISO 8601 format (eg. 2018-06-06T01:46:40Z) か Unix time (eg. 1528249600)  

## ベストプラクティス  
- 暗号通貨のシンボルじゃなくて CoinMarketCap の ID を使え  
  シンボルは一意じゃなかったり 変わったりするから  
- やりたいことに適した endpoint を使え  
  異なる format だけど 同じデータというのはある  
  あるシンボルの情報を 一覧から取得するか 選択的に要求するか使い分けろ  
- キャッシュ戦略を導入しろ  
  毎回 call するのではなく rate が 60 sec ごとしか変わらないなら 60 sec 間キャッシュするなど  
- 防御的なコーディングをしろ  
  - レスポンスの解析ロジックが堅牢なように 正規表現じゃなくて JSON として解析しろ  
  - 解析コードは 必要とするプロパティのみ解析するようにしろ(将来追加される可能性のある新しいフィールドが無視されるように)  
  - レスポンス解析ロジックでは 堅牢なフィールド検証を追加しろ(try/catch 文でラップしろ)  
  - REST API 呼び出しロジックに "Retry with exponential backoff" コーディングパターンを実装しろ  

Retry with exponential backoff とは  
リトライを 指数関数的に 遅くしていく手法  
[リトライ処理の効率的アプローチ「Exponential Backoff」の概要とGoによる実装](https://qiita.com/po3rin/items/c80dea298f16a2625dbe)  

## 使いそうな API
endpoint ごとに どのプランで使えるか など 書いてあった  
基本的に使いそうな API は  
- Quotes Latest v2  
GET /v2/cryptocurrency/quotes/latest  
- CoinMarketCap ID Map  
GET /v1/cryptocurrency/map  
- Metadata v2  
GET /v2/cryptocurrency/info  
- Key Info  
GET /v1/key/info  

### Quotes Latest v2 (GET /v2/cryptocurrency/quotes/latest)
銘柄に関する情報が取得できる  
クエリパラメータがあるから 必要に応じて ドキュメント を読もう  
```json
{
  "data": {
    "1": {
      "id": 1,
      "name": "Bitcoin",
      "symbol": "BTC",
      "slug": "bitcoin",
      "is_active": 1,
      "is_fiat": 0,
      "circulating_supply": 17199862,
      "total_supply": 17199862,
      "max_supply": 21000000,
      "date_added": "2013-04-28T00:00:00.000Z",
      "num_market_pairs": 331,
      "cmc_rank": 1,
      "last_updated": "2018-08-09T21:56:28.000Z",
      "tags": [
        "mineable"
      ],
      "platform": null,
      "self_reported_circulating_supply": null,
      "self_reported_market_cap": null,
      "quote": {
        "USD": {
          "price": 6602.60701122,
          "volume_24h": 4314444687.5194,
          "volume_change_24h": -0.152774,
          "percent_change_1h": 0.988615,
          "percent_change_24h": 4.37185,
          "percent_change_7d": -12.1352,
          "percent_change_30d": -12.1352,
          "market_cap": 852164659250.2758,
          "market_cap_dominance": 51,
          "fully_diluted_market_cap": 952835089431.14,
          "last_updated": "2018-08-09T21:56:28.000Z"
        }
      }
    }
  },
  "status": {
    "timestamp": "2022-03-10T16:13:24.430Z",
    "error_code": 0,
    "error_message": "",
    "elapsed": 10,
    "credit_count": 1
  }
}
```

### CoinMarketCap ID Map (GET /v1/cryptocurrency/map)
CoinMarketCap ID Map を調べるとき に使いそう  
1度 id が分かれば2度とアクセスしなくて良いかも  
クエリパラメータがあるから 必要に応じて ドキュメント を読もう  
_おそらく Sandbox 環境のレスポンスは間違っているから注意_  
クエリ文字列パラメータで渡した symbol が key になってたりするが Live 環境はそうではない  
Sandbox 環境は クエリ文字列パラメータが機能してないように思われる  

銘柄の name, Symbol, Platform まで分かるのは便利  
```json
{
  "data": [
    {
      "id": 1,
      "rank": 1,
      "name": "Bitcoin",
      "symbol": "BTC",
      "slug": "bitcoin",
      "is_active": 1,
      "first_historical_data": "2013-04-28T18:47:21.000Z",
      "last_historical_data": "2020-05-05T20:44:01.000Z",
      "platform": null
    },
    {
      "id": 1839,
      "rank": 3,
      "name": "Binance Coin",
      "symbol": "BNB",
      "slug": "binance-coin",
      "is_active": 1,
      "first_historical_data": "2017-07-25T04:30:05.000Z",
      "last_historical_data": "2020-05-05T20:44:02.000Z",
      "platform": {
        "id": 1027,
        "name": "Ethereum",
        "symbol": "ETH",
        "slug": "ethereum",
        "token_address": "0xB8c77482e45F1F44dE1745F52C74426C631bDD52"
      }
    }
  ],
  "status": {
    "timestamp": "2018-06-02T22:51:28.209Z",
    "error_code": 0,
    "error_message": "",
    "elapsed": 10,
    "credit_count": 1
  }
}
```
### Metadata v2 (GET /v2/cryptocurrency/info)
公式サイト, エクスプローラ, ロゴ, ホワイトペーパー, とか分かるの便利そう  
クエリパラメータがあるから 必要に応じて ドキュメント を読もう  
key が クエリ文字列パラメータ なのが Go では扱いにくい?
```json
{
  "data": {
    "1": {
      "urls": {
        "website": [
          "https://bitcoin.org/"
        ],
        "technical_doc": [
          "https://bitcoin.org/bitcoin.pdf"
        ],
        "twitter": [],
        "reddit": [
          "https://reddit.com/r/bitcoin"
        ],
        "message_board": [
          "https://bitcointalk.org"
        ],
        "announcement": [],
        "chat": [],
        "explorer": [
          "https://blockchain.coinmarketcap.com/chain/bitcoin",
          "https://blockchain.info/",
          "https://live.blockcypher.com/btc/"
        ],
        "source_code": [
          "https://github.com/bitcoin/"
        ]
      },
      "logo": "https://s2.coinmarketcap.com/static/img/coins/64x64/1.png",
      "id": 1,
      "name": "Bitcoin",
      "symbol": "BTC",
      "slug": "bitcoin",
      "description": "Bitcoin (BTC) is a consensus network that enables a new payment system and a completely digital currency. Powered by its users, it is a peer to peer payment network that requires no central authority to operate. On October 31st, 2008, an individual or group of individuals operating under the pseudonym \"Satoshi Nakamoto\" published the Bitcoin Whitepaper and described it as: \"a purely peer-to-peer version of electronic cash would allow online payments to be sent directly from one party to another without going through a financial institution.\"",
      "date_added": "2013-04-28T00:00:00.000Z",
      "date_launched": "2013-04-28T00:00:00.000Z",
      "tags": [
        "mineable"
      ],
      "platform": null,
      "category": "coin"
    },
    "1027": {
      "urls": {
        "website": [
          "https://www.ethereum.org/"
        ],
        "technical_doc": [
          "https://github.com/ethereum/wiki/wiki/White-Paper"
        ],
        "twitter": [
          "https://twitter.com/ethereum"
        ],
        "reddit": [
          "https://reddit.com/r/ethereum"
        ],
        "message_board": [
          "https://forum.ethereum.org/"
        ],
        "announcement": [
          "https://bitcointalk.org/index.php?topic=428589.0"
        ],
        "chat": [
          "https://gitter.im/orgs/ethereum/rooms"
        ],
        "explorer": [
          "https://blockchain.coinmarketcap.com/chain/ethereum",
          "https://etherscan.io/",
          "https://ethplorer.io/"
        ],
        "source_code": [
          "https://github.com/ethereum"
        ]
      },
      "logo": "https://s2.coinmarketcap.com/static/img/coins/64x64/1027.png",
      "id": 1027,
      "name": "Ethereum",
      "symbol": "ETH",
      "slug": "ethereum",
      "description": "Ethereum (ETH) is a smart contract platform that enables developers to build decentralized applications (dapps) conceptualized by Vitalik Buterin in 2013. ETH is the native currency for the Ethereum platform and also works as the transaction fees to miners on the Ethereum network.Ethereum is the pioneer for blockchain based smart contracts. When running on the blockchain a smart contract becomes like a self-operating computer program that automatically executes when specific conditions are met. On the blockchain, smart contracts allow for code to be run exactly as programmed without any possibility of downtime, censorship, fraud or third-party interference. It can facilitate the exchange of money, content, property, shares, or anything of value. The Ethereum network went live on July 30th,2015 with 72 million Ethereum premined.",
      "notice": null,
      "date_added": "2015-08-07T00:00:00.000Z",
      "date_launched": "2015-08-07T00:00:00.000Z",
      "tags": [
        "mineable"
      ],
      "platform": null,
      "category": "coin",
      "self_reported_circulating_supply": null,
      "self_reported_market_cap": null,
      "self_reported_tags": null
    }
  },
  "status": {
    "timestamp": "2022-03-10T16:13:24.430Z",
    "error_code": 0,
    "error_message": "",
    "elapsed": 10,
    "credit_count": 1
  }
}
```

### Key Info (GET /v1/key/info)
使用状況やいつ回数リセットか など 確認できる
