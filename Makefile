

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