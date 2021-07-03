package server

import (
	"echoServer/models"
	"echoServer/server/connectionQueue"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestEchoServer_Listen(t *testing.T) {
	type fields struct {
		listener     net.Listener
		connQueue    connectionQueue.ConnectionQueue
		idleTimeout  time.Duration
		maxConns     int
		maxReadBytes int
		currentConns int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*s := &EchoServer{
				listener:     tt.fields.listener,
				connQueue:    tt.fields.connQueue,
				idleTimeout:  tt.fields.idleTimeout,
				maxConns:     tt.fields.maxConns,
				maxReadBytes: tt.fields.maxReadBytes,
				currentConns: tt.fields.currentConns,
			}*/
		})
	}
}

func TestEchoServer_closeLeastUpdConn(t *testing.T) {
	type fields struct {
		listener     net.Listener
		connQueue    connectionQueue.ConnectionQueue
		idleTimeout  time.Duration
		maxConns     int
		maxReadBytes int
		currentConns int
	}
	tests := []struct {
		name   string
		fields fields
		client net.Conn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/*s := &EchoServer{
				listener:     tt.fields.listener,
				connQueue:    tt.fields.connQueue,
				idleTimeout:  tt.fields.idleTimeout,
				maxConns:     tt.fields.maxConns,
				maxReadBytes: tt.fields.maxReadBytes,
				currentConns: tt.fields.currentConns,
			}*/
		})
	}
}

func TestEchoServer_handleRequest(t *testing.T) {
	type fields struct {
		listener     net.Listener
		connQueue    connectionQueue.ConnectionQueue
		idleTimeout  time.Duration
		maxConns     int
		maxReadBytes int
		currentConns int
	}
	type args struct {
		c *models.Connection
	}
	l, _ := net.Listen("tcp", "127.0.0.1:3333")

	client, _ := net.Dial("tcp", "127.0.0.1:3333")
	tests := []struct {
		name   string
		fields fields
		args   args
		client net.Conn
		want   string
		sleep  time.Duration
		msg    string
	}{
		{
			name: "normal read",
			fields: fields{
				listener:     l,
				connQueue:    connectionQueue.NewConnectionQueue(),
				idleTimeout:  100 * time.Second,
				maxConns:     3,
				maxReadBytes: 8192,
				currentConns: 0,
			},
			client: client,
			want:   "Hello",
			sleep:  time.Millisecond,
			msg:    "Hello",
		},
		{
			name: "read above max length",
			fields: fields{
				listener:     l,
				connQueue:    connectionQueue.NewConnectionQueue(),
				idleTimeout:  100 * time.Second,
				maxConns:     3,
				maxReadBytes: 10,
				currentConns: 0,
			},
			client: client,
			want:   "HelloHello",
			sleep:  time.Millisecond,
			msg:    "HelloHelloABCDEFG",
		},
		{
			name: "timeout on first read",
			fields: fields{
				listener:     l,
				connQueue:    connectionQueue.NewConnectionQueue(),
				idleTimeout:  3 * time.Second,
				maxConns:     3,
				maxReadBytes: 8192,
				currentConns: 0,
			},
			client: client,
			want:   "Exit due to idle.\n",
			sleep:  time.Second * 3,
		},
		{
			name: "quit signal",
			fields: fields{
				listener:     l,
				connQueue:    connectionQueue.NewConnectionQueue(),
				idleTimeout:  100 * time.Second,
				maxConns:     3,
				maxReadBytes: 8192,
				currentConns: 0,
			},
			client: client,
			want:   "Got quit signal. Aborting.\n",
			sleep:  time.Millisecond,
			msg:    "quit",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EchoServer{
				listener:     tt.fields.listener,
				connQueue:    tt.fields.connQueue,
				idleTimeout:  tt.fields.idleTimeout,
				maxConns:     tt.fields.maxConns,
				maxReadBytes: tt.fields.maxReadBytes,
				currentConns: tt.fields.currentConns,
			}
			sConn, _ := s.listener.Accept()

			go s.handleRequest(&models.Connection{
				Conn:       sConn,
				LastUpdate: time.Now(),
				Index:      0,
			})

			time.Sleep(tt.sleep)
			tt.client.Write([]byte(tt.msg))
			buf := make([]byte, tt.fields.maxReadBytes)
			n, _ := tt.client.Read(buf)
			assert.Equal(t, tt.want, string(buf[:n]))
		})
	}
}

func TestEchoServer_isQuit(t *testing.T) {
	type fields struct {
		listener     net.Listener
		connQueue    connectionQueue.ConnectionQueue
		idleTimeout  time.Duration
		maxConns     int
		maxReadBytes int
		currentConns int
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "is quit signal",
			fields: fields{},
			args:   args{data: []byte("quit")},
			want:   true,
		},
		{
			name:   "is not a quit signal",
			fields: fields{},
			args:   args{data: []byte("not a quit")},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &EchoServer{
				listener:     tt.fields.listener,
				connQueue:    tt.fields.connQueue,
				idleTimeout:  tt.fields.idleTimeout,
				maxConns:     tt.fields.maxConns,
				maxReadBytes: tt.fields.maxReadBytes,
				currentConns: tt.fields.currentConns,
			}
			got := s.isQuit(tt.args.data)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewEchoServer(t *testing.T) {
	type args struct {
		host string
		port string
	}
	tests := []struct {
		name    string
		args    args
		want    *EchoServer
		wantErr error
	}{
		{
			name: "empty server with no settings",
			args: args{
				host: "",
				port: "",
			},
			want:    nil,
			wantErr: models.ErrCreateServer,
		},
		{
			name: "empty server with only port",
			args: args{
				host: "",
				port: "8080",
			},
			want:    nil,
			wantErr: models.ErrCreateServer,
		},
		/*{*/
		/*	name: "normal server with all params",*/
		/*	args: args{*/
		/*		host: "127.0.0.1",*/
		/*		port: ":8080",*/
		/*	},*/
		/*	want: &EchoServer{*/
		/*		connQueue:    connectionQueue.NewConnectionQueue(),*/
		/*		idleTimeout:  30 * time.Second,*/
		/*		maxConns:     500000,*/
		/*		maxReadBytes: 8192,*/
		/*		currentConns: 0,*/
		/*	},*/
		/*},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEchoServer(tt.args.host, tt.args.port)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
