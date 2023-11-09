# forge

## 依存解決の注意

現在、[docker-image-disassembler](https://github.com/tklab-group/docker-image-disassembler)がprivate repositoryなので、最新バージョンの取得やビルド時に以下の設定が必要になることがある
- `GOPRIVATE=github.com/tklab-group`
- https://go.dev/doc/faq#git_https のいずれかの設定