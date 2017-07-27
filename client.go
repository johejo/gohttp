package main

import (
    "fmt"
    "net/url"
    "os"
    "net"
    "net/http"
    "bufio"
    "strconv"
)

const BUF_LEN int = 256

func main() {
    //fmt.Println("Hello! This is Golang HTTP client.")
    var err error
    length := 0
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
        os.Exit(1)
    }

    u, err := url.Parse(os.Args[1])
    if !check_error(err, "URL parsing") {
        os.Exit(1)
    }

    var conn net.Conn = nil
    if conn == nil {
        conn, err = net.Dial("tcp", u.Host)
        if !check_error(err, "create connection") {
            os.Exit(1)
        }
        defer conn.Close()
        length = get_content_length(conn, u.String())
    }

    result := get_and_write(conn, *u, length)
    if !result {
        fmt.Fprint(os.Stderr, "get and write failed")
    }
    //fmt.Println(string(all))
}

func check_error(err error, memo string, n ...int) bool {
    if err != nil {
        fmt.Fprintf(os.Stderr, memo+":"+err.Error())
        return false
    }

    if n != nil {
        for _, i := range []int(n) {
            fmt.Println(memo, i)
            if i == -1 {
                fmt.Fprintf(os.Stderr, memo+":"+err.Error())
                return false
            }
        }

    }

    return true
}


func get_and_write(conn net.Conn, u url.URL, length int) bool {
    var err error
    var method string = "GET"

    request, err := http.NewRequest(method, u.String(), nil)
    if !check_error(err, "set "+method+" request") {
        return false
    }

    err = request.Write(conn)
    if !check_error(err, "write request to socket") {
        return false
    }

    response, _ := http.ReadResponse(bufio.NewReader(conn), request)
    if !check_error(err, "read response from socket") {
        return false
    }

    total := 0
    //fmt.Println(length)
    var data []byte
    for {
        buf := make([]byte, BUF_LEN)
        response.Body.Read(buf)
        //fmt.Print(string(buf))
        total += len(buf)
        data = append(data, buf...)
        if total >= length {
            break
        }
    }
    filename := "."+u.Path
    fp, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0666)
    if !check_error(err, "file open") {
        return false
    } else {
        defer fp.Close()
    }
    fp.Write(data)
    return true
}

func get_content_length(conn net.Conn, u string) int {
    var err error
    method := "HEAD"
    length := 0

    request, err := http.NewRequest(method, u, nil)
    if !check_error(err, "set "+method+" request") {
        return -1
    }

    err = request.Write(conn)
    if !check_error(err, "write request to socket") {
        return -1
    }

    response , err := http.ReadResponse(bufio.NewReader(conn), request)
    if !check_error(err, "read response from socket") {
        return -1
    }

    length, err = strconv.Atoi(response.Header.Get("Content-Length"))
    if !check_error(err, "convert 'content-length' to integer") {
        return -1
    }

    return length
}