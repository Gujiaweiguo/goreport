package datasource

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHTunnel struct {
	client     *ssh.Client
	forwarder  net.Listener
	localAddr  string
	remoteAddr string
	closeOnce  sync.Once
}

type SSHTunnelConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Key      []byte
	Phrase   string
}

func NewSSHTunnel(config *SSHTunnelConfig) *SSHTunnel {
	return &SSHTunnel{
		localAddr:  "127.0.0.1:0",
		remoteAddr: fmt.Sprintf("127.0.0.1:%d", config.Port),
	}
}

func (t *SSHTunnel) Connect(ctx context.Context, config *SSHTunnelConfig, targetHost string, targetPort int) (string, error) {
	authMethods := []ssh.AuthMethod{}

	if config.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	if len(config.Key) > 0 {
		var signer ssh.Signer
		var err error

		if config.Phrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(config.Key, []byte(config.Phrase))
		} else {
			signer, err = ssh.ParsePrivateKey(config.Key)
		}

		if err != nil {
			return "", fmt.Errorf("failed to parse SSH private key: %w", err)
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return "", errors.New("no SSH authentication method provided")
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	dialCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	type dialResult struct {
		client *ssh.Client
		err    error
	}

	resultChan := make(chan dialResult, 1)

	go func() {
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), sshConfig)
		resultChan <- dialResult{client: client, err: err}
	}()

	select {
	case <-dialCtx.Done():
		return "", fmt.Errorf("SSH connection timeout: %w", dialCtx.Err())
	case result := <-resultChan:
		if result.err != nil {
			return "", fmt.Errorf("failed to establish SSH connection: %w", result.err)
		}
		t.client = result.client
	}

	listener, err := t.client.Listen("tcp", t.remoteAddr)
	if err != nil {
		t.client.Close()
		return "", fmt.Errorf("failed to listen on SSH remote address: %w", err)
	}

	t.forwarder = listener

	localListener, err := net.Listen("tcp", t.localAddr)
	if err != nil {
		listener.Close()
		t.client.Close()
		return "", fmt.Errorf("failed to listen on local address: %w", err)
	}

	t.localAddr = localListener.Addr().String()

	go t.forwardConnections(ctx, listener, localListener, targetHost, targetPort)

	return t.localAddr, nil
}

func (t *SSHTunnel) forwardConnections(ctx context.Context, remote net.Listener, local net.Listener, targetHost string, targetPort int) {
	defer remote.Close()
	defer local.Close()
	defer t.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			localConn, err := local.Accept()
			if err != nil {
				return
			}

			remoteConn, err := remote.Accept()
			if err != nil {
				localConn.Close()
				return
			}

			go func() {
				defer localConn.Close()
				defer remoteConn.Close()

				done := make(chan struct{}, 2)

				go func() {
					defer close(done)
					_, _ = copyData(localConn, remoteConn)
				}()

				go func() {
					defer close(done)
					_, _ = copyData(remoteConn, localConn)
				}()

				<-done
				<-done
			}()
		}
	}
}

func copyData(dst net.Conn, src net.Conn) (written int64, err error) {
	defer dst.Close()
	defer src.Close()

	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = errors.New("short write")
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func (t *SSHTunnel) Close() error {
	var errs []error

	t.closeOnce.Do(func() {
		if t.forwarder != nil {
			if err := t.forwarder.Close(); err != nil {
				errs = append(errs, err)
			}
		}
		if t.client != nil {
			if err := t.client.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	})

	if len(errs) > 0 {
		return fmt.Errorf("errors closing SSH tunnel: %v", errs)
	}
	return nil
}

func (t *SSHTunnel) LocalAddr() string {
	return t.localAddr
}
