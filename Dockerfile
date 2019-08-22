FROM alpine:latest AS build

RUN apk add --update sudo
RUN adduser -G abuild -g "Alpine Package Builder" -s /bin/sh -D builder \
  && echo "builder ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

USER builder
WORKDIR /home/builder
RUN sudo apk add --update alpine-sdk \
  && abuild-keygen -i -a \
  && git clone --depth 1 git://dev.alpinelinux.org/aports \
  && cd ~/aports/community/tor \
  && sed -i 's/--enable-transparent/--enable-transparent --enable-tor2web-mode/g' APKBUILD \
  # Disable tests
  && sed -i 's/make test//g' APKBUILD \
  && abuild verify && abuild -r \
  && cd ~/packages/community/x86_64 \
  && rm -Rf ~/aports \
  && sudo apk del alpine-sdk

FROM alpine
RUN mkdir /packages
COPY --from=build /home/builder/packages/community/x86_64/ /packages
RUN apk add --no-cache libevent bash
RUN apk add --allow-untrusted /packages/tor-0.*.apk && rm -Rf /packages
# ADD scallion /
ADD torrc /
# ENTRYPOINT ["/investigo"]
CMD ["/bin/bash"]