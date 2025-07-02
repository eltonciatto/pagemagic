FROM openresty/openresty:alpine

# Install additional Lua modules
RUN /usr/local/openresty/luajit/bin/luarocks install lua-cjson
RUN /usr/local/openresty/luajit/bin/luarocks install lua-resty-jwt
RUN /usr/local/openresty/luajit/bin/luarocks install lua-resty-http

# Copy Nginx configuration
COPY nginx/nginx.conf /etc/nginx/nginx.conf

# Copy Lua scripts
COPY lua/ /etc/nginx/lua/

# Create log directory
RUN mkdir -p /var/log/nginx

# Create non-root user
RUN adduser -D -s /bin/sh -u 1000 nginx-proxy

# Set permissions
RUN chown -R nginx-proxy:nginx-proxy /var/log/nginx /etc/nginx
RUN chmod +x /etc/nginx/lua/*.lua

USER nginx-proxy

EXPOSE 80

CMD ["/usr/local/openresty/bin/openresty", "-g", "daemon off;"]
