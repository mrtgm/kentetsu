# kentetsu

<p align="center">
  <img src="kentetsu_pixel.png" width="200" alt="Koh Kentetsu pixel art">
</p>

[Koh Kentetsu Kitchen](https://www.youtube.com/@kohkentetsu) からランダムにレシピを1つ推薦する CLI ツール。

## Install

```sh
brew install mrtgm/tap/kentetsu
```

or

```sh
go install github.com/mrtgm/kentetsu@latest
```

## Usage

```sh
kentetsu                  # ランダム1品
kentetsu -search 鶏肉     # キーワード絞り込み
kentetsu -open            # ブラウザで開く
kentetsu -update          # 最新データに更新
```
