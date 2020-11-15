package zyh

import (
	"net/http"
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

		ctx := createContext(w, r, methodTree)
		if ctx != nil {
			ctx.Next()
		}
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
func (engine *Engine) UseMiddleware(middlewares ...HanderFunc) {
	engine.middlewares = middlewares
}

//添加新的中间件, 之前的会保留, 只对后面发 POST, GET 方法起作用
func (engine *Engine) AddMiddleware(middlewares ...HanderFunc) {
	engine.middlewares = append(engine.middlewares, middlewares...)
}




func (engine *Engine) Run(port string) error {

	writeLog("HTTP Start")
	err := http.ListenAndServe(port, engine)

	return err
}