package main

import (
    "fmt"
    "net/url"
    "os"
    "net"
    "net/http"
    "bufio"
    "strings"
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
        length = get_content_length(conn, u.String())
    }

    data := get(conn, u.String(), length)
    if data == nil {
        fmt.Fprint(os.Stderr, "get data failed")
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


func get(conn net.Conn, u string, length int) []byte {
    var err error
    var bodies string = ""
    var method string = "GET"

    request, err := http.NewRequest(method, u, strings.NewReader(bodies))
    if !check_error(err, "set "+method+" request") {
        return nil
    }

    err = request.Write(conn)
    if !check_error(err, "write request to socket") {
        return nil
    }

    response, _ := http.ReadResponse(bufio.NewReader(conn), request)
    if !check_error(err, "read response from socket") {
        return nil
    }

    total := 0
    //fmt.Println(length)
    var data []byte
    for {
        buf := make([]byte, BUF_LEN)
        response.Body.Read(buf)
        fmt.Print(string(buf))
        total += len(buf)
        data = append(data, buf...)
        if total >= length {
            break
        }
    }
    return data
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

    length, _ = strconv.Atoi(response.Header.Get("Content-Length"))

    return length
}