FROM dockware/dev:latest
LABEL authors="mr. pixel KG"

# Ensure the config directory exists
RUN mkdir -p /var/www/html/config/packages

# Add the framework.yaml file with the trusted_proxies setting
# Uses proxy ip automatically. Be careful that the container is in a secure network and only accessible through proxy.
RUN echo "framework:\n    trusted_proxies: '127.0.0.1,REMOTE_ADDR'" > /var/www/html/config/packages/framework.yaml

# Set working directory
WORKDIR /var/www/html