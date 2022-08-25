ifeq ($(OS),Windows_NT)
	BIN_EXT = .exe
endif

FS = fs$(BIN_EXT)

$(FS):
	go build -o $(FS)

sandbox:
	mkdir sandbox

sandbox/structure: sandbox
	cd sandbox && \
		mkdir structure && cd structure && \
			echo 1 > 1.txt && \
			echo 2 > 2.txt && \
			mkdir link && cd link && \
				mklink 1.txt ..\1.txt

test: $(FS) sandbox/structure
	$(FS) copy sandbox/structure sandbox/structure-copied
	$(FS) delete sandbox/structure
	$(FS) move sandbox/structure-copied sandbox/structure

clean:
	$(FS) delete sandbox
