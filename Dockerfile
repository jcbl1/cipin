FROM golang:1.20.3-bullseye
COPY . ./cipin
# RUN ["go", "env", "-w", "GOPROXY=https://goproxy.cn,direct"] # Uncomment this line if needed
RUN ["go","build","-o","/testbuild","-C","./cipin"]
FROM ubuntu:22.10
COPY --from=0 /testbuild .
CMD ["./testbuild"]
