package MySQL

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func copyHeader(dst, src http.Header) {
	for k, w := range src {
		for _, v := range w {
			dst.Add(k, v)
		}
	}
}

func copyResponse(r *http.Response, w http.ResponseWriter) {
	copyHeader(w.Header(), r.Header)
	w.WriteHeader(r.StatusCode)
	io.Copy(w, r.Body)
}



func AppHandler(c *gin.Context) {
	var (
		resp *http.Response
		err  error
	)
	var path string
	bindID := c.Param("id")
	path = strings.TrimPrefix(c.Request.RequestURI, "/app/"+bindID)
	var searchStruct RunningStruct
	SqlDB.Where(&RunningStruct{BindID: bindID }).First(&searchStruct)
	host := searchStruct.ClusterIP + ":" + strconv.Itoa(int(searchStruct.Port))
	log.Println(host)
	log.Println(path)
	if path == "/ws" {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		var header string
		for k, v := range c.Request.Header {
			header += k + ": " + strings.Join(v, ";") + "\n"
		}
		if _, err = conn.Write([]byte("GET /ws HTTP/1.1\nHost: " + host + "\n" + header + "\n")); err != nil {
			fmt.Println(err)
			return
		}
		tcpConn, _, err := c.Writer.Hijack()
		if err != nil {
			fmt.Println(tcpConn)
			return
		}
		defer tcpConn.Close()
		go io.Copy(conn, tcpConn)
		io.Copy(tcpConn, conn)
	} else {
		resp, err = http.Get("http://" + host + path)
		if err != nil {
			// handle error
			fmt.Println(err)
			return
		}
		copyResponse(resp, c.Writer)
	}
}
