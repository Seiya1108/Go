package main

import (
    "io"
    "os"
)

// CopyFile はファイルを安全にコピーする関数です。
// srcPath: コピー元ファイルパス
// dstPath: コピー先ファイルパス
// 戻り値: エラー情報（エラーがなければ nil）
func CopyFile(srcPath, dstPath string) error {

    // 1. ソースファイルを読み取り専用モードで開く
    // os.Open は読み取り専用でファイルを開く。ファイルが存在しない場合はエラーを返す
    src, err := os.Open(srcPath)
    if err != nil {
        return err // ファイルが開けなかった場合はここでエラーを返して終了
    }
    // 関数終了時に確実にファイルが閉じられるように defer を使う
    defer src.Close()

    // 2. ソースファイルのメタデータ（サイズやパーミッションなど）を取得
    // Stat() は os.FileInfo を返す
    srcInfo, err := src.Stat()
    if err != nil {
        return err // ファイルメタ情報が取得できない場合はエラーを返して終了
    }

    // 3. デスティネーションファイルを作成または上書き
    // os.O_CREATE: ファイルが無ければ作成
    // os.O_WRONLY: 書き込み専用
    // os.O_TRUNC : 既存ファイル内容を空にしてから新たに書き込む
    // 第3引数でソースファイルと同じパーミッションを指定（srcInfo.Mode()）
    dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
    if err != nil {
        return err // ファイルが作れない/書き込みできない場合はエラーを返して終了
    }
    // 関数終了時に確実にファイルを閉じる
    defer dst.Close()

    // 4. ソースファイルの内容をデスティネーションファイルにコピー
    // io.Copy(dst, src) は、src から dst にバッファリングしながらデータをコピーする
    if _, err := io.Copy(dst, src); err != nil {
        return err // コピー中にエラーがあれば返す
    }

    // 5. コピーしたファイルのパーミッションを再設定
    // (既存ファイルを上書きする場合、OSや環境によってパーミッションが変化する可能性があるため)
    if err := os.Chmod(dstPath, srcInfo.Mode()); err != nil {
        return err
    }

    // 6. ファイルをディスクにフラッシュし、書き込みが完了したことを保証
    if err := dst.Sync(); err != nil {
        return err
    }

    // エラーがなければ nil を返して終了
    return nil
}

func main() {
    // 使用例
    // source.txt を destination.txt にコピーする
    if err := CopyFile("source.txt", "destination.txt"); err != nil {
        // エラーが発生した場合の処理
        panic(err)
    }
}