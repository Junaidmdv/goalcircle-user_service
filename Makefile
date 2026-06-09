gen:
	protoc \
		--proto_path=./proto \
		--go_out=./proto/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./proto/pb \
		--go-grpc_opt=paths=source_relative \
		user_service.proto    



KEY_DIR = ./assets/secrets

genkey:
	@if [ -f $(KEY_DIR)/private.pem ]; then \
		echo "private.pem already exists. Run 'make clean' first."; \
		exit 1; \
	fi
	mkdir -p $(KEY_DIR)
	openssl genrsa -out $(KEY_DIR)/private.pem 2048
	openssl rsa -in $(KEY_DIR)/private.pem -pubout -out $(KEY_DIR)/public.pem
	chmod 600 $(KEY_DIR)/private.pem
	echo "$(KEY_DIR)/private.pem" >> .gitignore 

addkeytoenv: 
	@{ \
    echo 'JWT_PRIVATE_KEY="'$$(cat $(KEY_DIR)/private.pem)'"'; \
    echo 'JWT_PUBLIC_KEY="'$$(cat $(KEY_DIR)/public.pem)'"'; \
    } >> .env  
	@echo "Keys written to .env"
clean:
	rm -f $(KEY_DIR)/private.pem $(KEY_DIR)/public.pem  