package GeeCache

import (
	"fmt"
	"github.com/linyerun/GeeCache/cache"
	"github.com/linyerun/GeeCache/protobuf"
	"github.com/linyerun/GeeCache/utils"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpServer struct {
	ip       string
	port     uint
	basePath string
}

func (hs *HttpServer) Run() error {
	defer utils.Logger().Infof("Your Request should be [https://%v:%d/%v/<GroupName>/<key>]", hs.ip, hs.port, hs.basePath)
	return http.ListenAndServe(fmt.Sprintf("%v:%d", hs.ip, hs.port), hs)
}

func (hs *HttpServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// 记录过来的请求记录
	utils.Logger().Infof("[%s] path=%s addr=%s is comming at %v", req.Method, req.URL.Path, req.RemoteAddr, time.Now().Format("2006-01-02 15:04:05"))

	// 不以 BasePath 开头, 拒绝处理
	if !strings.HasPrefix(req.URL.Path, "/"+hs.basePath) {
		utils.Logger().Errorln("HttpServer serving unexpected path: " + req.URL.Path)
		http.Error(resp, "HttpServer serving unexpected path: "+req.URL.Path, http.StatusBadRequest)
		return
	}

	// basePath/<GroupName>/<key>
	parts := strings.SplitN(req.URL.Path[len(decodeBasePath(hs.basePath)):], "/", 2)
	if len(parts) != 2 {
		http.Error(resp, "bad request, because your request <GroupName/key> is illegal.", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group, ok := GetGroup(groupName)

	if !ok {
		http.Error(resp, "no such group that group_name = "+groupName, http.StatusNotFound)
		return
	}
	if len(key) == 0 {
		http.Error(resp, "key can not be null", http.StatusBadRequest)
		return
	}

	method := req.Method
	resp.Header().Set("Content-Type", binaryContentType)
	if method == http.MethodGet {
		httpGetFunc(resp, group, key)
		return
	} else if method == http.MethodPost {
		httpPostFunc(resp, req, group, key)
		return
	} else if method == http.MethodDelete {
		httpDeleteFunc(resp, group, key)
		return
	}

	// 未知请求方法
	utils.Logger().Errorln("HttpServer serving unexpected method: " + method)
	http.Error(resp, "HttpServer serving unexpected method: "+method, http.StatusNotFound)
}

func httpGetFunc(resp http.ResponseWriter, group IGroup, key string) {
	view, err := group.Get(key)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	// 使用Protobuf压缩一下
	data, err := proto.Marshal(&protobuf.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = resp.Write(data); err != nil {
		utils.Logger().Errorln(err.Error())
	}
}

func httpPostFunc(resp http.ResponseWriter, req *http.Request, group IGroup, key string) {
	// 获取请求体数据
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	err = group.AddOrUpdate(key, cache.NewByteView(bytes))
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		_, err := resp.Write([]byte(err.Error()))
		if err != nil {
			utils.Logger().Errorln(err.Error())
		}
		return
	}
	resp.WriteHeader(http.StatusOK)
}

func httpDeleteFunc(resp http.ResponseWriter, group IGroup, key string) {
	// 删除操作
	if err := group.Delete(key); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		_, err := resp.Write([]byte(err.Error()))
		if err != nil {
			utils.Logger().Errorln(err.Error())
		}
		return
	}
	resp.WriteHeader(http.StatusOK)
}
