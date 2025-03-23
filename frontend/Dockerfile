FROM node:18
WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy source code and expose Vite dev server
COPY . .
EXPOSE 5173

# Start Vite dev server
CMD ["npm", "run", "dev", "--", "--host"]
