package sftpUtils

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/wuchunfu/CloudSync/config"
	"github.com/wuchunfu/CloudSync/middleware/logUtils"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"
)

// SftpHandler 全局结构体
type SftpHandler struct {
	SftpClient *sftp.Client
}

// NewSftpHandler 初始化
func NewSftpHandler() *SftpHandler {
	sftpHandler := new(SftpHandler)
	sftpHandler.SftpClient, _ = connect(config.GlobalObject.Sftp.Hostname, config.GlobalObject.Sftp.SSHPort, config.GlobalObject.Sftp.Username, config.GlobalObject.Sftp.Password)
	return sftpHandler
}

// connect 生成链接 对象
func connect(host string, port int, username string, password string) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    username,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

// uploadFile 首先上传文件的方法
func (sftpHandler *SftpHandler) uploadFile(localPath string, remotePath string) {
	// 打开本地文件流
	srcFile, err := os.Open(localPath)
	if err != nil {
		fmt.Println("os.Open error : ", localPath)
		log.Fatal(err)
	}
	// 关闭文件流
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(srcFile)

	// 上传到远端服务器的文件名,与本地路径末尾相同
	var remoteFileName = path.Base(localPath)

	// 判断当前目录是否存在
	if _, err := sftpHandler.SftpClient.Stat(remotePath); err != nil {
		mkdirErr := sftpHandler.SftpClient.Mkdir(remotePath)
		if mkdirErr != nil {
			log.Fatal(mkdirErr)
			return
		}
	}

	//打开远程文件,如果不存在就创建一个
	dstFile, err := sftpHandler.SftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)
	}
	// 关闭远程文件
	defer func(dstFile *sftp.File) {
		err := dstFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dstFile)

	// 读取本地文件,写入到远程文件中(这里没有分块传, 自己写的话可以改一下,防止内存溢出)
	readFile, err := io.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localPath)
		log.Fatal(err)
	}
	_, writeErr := dstFile.Write(readFile)
	if writeErr != nil {
		log.Fatal(writeErr)
		return
	}
	logUtils.Logger.Println("上传：" + localPath)
}

//uploadDirectory 遍历上传远程文件夹
func (sftpHandler *SftpHandler) uploadDirectory(localPath string, remotePath string) {
	//打开本地文件夹流
	localFiles, err := os.ReadDir(localPath)
	if err != nil {
		log.Fatal("路径错误 ", err)
	}
	// 先创建最外层文件夹
	mkdirErr := sftpHandler.SftpClient.Mkdir(remotePath)
	if mkdirErr != nil {
		log.Fatal(mkdirErr)
		return
	}
	// 遍历文件夹内容
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		// 判断是否是文件,是文件直接上传.是文件夹,先远程创建文件夹,再递归复制内部文件
		if backupDir.IsDir() {
			err := sftpHandler.SftpClient.Mkdir(remoteFilePath)
			if err != nil {
				log.Println(err)
				return
			}
			sftpHandler.uploadDirectory(localFilePath, remoteFilePath)
		} else {
			sftpHandler.uploadFile(path.Join(localPath, backupDir.Name()), remotePath)
		}
	}

	logUtils.Logger.Println("上传本地目录：" + localPath + "远端目录：" + remotePath)
}

// Upload 判断是否是路径属性
func (sftpHandler *SftpHandler) Upload(localPath string, remotePath string) {
	// 获取路径的属性
	s, err := os.Stat(localPath)
	if err != nil {
		fmt.Println("文件路径不存在", err)
		return
	}

	// 判断是否是文件夹
	if s.IsDir() {
		sftpHandler.uploadDirectory(localPath, remotePath)
	} else {
		sftpHandler.uploadFile(localPath, remotePath)
	}
}
