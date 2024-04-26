// 主程序包
package main

// 导入所需的库
import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net"
	"strings"
	"sync"
)

// 使用go:embed将Dic.txt文件嵌入到二进制文件中
//
//go:embed Dic.txt
var dic string

// wg1和wg2是WaitGroup对象，用于等待goroutine完成
var wg1 sync.WaitGroup
var wg2 sync.WaitGroup

//常见端口
/*
21端口：FTP 文件传输服务
22端口：SSH协议、SCP（文件传输）、端口号重定向
23/tcp端口：TELNET 终端仿真服务
25端口：SMTP 简单邮件传输服务
53端口：DNS 域名解析服务
69/udp：TFTP
80/8080/3128/8081/9098端口：HTTP协议代理服务器
110/tcp端口：POP3（E-mail）
119端口：Network
123端口：NTP（网络时间协议）
135、137、138、139端口： 局域网相关默认端口，应关闭
161端口：SNMP（简单网络管理协议）
389端口：LDAP（轻量级目录访问协议）、ILS（定位服务）
443/tcp 443/udp：HTTPS服务器
465端口：SMTP（简单邮件传输协议）
873端口：rsync
1080端口：SOCKS代理协议服务器常用端口号、QQ
1158端口：ORACLE EMCTL
1433/tcp/udp端口：MS SQL*SERVER数据库server、MS SQL*SERVER数据库monitor
1521端口：Oracle 数据库
2100端口：Oracle XDB FTP服务
3389端口：WIN2003远程登录
3306端口：MYSQL数据库端口
5432端口：postgresql数据库端口
5601端口：kibana
6379端口：Redis数据库端口
8080端口：TCP服务端默认端口、JBOSS、TOMCAT、Oracle XDB（XML 数据库）
8081端口：Symantec AV/Filter for MSE
8888端口：Nginx服务器的端口
9000端口：php-fpm
9080端口：Webshpere应用程序
9090端口：webshpere管理工具
9200端口：Elasticsearch服务器端口
10050端口：zabbix_server 10050
10051端口：zabbix_agent
11211端口：memcache（高速缓存系统）
27017端口：mongoDB数据库默认端口
22122端口：fastdfs服务器默认端口
*/

// scanner1函数是端口扫描的实现
func scanner1(wg sync.WaitGroup, host string, ports chan int) {
	//ch := make(chan string, 1024)
	//函数结束时,计数器减1
	defer wg.Done()
	//管道ports导出到port
	port := <-ports
	//拼接目标
	goal := fmt.Sprintf("%s:%d", host, port)
	//开始扫描
	_, err := net.Dial("tcp", goal)
	//错误处理
	if err != nil {
		return
	}
	sprintf := fmt.Sprintf("%s:%d", "[+]端口存活", port)
	fmt.Println(sprintf)
}

// scanner2函数用于检查HTTP/HTTPS服务是否存活
func scanner2(ch chan string) {
	for x := range ch {
		// 声明并初始化 HTTP 和 HTTPS 的完整 URL
		http := fmt.Sprintf("http://%s", x)
		https := fmt.Sprintf("https://%s", x)

		// 发起 HTTP 请求
		request := gorequest.New()
		httpresp, _, errors := request.Head(http).End()
		if errors != nil {
			goto HTTPS // 如果 HTTP 请求发生错误，则跳转到 HTTPS 请求
		}
		// 如果 HTTP 响应状态码为 200 或 403，则打印存活消息
		if httpresp.StatusCode == 200 || httpresp.StatusCode == 403 {
			fmt.Println("[+] 存活", http)
		}

	HTTPS:
		// 发起 HTTPS 请求
		httpsresp, _, i := request.Head(https).End()
		if i != nil {
			continue // 如果 HTTPS 请求发生错误，则跳过当前循环
		}
		// 如果 HTTPS 响应状态码为 200 或 403，则打印存活消息
		if httpsresp.StatusCode == 200 || httpresp.StatusCode == 403 {
			fmt.Println("[+] 存活", https)
		}

	}
}

// 接收参数结构体
type C struct {
	goal  string
	model string
}

// main函数是程序的入口点
func main() {
	c := C{}
	flag.StringVar(&c.goal, "u", "www.baidu.com", "目标")
	flag.StringVar(&c.model, "m", "1:常用端口扫描, 2:全端口扫描, 3:子域名扫描", "模式")
	flag.Parse()

	switch c.model {
	case "1":
		ints := make(chan int, 100)
		defer close(ints)
		ports := []int{21, 22, 23, 25, 53, 69, 80, 8080, 3128, 9098, 110, 119, 123, 135, 137, 138, 139, 161, 389, 443, 465, 873, 1080, 1158, 1433, 1521, 2100, 3389, 3306, 5432, 5601, 6379, 8081, 8888, 9000, 9080, 9090, 9200, 10050, 10051, 11211, 27017, 22122}
		for _, port := range ports {
			wg1.Add(1)
			ints <- port
			go scanner1(wg1, string(c.goal), ints)
		}
		wg1.Wait()
		fmt.Println("扫描完成")
		return
	case "2":
		ints := make(chan int, 100)
		defer close(ints)
		for i := 0; i < 65535; i++ {
			wg1.Add(1)
			ints <- i
			go scanner1(wg1, string(c.goal), ints)
		}
		wg1.Wait()
		fmt.Println("扫描完成")
		return
	case "3":
		ch := make(chan string, 1024)
		go scanner2(ch)
		var buffer bytes.Buffer
		buffer.WriteString(dic)
		for {
			line, err := buffer.ReadString('\n')
			if err != nil {
				break
			}
			wg2.Add(1)
			line = strings.TrimSpace(line)
			sprintf := fmt.Sprintf("%s%s%s", line, ".", string(c.goal))
			ch <- sprintf
		}
		wg2.Wait()
	}
}
