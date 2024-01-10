#include <stdio.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <arpa/inet.h>

int get() {
    int sockfd;
    struct sockaddr_in server_addr;

    // 创建套接字
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd == -1) {
        perror("socket");
        exit(EXIT_FAILURE);
    }

    // 设置服务器地址
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(8080);  // 服务器端口号
    server_addr.sin_addr.s_addr = inet_addr("127.0.0.1");  // 服务器IP地址

    // 连接到服务器
    if (connect(sockfd, (struct sockaddr *)&server_addr, sizeof(server_addr)) == -1) {
        perror("connect");
        exit(EXIT_FAILURE);
    }

    // 将套接字转换为FILE类型
    FILE *file = fdopen(sockfd, "r+");  // 或者使用"w+"来支持读写

    // 使用FILE类型进行读写操作
    fprintf(file, "GET /albums HTTP/1.0\r\n\r\n");
    fflush(file);

    char buffer[1024];
    
    while (fgets(buffer, sizeof(buffer), file) != NULL)
    {
        printf("%s\n", buffer);
    }

    // 关闭套接字和FILE
    fclose(file);

    return 0;
}
