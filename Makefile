
CC= gcc
AR= ar rcu
RM= rm -f
ALL_O= *.o
ALL_A= *.a lib/*.a
CFLAGS= -DHAS_SOCKLEN_T
C_INCLUDES= -Isrc/

all:
	$(CC) $(CFLAGS) $(C_INCLUDES) -g -fPIC -Isrc/ -c src/*.c 
	mkdir -p ./lib
	$(AR) lib/libenet.a *.o
	$(RM) $(ALL_O)
	go install
	