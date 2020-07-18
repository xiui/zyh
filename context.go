package zyh

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
)

type File struct {
	File multipart.File
	FileHeader *multipart.FileHeader
}

type Context struct {

	w http.ResponseWriter
	r *http.Request
	handlers []HanderFunc
	currentMethodIndex int

	Params map[string]string

}

func (ctx *Context) ValueWithDefault(key string, defaultValue string) string {

	val := ctx.r.FormValue(key)

	if len(val) > 0{
		return val
	}

	return defaultValue

}

func (ctx *Context) Value(key string) string {

	val := ctx.r.FormValue(key)

	return val
}

func (ctx *Context) ValueInt(key string) (int, error) {

	val := ctx.r.FormValue(key)

	valInt, err := strconv.ParseInt(val, 10, 64)

	return int(valInt), err
}

func (ctx *Context) ValueUint(key string) (int, error) {

	val := ctx.r.FormValue(key)

	valInt, err := strconv.ParseUint(val, 10, 64)

	return int(valInt), err
}

func (ctx *Context) ValueFloat(key string) (float64, error) {

	val := ctx.r.FormValue(key)

	return strconv.ParseFloat(val, 64)

}

func (ctx *Context) ValueBool(key string) (bool, error) {

	val := ctx.r.FormValue(key)

	return strconv.ParseBool(val)

}

func (ctx *Context) FileValues() ([]File, error) {

	var newFiles []File

	err := ctx.r.ParseMultipartForm(32 << 20)

	if err != nil {
		return newFiles, err
	}

	files := ctx.r.MultipartForm.File

	for k,_ := range files {

		file, head, err := ctx.r.FormFile(k)
		defer file.Close()

		if err != nil {
			return newFiles, err
		}

		newFiles = append(newFiles, File{File:file, FileHeader:head})

	}

	return newFiles, nil

}

//判断参数中是否有某个key
func (ctx *Context) HasParamsKey(key string) bool {
	_, has := ctx.Params[key]
	return has
}

//判断一个参数是否是空的, 没有此字段/空字符串/0
func (ctx *Context) IsEmpty(key string) bool {

	if !ctx.HasParamsKey(key) {
		return true
	}
	v := ctx.Value(key)
	if v == "" || v == "0" {
		return true
	}
	return false
}

func (ctx *Context) Next() {


	//因为ctx.currentMethodIndex 初始化的时候设置的是 -1
	//每次 if 中会先 +1, 所以这里判断条件的时候用了 len(ctx.handlers) - 1, 防止超出数组
	if ctx.currentMethodIndex < len(ctx.handlers) - 1 {

		//加一, 方法树会继续调用, 调用时走下一个方法
		ctx.currentMethodIndex ++

		ctx.handlers[ctx.currentMethodIndex](ctx)

	} else {
		//TODO: 这里需要考虑怎么处理, 方法超出了的问题
		writeLog("函数调用次数超出 handlers 中设置的方法, 请查看 ctx.Next() 调用情况, 如果是逻辑中最后一个方法, 不需要调用 ctx.Next()方法")
	}

}

//返回 json 数据
func (ctx *Context) JSON(code int, data interface{}) {

	b, e := json.Marshal(data)

	var err error
	if e != nil {

		writeLog("context.go->JSON, json 格式化错误: " + e.Error())

		ctx.w.WriteHeader(500)

		_, err = ctx.w.Write([]byte("bad to JSON string"))

	} else {
		ctx.w.WriteHeader(code)
		_, err = ctx.w.Write(b)

	}

	if err != nil {
		writeLog("context.go->JSON, 返回失败: " + e.Error())
	}

}

//返回普通文本
func (ctx *Context) String(code int, data string) {
	ctx.w.WriteHeader(code)
	_, err := ctx.w.Write([]byte(data))
	if err != nil {
		writeLog("context.go->String, 返回失败: " + err.Error())
	}
}

func (ctx *Context) Request() *http.Request {
	return ctx.r
}

func (ctx *Context) ResponseWriter() http.ResponseWriter {
	return ctx.w
}
