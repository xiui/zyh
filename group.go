package zyh

import "net/http"

type Group struct {
	path string
	engine *Engine
	middleware HanderFunc
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
	if group.engine.middleware != nil {
		handerTree = append(handerTree, group.engine.middleware)
	}

	//再添加 group 的中间函数
	if group.middleware != nil {
		handerTree = append(handerTree, group.middleware)
	}

	for i := len(handles) - 1; i >= 0; i -- {
		handerTree = append(handerTree, handles[i])
	}

	if method == http.MethodPost {
		group.engine.postMethodTrees[path] = handerTree
	} else if method == http.MethodGet {
		group.engine.getMethodTrees[path] = handerTree
	}

}

func (group *Group) Use(middleware HanderFunc) {
	group.middleware = middleware
}