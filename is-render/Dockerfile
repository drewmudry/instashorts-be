# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY is-render/package*.json ./
# Use npm install instead of npm ci since we don't have package-lock.json
RUN npm install

# Copy source code
COPY is-render/ ./

# Build TypeScript
RUN npm run build

# Final stage
FROM node:20-alpine

WORKDIR /app

# Install runtime dependencies
COPY is-render/package*.json ./
# Use npm install with --omit=dev instead of npm ci
RUN npm install --omit=dev

# Copy built files
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules

EXPOSE 3000

CMD ["node", "dist/index.js"]

