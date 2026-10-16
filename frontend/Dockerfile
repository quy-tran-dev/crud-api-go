# ----------------------------------------------------
# Stage 1: Build the React Application
# ----------------------------------------------------
FROM node:22-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Copy toàn bộ mã nguồn
COPY . .

# Build production bundle
RUN npm run build

# ----------------------------------------------------
# Stage 2: Serve with Nginx Alpine
# ----------------------------------------------------
FROM nginx:alpine

# Copy cấu hình nginx tùy chỉnh
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy production build từ stage 1
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
