# This is for ARM 32 v6 (eg, rpi zero w)
FROM arm32v6/nginx:alpine

# The nginx config file: forwards /api request to backend
COPY nginx-default.conf /etc/nginx/conf.d/default.conf

# The Vue app in nginx's doc root
COPY ./dist /usr/share/nginx/html

EXPOSE 80

# Run nginx in the foreground
CMD ["nginx", "-g", "daemon off;"]