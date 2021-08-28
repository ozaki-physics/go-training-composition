FROM golang:1.17

# ENV は コンテナ内でも環境変数として利用したい値(PATH環境変数など)
# ARG は ビルド時にのみに利用したい値(インストールするミドルウェアのバージョンなど)
ENV REPOSITORY=github.com/ozaki-physics/go-training-composition

# go mod を使うために GOPATH を削除するために path を上書き
ENV GOPATH=""
WORKDIR /go/src/$REPOSITORY

# 本当は build だけで使える image を作成すべき
# 相対パスで Dockerfile より上の階層は指定できないからコピーできなかった
COPY . .
# 毎回 build すると余計な image ばかり作っちゃうが go.mod を更新したら imgae 作り直してもいいかも
RUN go mod download
# CMD ["go" "run" "main.go"]
CMD [ "bash" ]
