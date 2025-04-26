# Build application
FROM maven:3.9.9-sapmachine-21 AS builder

WORKDIR /app
COPY . .

RUN mvn clean package -DskipTests

# Run application
#FROM openjdk:21-jdk-slim
#
#WORKDIR /app
#COPY --from=builder /app/target/*.jar app.jar
#COPY *.mv.db .
#
#EXPOSE 8080
#
#ENTRYPOINT ["java", "-jar", "app.jar"]

FROM alpine:3.20
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk --update --no-cache add openjdk21-jre
WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar
COPY *.mv.db .

EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]

# 使用构建指令
# docker build -t springboot-backend:v0.1 .
# docker run -itd -p8081:8080 springboot-backend:v0.1

