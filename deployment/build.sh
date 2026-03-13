#!/bin/bash

# set default environment (if no argument passed, default to pre, support prod, pre)
MODE="pre"

# set default version (if no argument passed, default to latest)
VERSION="latest"

MAKE_MODE="debug"
IMAGE="ai-btc-contract"
# parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --mode=*)
            MODE="${1#*=}"
            shift
            ;;
        --version=*)
            VERSION="${1#*=}"
            shift
            ;;
        *)
            echo "❌ Unknown option: $1"
            echo "Usage: ./build.sh [--mode=prod|pre] [--version=xxx]"
            exit 1
            ;;
    esac
done

# check arguments
if [ "$MODE" != "prod" ] && [ "$MODE" != "pre" ]; then
    echo "❌  invalid argument: $MODE 'prod' or 'pre'"
    exit 1
fi

if [ "$MODE" == "prod" ]; then
    MAKE_MODE=""
fi

# ensure buildx is usable
docker buildx inspect --bootstrap > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "⚙️  Creating buildx builder..."
    docker buildx create --use
fi

# build image

IMAGE_NAME="${IMAGE}:${VERSION}"
echo "🔨 start building $IMAGE_NAME..."
# --platform linux/amd64 指定构建目标平台（Linux x86_64） 用于交叉编译或在 macOS M 芯片上构建 x86 镜像
# --network=host  # 允许容器访问主机网络
# --build-arg BUILD_MODE=$MODE 设置构建参数
# --tag $IMAGE_NAME 设置镜像标签
# --file ./Dockerfile 构建文件
# ./ 构建上下文  回到 project/ 根目录 才能找到对应文件
# --load     # 将构建好的镜像 直接加载到本地 Docker 引擎 中

if ! docker buildx build \
    --platform linux/amd64 \
    --network=host \
    --no-cache \
    -f ./Dockerfile \
    --build-arg BUILD_MODE=$MAKE_MODE \
    -t $IMAGE_NAME \
    ../\
    --load; then
    echo "❌ image build failed"
    exit 1
fi
echo "✅ image built successfully: $MAKE_MODE"

# login image registry
echo "🔐 login image registry..."
if ! docker login -u chenyvhang2@gmail.com -p gycyh1106; then
    echo "❌ login failed"
    exit 1
fi

echo "📤 start pushing $IMAGE_NAME..."
docker tag $IMAGE:$VERSION nebverifier/$IMAGE:$VERSION
if ! docker push nebverifier/$IMAGE_NAME; then
    echo "❌ image push failed"
    docker logout
    exit 1
fi

# build complete
echo "✅ images built and pushed successfully: $IMAGE_NAME"