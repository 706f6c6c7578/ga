package main

import (
    "bufio"
    "flag"
    "fmt"
    "net"
    "strings"
)

func getArticle(server string, port int, msgID string, fullArticle bool) (string, error) {
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
    if err != nil {
        return "", err
    }
    defer conn.Close()

    reader := bufio.NewReader(conn)
    _, err = reader.ReadString('\n')
    if err != nil {
        return "", err
    }

    fmt.Fprintf(conn, "article %s\r\n", msgID)

    // Read and discard the server's response to the article command
    _, err = reader.ReadString('\n')
    if err != nil {
        return "", err
    }

    var builder strings.Builder
    inBody := false
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            return "", err
        }

        if line == ".\r\n" {
            break
        }

        if fullArticle || inBody {
            builder.WriteString(line)
        } else if line == "\r\n" {
            inBody = true
        }
    }

    fmt.Fprintf(conn, "quit\r\n")

    return builder.String(), nil
}

func main() {
    fullArticle := flag.Bool("f", false, "Retrieve full article including headers")
    server := flag.String("s", "news.i2pn2.org", "NNTP server address")
    port := flag.Int("p", 119, "NNTP server port")
    flag.Parse()

    args := flag.Args()
    if len(args) < 1 {
        fmt.Println("Please provide a message ID")
        fmt.Println("Usage: ga [-f] [-s server] [-p port] <message-id>")
        return
    }

    msgID := args[0]

    if !strings.HasPrefix(msgID, "<") || !strings.HasSuffix(msgID, ">") {
        msgID = "<" + msgID + ">"
    }

    article, err := getArticle(*server, *port, msgID, *fullArticle)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Print(article)
}

