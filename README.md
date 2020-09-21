# kubectl-confirm

kubectl コマンド実行前に、コンテキスと及び実行するコマンドの確認プロンプトを表示する

実行例

```shell
❯ kubectl get pods
#--------------------
# Context: dev-eks
# Command: kubectl get pods
#--------------------
Are you sure? [Y/n]: Y
NAME                                             READY   STATUS    RESTARTS   AGE
test-nginx-59c9f8dff-86h48                       1/1     Running   0          21h
```

## インストール方法

```shell
make install
```

## aliasの登録

```shell
echo "alias kubectl='kubectl confirm'" > ~/.zshrc
source ~/.zshrc
```

※ alias登録後、一時的にconfirmオプションなしの状態でkubectlを実行したい場合には、
kubectlの前に `¥` を付けて実行する。

## カスタマイズ

### 確認対象外コマンドの登録

確認無しで実行可能なコマンドを `config/exclude_commands.conf` に追加/削除をする
exclude_commands.conf 変更後は `make install` で再インストールが必要

デフォルトの設定

```shell
get
describe
top
logs
exec
diff
kustomize
config
```
