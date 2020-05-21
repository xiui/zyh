package zyh

import (
	"net/http"
)

type Group struct {
	path string
	engine *Engine
	middlewares []HanderFunc
}

/**
	path: 接口路径 "/test"
	handles: 实现方法, 如果数量大于 1, 则会自动取消掉 Use 设置的 middleware
 */
func (group *Group) POST(path string, handles ...HanderFunc) {
	group.registerHandleTree("POST", path, handles...)
}
func (group *Group) GET(path string, handles ...HanderFunc) {
	group.registerHandleTree("GET", path, handles...)
}


func (group *Group) registerHandleTree(method string, path string, handles ...HanderFunc) {

	if len(handles) < 1 {
		panic("接口必须要有一个函数处理函数")
	}

	path = group.path + path

	var handerTree []HanderFunc

	//先添加 engine 的中间函数
	handerTree = append(handerTree, group.engine.middlewares...)

	//再添加 group 的中间函数
	handerTree = append(handerTree, group.middlewares...)


	for _, handle := range handles {
		handerTree = append(handerTree, handle)
	}

	if method == http.MethodPost {
		group.engine.postMethodTrees[path] = handerTree
	} else if method == http.MethodGet {
		group.engine.getMethodTrees[path] = handerTree
	}

}

//调用直接覆盖之前设置的, 但是调用之前的 POST,GET 等, 都使用了之前设置的
func (group *Group) UseMiddleware(middlewares ...HanderFunc) {
	group.middlewares = middlewares
}

//添加新的中间件, 之前的会保留, 只对后面发 POST, GET 方法起作用
func (group *Group) AddMiddleware(middlewares ...HanderFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}