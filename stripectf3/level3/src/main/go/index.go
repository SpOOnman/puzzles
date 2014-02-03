package main

// http://blog.golang.org/go-maps-in-action
// http://golang.org/doc/codewalk/functions/
// http://www.golang-book.com/13

import (
    "fmt"
    "os"
    "bufio"
    "path/filepath"
    "time"
    "strings"
    "bytes"
	"net/http"
    "flag"
    "encoding/binary"
    "runtime/pprof"
    "runtime/debug"
)

const (
    Path = "/home/tkl/dev/stripectf3/level3/test/data/input"
)

type Position struct {
    fileIdx, row, col uint16
}

type Match struct {
    filename string
    line int
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

var working = false
var root = ""

var trigram2pos = make(map[uint32][]Position)
var filenames = make([]string, 0)

var filename_counter uint16 = 0
var filesize_counter int64 = 0
var position_counter = 0

func index(path string) {
    fmt.Printf("Start indexing path %s\n", path)
    root = fmt.Sprintf("%s/", path)

    startTime := time.Now()
    err := filepath.Walk(path, visit)
    fmt.Printf("filepath.Walk() returned %v took %s\n", err, time.Since(startTime))
}

func visit(path string, f os.FileInfo, err error) error {
    if !f.IsDir() {
        visitFile(path, f)
    }

    // Without this it doesn't free up memory.
    // Objects are GCed, but memory is not given back to kernel.
    debug.FreeOSMemory()

    return nil
}

func visitFile(path string, f os.FileInfo) {
    fmt.Printf("%v. File: %v, size: %v\n", filename_counter, path, f.Size() / 1024)
    filenames = append(filenames, path)
    filename_counter++
    filesize_counter += f.Size()
    parseFile(path)
    fmt.Printf("Total files: %v, size: %v, positions: %v, keys: %v\n", filename_counter, filesize_counter / 1024, position_counter, len(trigram2pos))
}

func parseFile(path string) {
    file, _ := os.Open(path)
    scanner := bufio.NewScanner(file)
    lineIdx := 0
    for scanner.Scan() {
        parseLine(scanner.Bytes(), lineIdx)
        lineIdx++
    }
}

//http://stackoverflow.com/questions/15117513/runtime-error-assignment-to-entry-in-nil-map
func parseLine(line []byte, lineIdx int) {
    length := len(line) - 2
    for i := 0 ; i < length ; i++ {
        trigram := line[i:i+3]
        if (trigram[0] == 0x20 || trigram[1] == 0x20 || trigram[2] == 0x20 ||
            trigram[0] == 0x2e || trigram[1] == 0x2e || trigram[2] == 0x2e) {
            continue
        }
        code := binary.LittleEndian.Uint32([]byte{line[i], line[i+1], line[i+2], 0})
        pos := Position{fileIdx: filename_counter, row: uint16(lineIdx), col: uint16(i)}
//        fmt.Printf("Trigram is '%s', position is %v", trigram, pos)

        trigram2pos[uint32(code)] = append(trigram2pos[uint32(code)], pos)
        position_counter++
    }
}

func query(text string) []Match {
    bs := []byte(text)
    length := len(bs)
    result := make([]Match, 0)
    trigrams := make([][]Position, length - 2)
    heads := make([]int, length - 2)
    var lastFoundFile uint16 = 9999
    var lastFoundRow uint16 = 9999
    for i := 0 ; i < length - 2 ; i++ {
//        buffer := bytes.NewBuffer(bs[i:i+3])
//        code, _ := binary.ReadVarint(buffer)
//        code := binary.BigEndian.Uint32([]byte{bs[i], bs[i+1], bs[i:i+3])
        code := binary.LittleEndian.Uint32([]byte{bs[i], bs[i+1], bs[i+2], 0})
        trigrams[i] = trigram2pos[uint32(code)]
//        fmt.Printf("Bytes %v, code %v, uint32 %v, trigram", bs[i:i+3], code, uint32(code), trigrams[i])
    }

//    fmt.Printf("START %v\n", text)
    for _, poss := range trigrams[0] {
//        fmt.Printf("We're looking for %v at %v\n", poss, strings.Replace(filenames[poss.fileIdx - 1], root, "", -1))
        if (poss.fileIdx == lastFoundFile && poss.row == lastFoundRow) {
//            fmt.Printf("Skipping matching on same file and line\n")
            continue
        }

        // move all pointer right
        for index, trigram := range trigrams {
//            h := heads[index]
//            fmt.Printf("Head index is %v.%v \n", index, h)
            pos := trigram[heads[index]]
//            fmt.Printf("Head index is %v at %v/%v points to %v\n", index, heads[index], len(trigram), h, pos)
            for ((pos.fileIdx < poss.fileIdx) ||
                (pos.fileIdx == poss.fileIdx && pos.row < poss.row) ||
                (pos.fileIdx == poss.fileIdx && pos.row == poss.row && pos.col < poss.col + uint16(index))) &&
                heads[index] < len(trigram) -1 {
                heads[index]++
                pos = trigram[heads[index]]
//                fmt.Printf("Move right to %v/%v: %v\n",heads[index], len(trigram), pos)
            }
        }
        // check if all pointer's values create a valid chain
        found := true
        for index, trigram := range trigrams {
            pos := trigram[heads[index]]
            if pos.fileIdx != poss.fileIdx || pos.row != poss.row || pos.col != poss.col + uint16(index) {
                found = false
                break
            }
        }
        if found {
            result = append(result, Match{filename: strings.Replace(filenames[poss.fileIdx - 1], root, "", -1), line: int(poss.row) + 1})
            lastFoundFile = poss.fileIdx
            lastFoundRow = poss.row
        }
    }
    return result
}


//func main() {
//    index(Path)
//    q := "ramshackly"
//    fmt.Printf("Entries for %v: %v", q, query(q))
//    main2()
//}

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, _ := os.Create(*cpuprofile)
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
    http.HandleFunc("/healthcheck", healthcheckHandler)
    http.HandleFunc("/index", indexHandler)
    http.HandleFunc("/isIndexed", isIndexedHandler)
    http.HandleFunc("/", queryHandler)
    http.ListenAndServe(":9090", nil)
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "{\"success\": \"true\"}")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    working = true
    index(r.URL.Query()["path"][0])
    working = false
//    fmt.Printf("Positions for eek, %v", trigram2pos["eek"])
    if *memprofile != "" {
        f, _ := os.Create(*memprofile)
        pprof.WriteHeapProfile(f)
        f.Close()
    }
}

func isIndexedHandler(w http.ResponseWriter, r *http.Request) {
    if working {
        fmt.Fprintf(w, "{\"success\": \"false\"}")
    } else {
        fmt.Fprintf(w, "{\"success\": \"true\"}")
    }
}

//http://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go
func queryHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()["q"][0]
    startTime := time.Now()
    matches := query(q)
    delta := time.Since(startTime)
    fmt.Printf("Entries for %v: found %v in %v\n", q, len(matches), delta)

    var buffer bytes.Buffer
    first := true

    buffer.WriteString("{ \"success\": true, \"results\": [")
    for _, match := range matches {
        if !first {
            buffer.WriteString(",")
        }
        buffer.WriteString(fmt.Sprintf("\"%s:%d\"", match.filename, match.line))
        first = false
    }
    buffer.WriteString("]}")

    fmt.Fprintf(w, buffer.String())
}
