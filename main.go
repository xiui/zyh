package zyh

import (
	"net/http"
	"strings"
)

func Default() *Engine {
	return &Engine{postMethodTrees: map[string][]HanderFunc{}, getMethodTrees: map[string][]HanderFunc{}}
}

type Engine struct {

	middlewares []HanderFunc

	postMethodTrees map[string][]HanderFunc
	getMethodTrees map[string][]HanderFunc

}

type HanderFunc func(ctx *Context)


func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var methodTree []HanderFunc
	hadMethodTree := false	//

	if r.Method == http.MethodPost {
		methodTree, hadMethodTree = engine.postMethodTrees[r.URL.Path]
	} else if r.Method == http.MethodGet {
		methodTree, hadMethodTree = engine.getMethodTrees[r.URL.Path]
	}

	if hadMethodTree {

		var err error

		contentType := r.Header["Content-Type"]
		isUploadFile := false

		for _, v := range contentType {

			//判断是不是传文件
			if strings.Index(v, "multipart/form-data") >= 0 {
				isUploadFile = true
				break
			}
		}

		if isUploadFile {
			err = r.ParseMultipartForm(32 << 20)
		} else {
			err = r.ParseForm()
		}


		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("params analysis is wrong"))
			return
		}

		params := map[string]string{}
		for k, value := range r.Form {
			if len(value) > 0 {
				//重名的参数只取第一个, 尽量不要使用同名参数
				params[k] = value[0]
			}
		}

		ctx := &Context{
			handlers: methodTree,
			r: r,
			w: w,
			currentMethodIndex: -1,	//调用方法树时, 会提前加一, 所以在这里设置了-1
			Params: params,
		}

		ctx.Next()
	} else {
		w.WriteHeader(404)
		w.Write([]byte("not this method"))

	}


}

func (engine *Engine) Group(path string) Group {
	return Group{
		path:path,
		engine: engine,
	}
}

/**
	path: 接口路径 "/test"
	handles: 实现方法, 如果数量大于 1, 则会自动取消掉 Use 设置的 middleware
 */
func (engine *Engine) POST(path string, handles ...HanderFunc) {
	engine.registerHandleTree("POST", path, handles...)
}
func (engine *Engine) GET(path string, handles ...HanderFunc) {
	engine.registerHandleTree("GET", path, handles...)
}


func (engine *Engine) registerHandleTree(method string, path string, handles ...HanderFunc) {

	if len(handles) < 1 {
		panic("接口必须要有一个函数处理函数")
	}

	var handerTree []HanderFunc

	//添加中间函数
	handerTree = append(handerTree, engine.middlewares...)


	for i := len(handles) - 1; i >= 0; i -- {
		handerTree = append(handerTree, handles[i])
	}

	if method == http.MethodPost {
		engine.postMethodTrees[path] = handerTree
	} else if method == http.MethodGet {
		engine.getMethodTrees[path] = handerTree
	}

}

//调用直接覆盖之前设置的, 但是调用之前的 POST,GET 等, 都使用了之前设置的
func (engine *Engine) Use(middlewares ...HanderFunc) {
	engine.middlewares = middlewares
}


func (engine *Engine) Run(port string) error {

	writeLog("HTTP Start")
	err := http.ListenAndServe(port, engine)

	return err
}