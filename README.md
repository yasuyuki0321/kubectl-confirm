# kubectl-confirm

## 概要

kubectlコマンド実行時に、コンテキス及び実行するコマンドの確認プロンプトを表示する

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

### 前提

kubectlがインストールされていること　　
https://kubernetes.io/ja/docs/tasks/tools/install-kubectl/

### インストール

```shell
go get github.com/yasuyuki0321/kubectl-confirm
cd $GOPATH/src/github.com/yasuyuki0321/kubectl-confirm
make install
```

### aliasの登録

```shell
echo "alias kubectl='kubectl confirm'" >> ~/.zshrc
source ~/.zshrc
```

※ alias登録後、一時的にconfirmオプションなしの状態でkubectlを実行したい場合には、
kubectlの前に `¥` を付けて実行する。

※ kubectl-confirm コマンドが実行できない場合には、下記のようにPATHを追加する

```shell
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
source ~/.zshrc
```

### 動作確認

kubectl edit コマンドを実行し、下記のように確認のメッセージが表示されることを確認する

```shell
❯ kubectl edit
#--------------------
# Context: clp-eks
# Command: kubectl edit
#--------------------
Are you sure? [Y/n]:
```

## カスタマイズ

### 確認対象外コマンドの登録

確認無しで実行可能なコマンドを `config/exclude_commands.conf` に追加/削除をする
exclude_commands.conf 変更後は `make install` で再インストールがを実施

デフォルトの設定

```shell
config
describe
diff
exec
explain
get
kubectl
kustomize
logs
top
```
