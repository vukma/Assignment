 FROM nginx:1.24.0
 
 RUN apt-get update && apt-get upgrade -y
 
 COPY index.html /usr/share/nginx/html
 
 EXPOSE 9000
 
 CMD ["nginx", "-g", "daemon off;"]