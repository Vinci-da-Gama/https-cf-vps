package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	utls "github.com/refraction-networking/utls"
)

const (
	minPort     = 1
	maxPort     = 65535
	minChunk    = 1
	maxChunk    = 1024
	dialTimeout = 30 * time.Second
)

var (
	errLogger     *log.Logger
	loggerOnce    sync.Once
	port          = flag.Int("port", 8080, "HTTP代理端口,可选值:1-65535")
	passwd        = flag.String("pwd", "testPASSword", "密码")
	wssHost       = flag.String("wss", "", "websocket地址,[域名]:[端口](非443)")
	ckSize        = flag.Int("chunk", 64, "websocket每一帧的数据大小(KB),可选值:1-1024")
	debug         = flag.Bool("debug", false, "是否输出调试信息")
	wssHostRegexp = regexp.MustCompile(`^[a-zA-Z0-9.-]+(:\d+)?(/.*)?$`)
)

func Debug(err error) {
	if *debug && err != nil {
		loggerOnce.Do(func() {
			errLogger = log.New(os.Stderr, "\033[31m[ERROR]\033[0m ", log.LstdFlags|log.Lshortfile)
		})
		errLogger.Println(err)
	}
}

func PipeConn(ws *websocket.Conn, conn net.Conn) {
	defer ws.Close()
	defer conn.Close()

	buf := make([]byte, *ckSize*1024)
	done := make(chan struct{})

	// Websocket to TCP
	go func() {
		defer close(done)
		for {
			mt, r, err := ws.NextReader()
			if err != nil {
				Debug(err)
				return
			}
			if mt != websocket.BinaryMessage {
				io.Copy(io.Discard, r)
				continue
			}
			if _, err := io.CopyBuffer(conn, r, buf); err != nil {
				Debug(err)
				return
			}
		}
	}()

	// TCP to Websocket
	for {
		select {
		case <-done:
			return
		default:
			n, err := conn.Read(buf)
			if err != nil {
				Debug(err)
				return
			}
			if n > 0 {
				if err := ws.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					Debug(err)
					return
				}
			}
		}
	}
}

func utlsDialTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	var d net.Dialer
	tcpConn, err := d.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	if host == "" {
		host = addr
	}

	uconn := utls.UClient(tcpConn, &utls.Config{ServerName: host}, utls.HelloRandomized)

	if dl, ok := ctx.Deadline(); ok {
		uconn.SetDeadline(dl)
	}

	if err := uconn.Handshake(); err != nil {
		tcpConn.Close()
		return nil, err
	}

	uconn.SetDeadline(time.Time{})
	return uconn, nil
}

func SetUpTunnel(ctx context.Context, client net.Conn, target string) {
	defer client.Close()

	header := http.Header{
		"X-Target":   []string{target},
		"X-Password": []string{*passwd},
	}

	dialer := websocket.Dialer{
		NetDialTLSContext: utlsDialTLSContext,
		HandshakeTimeout:  dialTimeout,
	}

	ws, resp, err := dialer.DialContext(ctx, "wss://"+*wssHost, header)
	if err != nil {
		Debug(err)
		if resp != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			log.Printf("连接websocket出错: %s", body)
		}
		return
	}

	PipeConn(ws, client)
}

func main() {
	flag.Parse()

	if !wssHostRegexp.MatchString(*wssHost) {
		log.Fatalln("websocket地址,[域名]:[端口](非443)")
	}

	if *port < minPort || *port > maxPort {
		log.Fatalf("HTTP代理端口,可选值:%d-%d", minPort, maxPort)
	}

	if *ckSize < minChunk || *ckSize > maxChunk {
		log.Fatalf("websocket每一帧的数据大小(KB),可选值:%d-%d", minChunk, maxChunk)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodConnect {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "不支持 Hijacking", http.StatusInternalServerError)
			return
		}

		client, _, err := hijacker.Hijack()
		if err != nil {
			Debug(err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		log.Printf("访问: %s", r.Host)

		ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
		defer cancel()

		go SetUpTunnel(ctx, client, r.Host)
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("开启HTTP代理,端口:%d", *port)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal("开启HTTP代理失败:", err)
	}
}
