# 開発について

## リリース

[goreleaser/goreleaser](https://github.com/goreleaser/goreleaser)を使用したCIでのリリースの設定をしている。

関連ファイル:
- `.github/workflows/release.yml`
- `.goreleaser.yaml`

新バージョンをリリースするときは、`v0.0.0`形式(e.g. `v1.2.3`)のgitタグを作成してpushするだけ。
GitHub Actionsが回ってバイナリの作成や[Releases](https://github.com/tklab-group/forge/releases)へのアップロードが行われる。


## テスト

実行コマンドの管理に[Task](https://taskfile.dev)を使用している。

関連ファイル:
- `Taskfile.yaml`

### Golden Testの更新

多くのテストは[sebdah/goldie](https://github.com/sebdah/goldie)によるGoledn Testを採用しているため、更新が必要なときは`task update-golden-test`で一括更新ができる。

`go test`コマンドに対しての`-update`や`-clean`のフラグを引数(`./...`)の前と後のどちらに置くべきかが環境依存がありそう(根本原因は未特定)なため、意図通り更新が走らない場合はフラグの位置を動かしてみること。