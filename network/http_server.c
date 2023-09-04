#include <stdio.h>
#include <ctype.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <unistd.h>


#define SERVER_PORT 9999
#define MAX_CONCURRENT_CLIENT 10
#define MAX_BUFFER_LEN 10000


char *convert_to_json(char *buffer) {
    printf("buffer received from browser: %s", buffer);
    return buffer;
}
void echoClient(int clientfd) {
    char *buffer, *json_buffer;
    buffer = malloc(MAX_BUFFER_LEN);
    int recv_bytes;
    // take http request headers and convert it to json
    while (1) {
        memset(buffer, 0, MAX_BUFFER_LEN);
        recv_bytes = recv(clientfd, buffer, MAX_BUFFER_LEN, 0);
        if (recv_bytes == 0) {
            break;
        }
        buffer[recv_bytes] = '\0';
        json_buffer = convert_to_json(buffer);
        send(clientfd, json_buffer, strlen(json_buffer), 0);
    }
    free(buffer);
    close(clientfd);
}
int main() {
    int fd, clientfd, recv_bytes;
    socklen_t socklen;
    char *buffer;
   
    struct sockaddr_in server_addr, client_addr;

    fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) {
        perror("socket creation failed!");
        return -1;
    }
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(SERVER_PORT);
    server_addr.sin_family = AF_INET;
    socklen = sizeof(server_addr);
    if (bind(fd, (const struct sockaddr *)&server_addr, socklen) < 0) {
        perror("socket bind failed!");
        return -1;
    }
    if (listen(fd, 5) < 0) {
        perror("socket listen failed!");
    }
    while (1) {
        memset(&client_addr, 0, sizeof(client_addr));
        clientfd = accept(fd, (struct sockaddr *) &client_addr, &socklen);
        echoClient(clientfd); 
    }
    close(fd);
}