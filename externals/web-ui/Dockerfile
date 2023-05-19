FROM node:18.15.0 as build 


WORKDIR /web-ui


COPY package*.json .


RUN yarn install


COPY . .


RUN yarn run build


FROM nginx:1.19


COPY ./nginx/nginx.conf /etc/nginx/nginx.conf


COPY --from=build /web-ui/build /usr/share/nginx/html