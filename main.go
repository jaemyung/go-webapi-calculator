package main
 
import (
    "fmt"
    "net/http"
    "strings"
    "strconv"
    "math"
)

type Context struct {
    Params map[string]interface{}
     
    ResponseWriter http.ResponseWriter
    Request        *http.Request
}

type HandlerFunc func(*Context)

type router struct {
    // 키: http 메서드
    // 값: URL 패턴별로 실행할 HandlerFunc
    handlers map[string]map[string]HandlerFunc
}
 
func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {
    // http 메서드로 등록된 맵이 있는지 확인
    m, ok := r.handlers[method]
    if !ok {
        // 등록된 맵이 없으면 새 맵을 생성
        m = make(map[string]HandlerFunc)
        r.handlers[method] = m
    }
    // http 메서드로 등록된 맵에 URL 패턴과 핸들러 함수 등록
    m[pattern] = h
}

func match(pattern, path string) (bool, map[string]string) {
    // 패턴과 패스가 정확히 일치하면 바로 true를 반환
    if pattern == path {
        return true, nil
    }
 

    // 패턴과 패스를 “/" 단위로 구분
    patterns := strings.Split(pattern, "/")
    paths := strings.Split(path, "/")
 

    // 패턴과 패스를 “/“로 구분한 후 부분 문자열 집합의 개수가 다르면 false를 반환
    if len(patterns) != len(paths) {
        return false, nil
    }
 

    // 패턴에 일치하는 URL 매개변수를 담기 위한 params 맵 생성
    params := make(map[string]string)
 

    // “/“로 구분된 패턴/패스의 각 문자열을 하나씩 비교
    for i := 0; i < len(patterns); i++ {
        switch {
        case patterns[i] == paths[i]:
            // 패턴과 패스의 부분 문자열이 일치하면 바로 다음 루프 수행
        case len(patterns[i]) > 0 && patterns[i][0] == ':':
            // 패턴이 ‘:’ 문자로 시작하면 params에 URL params를 담은 후 다음 루프 수행
            params[patterns[i][1:]] = paths[i]
        default:
            // 일치하는 경우가 없으면 false를 반환
            return false, nil
        }
    }
 

    // true와 params를 반환
    return true, params
}
 
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // http 메서드에 맞는 모든 handers를 반복하면서 요청 URL에 해당하는 handler를 찾음
    for pattern, handler := range r.handlers[req.Method] {
        if ok, params := match(pattern, req.URL.Path); ok {
            // Context 생성
            c := Context{
                Params:         make(map[string]interface{}),
                ResponseWriter: w,
                Request:        req,
            }
            for k, v := range params {
                c.Params[k] = v
            }
            // 요청 url에 해당하는 handler 수행
            handler(&c)
            return
        }
    }
    // 요청 URL에 해당하는 handler를 찾지 못하면 NotFound 에러 처리
    http.NotFound(w, req)
    return
}

func convParams2Int(a, b interface {}) (int, int) {

    n1, _ := strconv.Atoi(fmt.Sprintf("%v", a))
    n2, _ := strconv.Atoi(fmt.Sprintf("%v", b))

    return n1, n2
}
 
func main() {
    r := &router{make(map[string]map[string]HandlerFunc)}
 

    r.HandleFunc("GET", "/", func(c *Context) {
        fmt.Fprintln(c.ResponseWriter, "Welcome!")
    })
 
    // 더하기 핸들러
    r.HandleFunc("GET", "/plus/:number1/:number2", func(c *Context) {

        n1, n2 := convParams2Int(c.Params["number1"], c.Params["number2"])
        
        fmt.Fprintf(c.ResponseWriter, "%v + %v = %v",
            n1, n2, n1 + n2)
    })

    // power will call math.Pow(number1,number2)
    r.HandleFunc("GET", "/power/:number1/:number2", func(c *Context) {

        n1, n2 := convParams2Int(c.Params["number1"], c.Params["number2"])
        
        fmt.Fprintf(c.ResponseWriter, "%v ^ %v = %v",
            n1, n2, math.Pow(float64(n1),float64(n2)))
    })

    http.ListenAndServe(":8080", r)
}