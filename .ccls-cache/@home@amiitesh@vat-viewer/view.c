#include <stdio.h>
#include <stdlib.h>

void print_file(const char *filename) {
    FILE *file = fopen(filename, "r");
    if (file == NULL) {
        perror("Error opening file");
        return;
    }

    char buffer[1024];
    size_t bytesRead;

    while ((bytesRead = fread(buffer, 1, sizeof(buffer), file)) > 0) {
        fwrite(buffer, 1, bytesRead, stdout);
    }

    fclose(file);
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        // No files provided; read from standard input
        char buffer[1024];
        while (fgets(buffer, sizeof(buffer), stdin)) {
            fputs(buffer, stdout);
        }
    } else {
        // Loop through all provided files
        for (int i = 1; i < argc; i++) {
            print_file(argv[i]);
        }
    }

    return 0;
}

