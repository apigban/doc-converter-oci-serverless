# ---- Final Stage ----
FROM nginx:1.25-alpine

# Copy the static HTML file into the Nginx public directory
COPY ./frontend/index.html /usr/share/nginx/html/index.html

# Expose port 80 for the web server
EXPOSE 80

# The default Nginx entrypoint will start the server
CMD ["nginx", "-g", "daemon off;"]