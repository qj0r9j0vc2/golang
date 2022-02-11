package echo

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestEchoServerUnix(t *testing.T) {
	dir, err := ioutil.TempDir("", "echo_unix") // 임시 디렉터리에 하위 디렉터리 생성
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil { //프로그램 종료시 파일 삭제
			t.Error(rErr)
		}
	}()

	ctx, _ := context.WithCancel(context.Background()) // 컨텍스트 생성
	defer func() {
		ctx.Done()
	}()

	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid())) // 생성했던 디렉터리에 socket File 저장
	rAddr, err := streamingEchoServer(ctx, "unix", socket)            //server 작동
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chmod(socket, os.ModeSocket|0666)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.Dial("unix", rAddr.String()) // rAddr 주소와 unix network type으로 연결
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = conn.Close()
	}()

	msg := []byte("ping")
	for i := 0; i < 3; i++ {
		_, err := conn.Write(msg) //msg 3번 전송
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf) // msg Read
	if err != nil {
		t.Fatal(err)
	}

	expected := bytes.Repeat(msg, 3) //msg 3번 들어왔는지 확인
	if !bytes.Equal(expected, buf[:n]) {
		t.Fatalf("expected reply %q; actual reply %q", expected, buf[:n])
	}

}
