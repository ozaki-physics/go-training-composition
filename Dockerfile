FROM golang:1.17

# ENV は コンテナ内でも環境変数として利用したい値(PATH環境変数など)
# ARG は ビルド時にのみに利用したい値(インストールするミドルウェアのバージョンなど)
ENV REPOSITORY=github.com/ozaki-physics/go-training-composition

# go mod を使うが一応 GOPATH を指定
# go.mod で指定されたパッケージの情報は /go/pkg に格納される
# GOPAHT を空にすると go.mod で指定されたパッケージの情報が GOROOT の方に格納されるっぽい
ENV GOPATH="/go"
WORKDIR /go/src/$REPOSITORY

# 本当は build だけで使える image を作成すべき
# 相対パスで Dockerfile より上の階層は指定できないからコピーできなかった
COPY ./go.mod .
# 毎回 build すると余計な image ばかり作っちゃうが go.mod を更新したら imgae 作り直してもいいかも
RUN go mod download
# もしかして COPY するファイルに変更が加わっているから 毎回 image が作り直される説がある
# よってはじめは go.mod だけコピーして その後に ディレクトリごとコピーする
COPY . .
# そもそも go build したものを実行するから go のコマンドは使わないはず
# CMD ["go" "run" "main.go"]
# CMD [ "bash" ]
# コンテナをバックグラウンドで実行してサーバ立てられるようになった
# RUN go build main.go
# CMD ["./main"]
# マウントして上書きされないように build する場所を指定
# 本来 /go/bin は go install でインストールした 外部パッケージの場所らしいから 思想的にここで合っているかは不明
RUN go build -o /go/bin main.go
CMD ["/go/bin/main"]
