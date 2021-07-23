package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "strings"
        "time"
)

// port configs
var tcpPort  = 9091
var httpPort = 9090


// our temporary storage
var dataBuffer = make(map[string]byte)




func main() {
    // start the controle port
    startControlePort()
    
    // open tcp server
    l, err := net.Listen("tcp", tcpPort)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer l.Close()

    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
    
}    
    
func startControlePort() {
    // run a http port
    h := http.NewServeMux()
    
	// Route: getIdent
	h.HandleFunc("/getIdent", func(w http.ResponseWriter, r *http.Request) {
        var uuid = ""
        var check = true
        for true == check {
            // get uuid
            uuid,err = getUUID()
            if err != nil {
                respond("error generating UUID", 500, w)
                return
            }
            
            // lets check if the uuid is already used,
            // if not we gonne create the uuid with a 0byte array
            if _, ok := dataBuffer[uuid]; ok {
                dataBuffer[uuid] = make([]byte,0)
                check = false
            }
            
        }
        
        // k we got a UUID lets return it
		respond(uuid, 200, w)
	})
    
    h.HandleFunc("/getData", func(w http.ResponseWriter, r *http.Request) {
        // prepare ret string
        ret := ""
        // lets get the uri qry params
        tmpParams := r.URL.Query()
        // read the uuid from params
        val, ok := tmpParams["uuid"]
        // if uuid param doesnt exist
        if !ok {
            respond("missing params",500,w)
            return
        }
    })
    
    return
}

func handleConnection(conn net.Conn){
    // first we read the first 10 bytes to determine a identifier delimnited by \n
    uuid, err := bufio.NewReader(buf).ReadString('\n')
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }
    // Send a response back to person contacting us.
    conn.Write([]byte("Message received."))
    // Close the connection when you're done with it.
    conn.Close()
}

func getUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}


func respond(message string, responseCode int, w http.ResponseWriter) {
	w.WriteHeader(responseCode)
	messageBytes := []byte(message)
	_, err := w.Write(messageBytes)
	if nil != err {
		config.Logger.Print(err)
	}
}

