# database/sql パッケージ の勉強
## 参考書籍
『詳解Go言語Webアプリケーション開発』
ISBNコード: 978-4-86354-372-0
著者: 清水 陽一郎
初版: 2022/08/01
## database/sql の使い方
DB へ接続するとき 最初の操作は `sql.Open()` [`func Open(driverName, dataSourceName string) (*DB, error)`](https://pkg.go.dev/database/sql#Open)  
`*sql.DB` 型の値を介してクエリを実行したリ トランザクション を開始する  
`*sql.DB` 型 の内部でコネクションプールを管理してくれているから `*sql.DB` 型の値は Web アプリ起動時に1回作成すればよい  

Xxx() と XxxContext() があるときは XxxContext() を優先して使う  
Xxx() は context が実装される前の後方互換のために存在する  

`sql.ErrNoRows` を返すのは `*Row` だけで `*Rows` は絶対に発生しない  
複数行のレコードを取得するメソッドでは発生しない  

トランザクションを使うときは `defer` で `Rollback()` を呼ぶ  
すると メソッドスコープが完了するときに 絶対実行されるため 下記忘れ や 何度も書く手間がなくなる
アンチパターン  
```go
func (r *Repository) Update(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err := tx.Exec(/* sql */)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err := tx.Exec(/* sql */)
	if err != nil {
		// ロールバック忘れ
		return err
	}

	return tx.Commit()
}
```
デザインパターン  
```go
func (r *Repository) Update(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
  defer tx.Rollback()

	_, err := tx.Exec(/* sql */)
	if err != nil {
		return err
	}

	_, err := tx.Exec(/* sql */)
	if err != nil {
		return err
	}

	return tx.Commit()
}
```

## database/sql の代わりによく使われる OSS  
- [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx)  
使用感は database/sql に近く SQL を手で書く  
- [github.com/ent/ent](https://github.com/ent/ent)  
SQL を書かずにスキーマや RDBMS の操作を自動生成  
- [gorm.io/gorm](https://gorm.io/)  
Go で有名な ORM  
- [github.com/kyleconroy/sqlc](https://github.com/kyleconroy/sqlc)  
SQL から RDBMS にアクセスできる型安全なコードを自動生成  
- [github.com/volatiletech/sqlboiler/v4](https://github.com/volatiletech/sqlboiler)  
RDBMS のテーブル定義から コードを自動生成  
