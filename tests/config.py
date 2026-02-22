"""测试配置常量"""

import os

# 服务地址，可通过环境变量覆盖
BASE_URL = os.getenv("FORMA_BASE_URL", "http://localhost:8888/api")

# 请求超时（秒）
TIMEOUT = int(os.getenv("FORMA_TIMEOUT", "10"))

# 响应码
CODE_SUCCESS = "200"
CODE_INVALID_PARAM = "10001"
CODE_NOT_FOUND = "10002"
CODE_UNAUTHORIZED = "10003"
CODE_FORBIDDEN = "10004"
CODE_INTERNAL = "99999"
