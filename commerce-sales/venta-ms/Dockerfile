FROM node:20-alpine

WORKDIR /app
COPY package*.json ./
# CHANGED: use install instead of ci, and skip audits/fund messages
RUN npm install --omit=dev --no-audit --no-fund

COPY . .
ENV NODE_ENV=production
EXPOSE 3000
CMD ["npm", "start"]