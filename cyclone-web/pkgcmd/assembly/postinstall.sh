#if [ `readlink -f /etc/nginx/conf.d/default.conf` != "@APP_DIR@/conf/nginx.conf" ]; then
#    mv /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.bak
#    ln -s @APP_DIR@/conf/nginx.conf /etc/nginx/conf.d/default.conf
#fi
#(service nginx status || service nginx start) >/dev/null 2>&1 ||:
