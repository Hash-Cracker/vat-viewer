#include <stdio.h>  // #88C0D0 for function names, #4C566A for comments
#include <stdlib.h>

void print_file(const char *filename) {  // #81A1C1 for keywords, #8FBCBB for constants
    FILE *file = fopen(filename, "r");
    if (file == NULL) {
        perror("Error opening file");
        return;
    }

    char buffer[1024];  // #A3BE8C for strings
    size_t bytesRead;

    while ((bytesRead = fread(buffer, 1, sizeof(buffer), file)) > 0) {
        fwrite(buffer, 1, bytesRead, stdout);
    }

    fclose(file);
}

int main(int argc, char *argv[]) {
    if (argc < 2) {  // #81A1C1 for keywords, #BF616A for errors
        char buffer[1024];
        while (fgets(buffer, sizeof(buffer), stdin)) {
            fputs(buffer, stdout);
        }
    } else {
        for (int i = 1; i < argc; i++) {
            print_file(argv[i]);
        }
    }

    return 0;
}

