// Package page_items contains parsed result by PageProcesser.
// The result is processed by Pipeline.
package page_items

import (
    "Spider/src/core/common/request"
)

//PageItems表示PageProcesser解析的实体保存结果，最终将被输出。
type PageItems struct {

    // req是包含解析结果的请求对象，它保存在PageItems中。
    req *request.Request

    // items是解析结果的容器。
    items map[string]string

    // skip表示是否将结果发送到调度程序。
    skip bool
}

// NewPageItems返回初始化的PageItems对象。
func NewPageItems(req *request.Request) *PageItems {
    items := make(map[string]string)
    return &PageItems{req: req, items: items, skip: false}
}

// GetRequest返回PageItems的请求。
func (this *PageItems) GetRequest() *request.Request {
    return this.req
}

//AddItem将结果保存到PageItems中。
func (this *PageItems) AddItem(key string, item string) {
    this.items[key] = item
}

// GetItem返回键的值。
func (this *PageItems) GetItem(key string) (string, bool) {
    t, ok := this.items[key]
    return t, ok
}

// GetAll返回所有键结果。
func (this *PageItems) GetAll() map[string]string {
    return this.items
}

// SetSkip set skip true使此页不被管道处理。
func (this *PageItems) SetSkip(skip bool) *PageItems {
    this.skip = skip
    return this
}

// GetSkip返回skip标签。
func (this *PageItems) GetSkip() bool {
    return this.skip
}
