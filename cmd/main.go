package main

import (
	"fmt"
	"net/http"
)

func main() {
	frame := &Frame{}
	frame.AddFilter(GetMethodFilter)
	frame.AddRoute("/hg", GetEcho)

	frame = &Frame{}
	frame.AddFilter(PostMethodFilter)
	frame.AddRoute("/hh", GetEcho)
	http.ListenAndServe(":8000", nil)

	fmt.Println("Hello from go")
}

func GetEcho(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Request was success"))
}

func GetMethodFilter(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if "GET" == r.Method {
			f.ServeHTTP(w, r)
		} else {
			w.Write([]byte("Method not mapped"))
			w.WriteHeader(405)
		}
	}
}

func PostMethodFilter(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if "POST" == r.Method {
			f.ServeHTTP(w, r)
		} else {
			w.Write([]byte("Method not mapped"))
			w.WriteHeader(405)
		}
	}
}

type Filter func(f http.HandlerFunc) http.HandlerFunc

type Frame struct {
	filters []Filter
}

func (frame *Frame) AddFilter(filter Filter) {
	frame.filters = append(frame.filters, filter)
}

func (frame *Frame) AddRoute(pattern string, f http.HandlerFunc) {
	frame.process(pattern, f, len(frame.filters)-1)
}

func (frame *Frame) process(pattern string, f http.HandlerFunc, index int) {
	if index == -1 {
		http.Handle(pattern, f)
		return
	}
	fwrap := frame.filters[index](f)
	index--
	frame.process(pattern, fwrap, index)
}

/*

type MethodToHandle map[string]func(http.Header, io.ReadCloser) (int, string)

type RequestMapper struct {
	handler *map[string]*MethodToHandle
}
func GetHH(headers http.Header, body io.ReadCloser) (int, string) {
	sb := strings.Builder{}
	sb.WriteString("Request Body was : \n")
	by, err := ioutil.ReadAll(body)
	if err == nil {
		for _, b := range by {
			sb.WriteByte(b)
		}
	} else {
		fmt.Println(fmt.Errorf("Error Occurred while read %e", err))
	}
	fmt.Printf("Called : %s", body)
	fmt.Println()
	return 500, sb.String()
}

func (d *RequestMapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if h, ok := (*d.handler)[r.URL.Path]; ok {
		// fmt.Printf("Endpoint URL %s\n", h)
		if m, ok := (*h)[r.Method]; ok {
			status, response := m(r.Header, r.Body)
			w.WriteHeader(status)
			if status != 200 {
				w.Write([]byte("<HTML><H1 align='center'>Server Error has occurred</H1></HTML>"))
			} else {
				w.Write([]byte(response))
			}
		}
	}
}

func NewRequestMapper() RequestMapper {
	outer_map := make(map[string]*MethodToHandle)
	return RequestMapper{
		handler: &outer_map,
	}
}

func (d *RequestMapper) MapRequestURL(url string, method string, handler func(http.Header, io.ReadCloser) (int, string)) {
	mapping, ok := (*d.handler)[url]
	if ok {
		(*mapping)[method] = handler
	} else {
		inner_map := make(MethodToHandle)
		inner_map[method] = handler
		(*d.handler)[url] = &inner_map
	}
} */
