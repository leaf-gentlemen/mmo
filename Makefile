
#工程初始化
setup:init-submodule init-toolchain

init-submodule:
	# 须先自行保存全局git账密
	@git submodule update --init --recursive

# 安装必要工具链,请仔细阅读注释自行安装部分
init-toolchain:
# 请自行下载安装 docker
	# 镜像站初始化可以用 sudo toolchain/docker/docker.sh 或根据自己需求修改执行
	#请自行下载安装 go https://go.dev/doc/install
	#请自行下载安装 protoc v3.21.8  url=https://github.com/protocolbuffers/protobuf/releases/download/v21.8/protoc-21.8-linux-x86_64.zip  then add bin to PATH
	#idea设置fileWatch: setting->tools->file watchers->import 选中 poker/toolchain/idea/filewatchers/watchers.xml,里面包含了gofmt goimports clangformat(.proto)
	#自行安装goland插件,推荐插件:.ignore ;diagrams.net ;gittoolbox ;goanno ;Conventional Commit;vscode keymap;
	#统一设置comments以空格开头 settings->Editor->Code style->Go->Other->Add a leading space to comments
	#自行设置goanno模板(参考其他同学)使符合golint
	go env -w GOPROXY=https://goproxy.cn,direct
	go env -w GOPRIVATE=steed.fun
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	#工程统一使用clang-format对.proto格式化
	sudo apt install clang-format -y
	#保存.go文件前必须用gofmt和goimports工具格式化代码
	go install golang.org/x/tools/cmd/goimports@latest
	#go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	#提交前必须用golangci-lint工具检查代码
	GOPROXY="" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
pre-commit:yaml-lint
	go mod tidy
	#fieldalignment -fix ./...
	golangci-lint run --issues-exit-code 1 -v "./..."
yaml-lint:
	 docker run --rm -v $(shell pwd):/code registry.gitlab.com/pipeline-components/yamllint yamllint .


# 通过--go_opt=M 在来替换.proto文件内的go_opt的功能,以避免go_opt写死完整包名使得其他项目无法复用,但是M参数不支持*匹配符,故需要用shell遍历.另外,要加module参数使得生成的路径扣除前缀
# 参考使用实例 https://maratrix.cn/post/2021/01/15/how-to-use-protoc-notes/
PROTO_PUBLIC_MAP=$(shell ls protocol/pubproto -l| grep .proto |awk '{print "--go_opt=Mpubproto/"$$9"=github.com/leaf-gentlemen/mmo/protos/pubproto "}')
protos-compile:
    # 公共写协议
	protoc --proto_path=protocol/  --go_opt=module=github.com/leaf-gentlemen/mmo  --go_out=.  \
	protocol/pubproto/*.proto $(PROTO_PUBLIC_MAP)
