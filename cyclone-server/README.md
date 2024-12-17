# 目录
- [Compile & Running](#compile-running)
- [Go code style guide](#go-code-style-guide)
- [REST API guide](#rest-api-guide)
- [Git commit message guide](#git-commit-message-guide)

# Compile & Running
- 安装go、git(略)

- 下载源码

	```shell
	$ cd /home/voidint
	$ git clone -b webank git@gitlab.idcos.com:CloudBoot/cloudboot.git cloudboot-webank
	```

- 将项目根目录纳入GOPATH环境变量

	```shell
	$ export GOPATH="$GOPATH:/home/voidint/workspace/cloudboot-webank"
	```

- 下载gbb(可选)

	```shell
	$ go get -u -v github.com/voidint/gbb
	```

- 使用gbb编译所有可执行文件

	```shell
	$ cd /home/voidint/workspace/cloudboot-webank
	$ gbb -a --debug
	```

- 或者使用go原生工具链编译可执行文件

	```shell
	$ cd /home/voidint/workspace/cloudboot-webank/src/idcos.io/cloudboot/cmd/cloudboot-server
	$ go install
	```

- 准备MySQL数据库环境
    - 创建数据库

    ```shell
    mysql> CREATE DATABASE IF NOT EXISTS `cloudboot_webank` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
    ```
    - 导入工程内的全量SQL文件

    ```shell
    $ mysql -h 127.0.0.1 -P 3307 -u root -p cloudboot_webank < ./doc/db/webank.sql
    ```

- 启动cloudboot-server

	```shell
	$ cloudboot-server -c /path/to/cloudboot-server.conf
	```

# Go code style guide
戳[这里](https://github.com/voidint/go-style-guide)

# REST API guide
- URL中集合资源的操作使用名词复数。

	Don't
	
	```
	// 查询所有用户
	GET /examples/user 
	GET /examples/user/list
	
	// 查询所有用户的所有图书
	GET /examples/user/book
	GET /examples/user/book/list
	```
	
	Do
	
	```
	// 查询所有用户
	GET/examples/users
	
	// 查询所有用户的所有图书
	GET /examples/users/books 
	```

- URL中对单个资源的操作可以使用单数。

	Don't
	
	```
	// 查询id为101的用户的学校信息
	GET /examples/users/101/schools
	```
	
	Do
	
	```
	// 查询id为101的用户的学校信息(假定某一时刻所在的学校只有一个)
	GET /examples/users/101/school
	```

- URL中资源的命名需要由多个单词构成时，使用中划线(`-`)连接。

	Don't
	
	```
	// 查询id为101的用户的用户名
	GET /examples/users/101/user_name
	GET /examples/users/101/userName
	```
	
	Do
	
	```
	// 查询id为101的用户的用户名
	GET /examples/users/101/user-name
	```


- URL中不要出现动词，使用`HTTP Methods`表示对URL资源的操作动作。

	Don't
	
	```
	// 更新用户信息
	POST /examples/users/update 
	POST /examples/updateUser
	```
	
	Do
	
	```
	// 更新用户信息
	PUT /examples/users
	```
- 对于资源集合中的某一具体资源的操作，建议将具体资源的标识符放入URL中。

	Don't
	
	```
	// 查询id为101的用户信息
	GET /examples/users?id=101
	
	// 查询id为101的用户的所有图书
	GET /examples/users/books?id=101
	```

	Do
	
	```
	// 查询id为101的用户信息
	GET /examples/users/101
	
	// 查询id为101的用户的所有图书
	GET /examples/users/101/books
	```
	
- URL中的每一级遵循资源范围上`从大到小`、`从广到窄`的排列原则，URL的末尾是当前操作的目标。
	
	Don't
	
	```
	// 查询id为101的用户的所有图书的所有作者
	GET /examples/books/authors/users/101 // 这个查询操作的目标应该是图书作者，而非用户。
	```
	
	Do
	
	```	
	// 查询id为101的用户的所有图书的所有作者
	// 资源范围从广到窄应该是"所有用户的集合"-->"某个用户"-->"某个用户拥有的所有图书"-->"某个用户拥有图书的作者集合"，而当前操作的目标应该是"作者"这一资源
	GET /examples/users/101/books/authors
	```

- `Request`和`Response`字段一律小写(若由多个单词组成，则用下划线(`_`)连接多个单词)。

	Don't
	
	```
	GET /examples/users?UserName=voidint 
	```

	Do
	
	```
	GET /examples/users?user_name=voidint
	```
	
- 操作失败情况下，不能使用`200`作为`HTTP Status code`，使用合适的[HTTP Status Codes](https://httpstatuses.com/)。

	Don't
	
	```
	// 查询系统中并不存在的id为101的用户
	GET /examples/users/101 
	Response HTTP Status Code: 200
	Response Body: {"status": "success", "message": "user(id=101) not found"}
	```

	Do
	
	```
	// 查询系统中并不存在的id为101的用户
	GET /examples/users/101 
	Response HTTP Status Code:404
	Response Body: {"status": "success", "message": "user(id=101) not found"}
	```

- MIME类型为`Application/json`的`Response Body`要求遵循以下结构

	```json
	{
		"status": "success/failure",
		"message": "balabala...",
		"content": {
			
		}
	}
	```

# Git commit message guide
### 前缀
- `Add feature`: 新增新特性/功能
- `Mofidy feature`: 修改原有特性/功能（功能/特性上发生了变化）
- `Fixbug`: bug修复
- `Add lib`: 新增第三方库
- `Upgrade lib`: 升级第三方库
- `Remove lib`: 移除第三方库
- `Add doc`: 添加文档
- `Modify doc`: 修改文档
- `Add comment`: 添加注释
- `Code refactor`: 代码重构（功能/特性保持不变，代码层面的重构）
- `Code optimization`: 代码优化
- `Add ut`: 添加单元测试
- `Fix ut`: 修复单元测试
- `Modify SQL`: 修改SQL
- `Fix conflict`: 冲突修复
- `Remove file`: 移除文件/目录
- ......


