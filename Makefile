MAKE=make
DESTINATION=dist
SERVER_SRC=src/server
CLIENT_SRC=src/client
CLEAN=rm -rf

all: clean client server

folder:
	mkdir -p $(DESTINATION)

server: folder
	cd $(SERVER_SRC) && make
	cp $(SERVER_SRC)/server $(DESTINATION)/.
	cp $(SERVER_SRC)/config.prod.json $(DESTINATION)/config.json

client: 
	cd $(CLIENT_SRC) && gulp

clean:
	$(CLEAN) $(DESTINATION)