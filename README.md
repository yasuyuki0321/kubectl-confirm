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

# インストール方法
```shell
make install
```

# aliasの登録
```shell
echo "alias kibectl='kubectl confirm'" > ~/.zshrc
source ~/.zshrc
```
※ alias登録後、一時的にconfirmオプションなしの状態でkubectlを実行したい場合には、
kubectlの前に `¥` を付けて実行する。
