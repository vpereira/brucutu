FROM nginx:alpine

RUN apk add apache2-utils

COPY nginx.conf /etc/nginx/nginx.conf
COPY create_htpasswd.sh /create_htpasswd.sh
COPY index.html /etc/nginx/html/index.html

# Make the script executable
RUN chmod +x /create_htpasswd.sh


# Run the script to generate the .htpasswd file using environment variables
RUN /create_htpasswd.sh

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
