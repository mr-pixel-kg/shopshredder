# Step 1: Build stage
FROM node:22-alpine as builder

WORKDIR /app

# Copy package.json and yarn.lock
COPY package.json yarn.lock ./

# Install dependencies
RUN yarn install

# Copy application
COPY . .

# Set build-time environment variable for backend URL
ARG VITE_BACKEND_URL
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL

# Build frontend for production
RUN yarn build

# Step 2: Production stage
FROM nginx:alpine

# Copy build directory from Vite to Nginx directory
COPY --from=builder /app/dist /usr/share/nginx/html

# Export Nginx server port
EXPOSE 80

# Configure health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD curl -f http://localhost || exit 1

# Run nginx server
CMD ["nginx", "-g", "daemon off;"]