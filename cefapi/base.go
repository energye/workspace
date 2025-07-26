package cefapi

/*
#cgo CFLAGS: -I./include
#cgo windows LDFLAGS: -L./lib -lcef -lwinmm -lversion
#cgo linux LDFLAGS: -L./lib -lcef -ldl -lm
#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit -framework CoreFoundation -framework CoreGraphics

#include "include/capi/cef_base_capi.h"
#include <stdlib.h>
#include <string.h>

// 定义函数指针类型（因为CGO不能直接使用cef_base_add_ref_t等类型）
typedef void (CEF_CALLBACK *add_ref_func)(struct _cef_base_ref_counted_t*);
typedef int (CEF_CALLBACK *release_func)(struct _cef_base_ref_counted_t*);
typedef int (CEF_CALLBACK *has_one_ref_func)(struct _cef_base_ref_counted_t*);
typedef int (CEF_CALLBACK *has_at_least_one_ref_func)(struct _cef_base_ref_counted_t*);
typedef void (CEF_CALLBACK *del_func)(struct _cef_base_scoped_t*);

// 声明C回调函数原型
void go_cef_base_add_ref(struct _cef_base_ref_counted_t* self);
int go_cef_base_release(struct _cef_base_ref_counted_t* self);
int go_cef_base_has_one_ref(struct _cef_base_ref_counted_t* self);
int go_cef_base_has_at_least_one_ref(struct _cef_base_ref_counted_t* self);
void go_cef_base_scoped_del(struct _cef_base_scoped_t* self);
*/
import "C"
import (
	"log"
	"sync"
	"unsafe"
)

// 引用计数对象类型
type RefCountedType int

const (
	OwnedByGo  RefCountedType = iota // Go管理生命周期
	OwnedByCEF                       // CEF管理生命周期
)

// GoRefCounted 表示Go端的引用计数对象
type GoRefCounted struct {
	ptr         unsafe.Pointer // 指向整个内存块的指针
	refCount    int            // 当前引用计数
	mutex       sync.Mutex     // 保护引用计数
	managedBy   RefCountedType // 生命周期管理方
	userDataPtr unsafe.Pointer // 用户数据指针
	userDataLen uintptr        // 用户数据长度
}

// GoBaseScoped 表示Go端的范围对象
type GoBaseScoped struct {
	ptr unsafe.Pointer // 指向cef_base_scoped_t结构
}

// CreateRefCounted 创建引用计数对象
// userDataSize: 用户数据区域大小
// managedBy: 生命周期管理方
func CreateRefCounted(userDataSize uintptr, managedBy RefCountedType) *GoRefCounted {
	// 计算总大小 = 指针大小 + cef_base_ref_counted_t大小 + 用户数据大小
	totalSize := unsafe.Sizeof(uintptr(0)) + C.sizeof_cef_base_ref_counted_t + userDataSize

	// 分配内存
	cMem := C.calloc(1, C.size_t(totalSize))
	if cMem == nil {
		log.Fatal("内存分配失败")
	}

	// 创建Go包装器
	grc := &GoRefCounted{
		ptr:         cMem,
		refCount:    1, // 初始引用计数为1
		managedBy:   managedBy,
		userDataPtr: nil,
		userDataLen: userDataSize,
	}

	// 存储指向Go对象的指针
	*(*uintptr)(cMem) = uintptr(unsafe.Pointer(grc))

	// 计算cef_base_ref_counted_t指针位置
	cefBasePtr := unsafe.Add(cMem, unsafe.Sizeof(uintptr(0)))

	// 初始化cef_base_ref_counted_t结构
	base := (*C.cef_base_ref_counted_t)(cefBasePtr)
	base.size = C.size_t(totalSize)

	// 设置虚函数表（使用重新定义的类型）
	base.add_ref = (C.add_ref_func)(C.go_cef_base_add_ref)
	base.release = (C.release_func)(C.go_cef_base_release)
	base.has_one_ref = (C.has_one_ref_func)(C.go_cef_base_has_one_ref)
	base.has_at_least_one_ref = (C.has_at_least_one_ref_func)(C.go_cef_base_has_at_least_one_ref)

	// 设置用户数据区域指针
	if userDataSize > 0 {
		userDataStart := unsafe.Add(cefBasePtr, int(C.sizeof_cef_base_ref_counted_t))
		grc.userDataPtr = userDataStart
		C.memset(userDataStart, 0, C.size_t(userDataSize))
	}

	return grc
}

func BaseRefCountedSize() uintptr {
	return C.sizeof_cef_base_ref_counted_t
}

// CreateBaseScoped 创建范围对象
func CreateBaseScoped() *GoBaseScoped {
	// 分配内存
	cMem := C.calloc(1, C.sizeof_cef_base_scoped_t)
	if cMem == nil {
		log.Fatal("内存分配失败")
	}

	// 创建Go包装器
	scoped := &GoBaseScoped{ptr: cMem}

	// 初始化cef_base_scoped_t结构
	base := (*C.cef_base_scoped_t)(cMem)
	base.size = C.sizeof_cef_base_scoped_t
	base.del = (C.del_func)(C.go_cef_base_scoped_del)

	return scoped
}

