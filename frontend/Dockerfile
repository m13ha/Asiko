FROM node:18-alpine

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application
COPY . .

# Expose the port Expo uses
EXPOSE 19000
EXPOSE 19006

# Start the application
CMD ["npm", "run", "web"]