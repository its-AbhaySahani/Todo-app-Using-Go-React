FROM node:18 AS build
WORKDIR /app
COPY . . 
RUN npm install
RUN npm run build

FROM nginx:alpine AS production

# Copy the built React app to Nginx's web server directory
COPY --from=build /app/dist /usr/share/nginx/html

# Expose port 80 for the Nginx server
EXPOSE 80

# Start Nginx when the container runs
CMD ["nginx", "-g", "daemon off;"]