// GetBaseRefCounted 获取cef_base_ref_counted_t指针
func (g *GoRefCounted) GetBaseRefCounted() *C.cef_base_ref_counted_t {
	if g.ptr == nil {
		return nil
	}
	return (*C.cef_base_ref_counted_t)(unsafe.Add(g.ptr, unsafe.Sizeof(uintptr(0))))
}

// GetUserData 获取用户数据区域
func (g *GoRefCounted) GetUserData() unsafe.Pointer {
	return g.userDataPtr
}

// AddRef 增加引用计数
func (g *GoRefCounted) AddRef() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.refCount++
	log.Printf("增加引用: 新计数=%d", g.refCount)
}

// Release 减少引用计数
func (g *GoRefCounted) Release() bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.refCount--
	log.Printf("减少引用: 新计数=%d", g.refCount)

	if g.refCount <= 0 {
		log.Printf("引用计数归零，释放资源")
		g.Free()
		return true
	}
	return false
}

// Free 释放分配的内存
func (g *GoRefCounted) Free() {
	if g.ptr != nil {
		C.free(g.ptr)
		g.ptr = nil
		g.userDataPtr = nil
	}
}

// FreeScoped 释放范围对象
func (g *GoBaseScoped) FreeScoped() {
	if g.ptr != nil {
		C.free(g.ptr)
		g.ptr = nil
	}
}

// 下面是导出的C回调函数实现

//export go_cef_base_add_ref
func go_cef_base_add_ref(self *C.cef_base_ref_counted_t) {
	// 获取关联的Go对象
	grc := getGoRefCountedFromBase(self)
	if grc != nil {
		grc.AddRef()
	}
}

//export go_cef_base_release
func go_cef_base_release(self *C.cef_base_ref_counted_t) C.int {
	grc := getGoRefCountedFromBase(self)
	if grc != nil {
		if grc.Release() {
			return 1 // 返回true表示引用计数归零
		}
		return 0
	}
	return 1 // 没有关联对象，返回true表示可以删除
}

//export go_cef_base_has_one_ref
func go_cef_base_has_one_ref(self *C.cef_base_ref_counted_t) C.int {
	grc := getGoRefCountedFromBase(self)
	if grc != nil {
		grc.mutex.Lock()
		defer grc.mutex.Unlock()
		if grc.refCount == 1 {
			return 1
		}
	}
	return 0
}

//export go_cef_base_has_at_least_one_ref
func go_cef_base_has_at_least_one_ref(self *C.cef_base_ref_counted_t) C.int {
	grc := getGoRefCountedFromBase(self)
	if grc != nil {
		grc.mutex.Lock()
		defer grc.mutex.Unlock()
		if grc.refCount >= 1 {
			return 1
		}
	}
	return 0
}

//export go_cef_base_scoped_del
func go_cef_base_scoped_del(self *C.cef_base_scoped_t) {
	// 对于范围对象，直接释放内存
	C.free(unsafe.Pointer(self))
}

// getGoRefCountedFromBase 从cef_base_ref_counted_t指针获取Go对象
func getGoRefCountedFromBase(base *C.cef_base_ref_counted_t) *GoRefCounted {
	if base == nil {
		return nil
	}

	// 计算存储Go指针的位置：base指针向前退一个指针大小的位置
	ptrAddr := uintptr(unsafe.Pointer(base)) - unsafe.Sizeof(uintptr(0))
	ptr := *(*uintptr)(unsafe.Pointer(ptrAddr))

	return (*GoRefCounted)(unsafe.Pointer(ptr))
}

func main() {
	// 示例1：创建由Go管理的对象
	goManagedObj := CreateRefCounted(128, OwnedByGo)
	defer func() {
		log.Println("释放Go管理对象")
		goManagedObj.Free()
	}()

	// 设置用户数据
	if dataPtr := goManagedObj.GetUserData(); dataPtr != nil {
		data := []byte("Go管理的数据")
		C.memcpy(dataPtr, unsafe.Pointer(&data[0]), C.size_t(len(data)))
		log.Printf("用户数据: %s", C.GoString((*C.char)(dataPtr)))
	}

	// 增加引用计数
	goManagedObj.AddRef()
	goManagedObj.Release()

	// 示例2：创建由CEF管理的对象
	cefManagedObj := CreateRefCounted(0, OwnedByCEF)

	// 模拟CEF增加引用
	C.go_cef_base_add_ref(cefManagedObj.GetBaseRefCounted())

	// 模拟CEF释放引用
	if C.go_cef_base_release(cefManagedObj.GetBaseRefCounted()) == 1 {
		log.Println("CEF管理对象引用计数归零")
	}

	// 示例3：创建范围对象
	scopedObj := CreateBaseScoped()
	defer scopedObj.FreeScoped()

	log.Println("所有对象创建完成")
}
