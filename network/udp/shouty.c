#include <stdio.h>
#include <ctype.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>

#define MAX_BUFFER_LEN 1024
#define SERVER_PORT 8888

int main() {
    // create socket connection
    int fd, recv_bytes, i;
    struct sockaddr_in server_addr, client_addr; 
    socklen_t socklen;
    char *buffer;

    buffer = malloc(MAX_BUFFER_LEN);
    if ((fd = socket(AF_INET, SOCK_DGRAM, 0)) < 0) {
        perror("socket creation failed\n");
        return -1;
    }
    memset(&server_addr, 0, sizeof(struct sockaddr_in));
    memset(&client_addr, 0, sizeof(struct sockaddr_in));

    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(SERVER_PORT);

    socklen = sizeof(client_addr);
    if (bind(fd, (const struct sockaddr *)&server_addr, sizeof(server_addr)) < 0) {
        perror("socket bind failed\n");
        return -1;
    }

    while (1) {
        // recv data from client
        recv_bytes = recvfrom(fd, buffer, MAX_BUFFER_LEN, 0,
                        (struct sockaddr *)&client_addr, &socklen);
        if (recv_bytes < 0) {
            perror("error while reciving from the socket");
        } else {
            buffer[recv_bytes] = '\0';
            for (i = 0; i < recv_bytes; i++) {
                buffer[i] = toupper(buffer[i]);
            }
            // send the same data back to the client
            sendto(fd, (const char *)buffer, strlen(buffer), 0, 
                (const struct sockaddr *) &client_addr, socklen);
        }
    }
    return 0;
}