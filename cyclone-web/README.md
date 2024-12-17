
## 安装

```sh
 $ git clone <git@..repo_name>
 $ cd repo_name
 $ yarn
```

## 开发

第一次运行需要配置 `proxy.json`，可以直接通过命令生成

```sh
 $ npm run proxy
```

proxy.json 样例，endpoint 代表需要代理的 path

```json
{
  "default": {
    "api": "http://127.0.0.1:8080",
    "endpoints": [
      "/api/*"
    ]
  }
}
```


运行

```
 $ npm run start
```

打开 http://localhost:3000 即可运行

## 编译发布

```sh
 // 生成 build.tar
 $ sh ./package.sh

 // 生产机器上获取 build.tar 然后解压放到 www 目录
 $ curl -o build.tar http://127.0.0.1:8000/build.tar
 $ tar xopf build.tar
```

编译过后的文件放置在 build 目录，为了测试编译后的结果（考虑到 proxy 的问题，否则直接启动一个静态文件服务器就能测试），可单独启动一个 node-server 来运行，也可以配置  Nginx ，将 api 访问转发给后端测试服务器。

一个 简单的 conf 为

```
 $ vim /usr/local/etc/nginx/servers/your-project-nginx-conf.conf
```

输入
your-project-nginx-conf.conf
```conf
server {
    set $api_server http://127.0.0.1:8080;
    set $build_path /Users/admin/www/idcos/react-boilerplate-master/build;
    listen       80;
    server_name  localhost;

    #charset koi8-r;

    #access_log  logs/host.access.log  main;

    location = / {
        root   html;
        index  index.html index.htm;
    }

    location / {

        root   $build_path; #config html dir

        if ( $content_type = "application/json" ) {
            proxy_pass $api_server;
        }

        if ( $uri ~ /login.html ) {
            proxy_pass $api_server;
        }

        if ( $uri ~ /api ) {
            proxy_pass $api_server;
        }

        if ( $uri ~ /static ) {
            proxy_pass $api_server;
        }

        index  index.html index.htm;
        client_max_body_size 1000m;
    }

    #error_page  404              /404.html;

    # redirect server error pages to the static page /50x.html
    #
    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
        root   html;
    }

    # proxy the PHP scripts to Apache listening on 127.0.0.1:80
    #
    #location ~ \.php$ {
    #    proxy_pass   http://127.0.0.1;
    #}

    # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
    #
    #location ~ \.php$ {
    #    root           html;
    #    fastcgi_pass   127.0.0.1:9000;
    #    fastcgi_index  index.php;
    #    fastcgi_param  SCRIPT_FILENAME  /scripts$fastcgi_script_name;
    #    include        fastcgi_params;
    #}

    # deny access to .htaccess files, if Apache's document root
    # concurs with nginx's one
    #
    #location ~ /\.ht {
    #    deny  all;
    #}
}
```

然后修改 root 为项目 build 目录地址。

执行

```sh
 $ sudo nginx -s stop && sudo nginx
```

这时候可直接访问: http://localhost







