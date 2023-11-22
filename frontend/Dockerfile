FROM node:lts as builder
ARG PRODUCTION
COPY public /app/public
COPY src /app/src
COPY package.json /app/
COPY package-lock.json /app/
ENV NODE_ENV=production
ENV REACT_APP_API_PREFIX=${PRODUCTION:+https://raw.ninja:8443}
ENV REACT_APP_API_PREFIX=${REACT_APP_API_PREFIX:-http://localhost:8080}
RUN echo "Building fronted with backend URL: $REACT_APP_API_PREFIX"
WORKDIR /app
RUN npm ci --force
RUN npm run build

FROM nginx:1.25.3-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]