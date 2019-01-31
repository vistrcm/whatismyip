.PHONY: deploy

deploy:
	gcloud functions deploy WhatIsMyIp --runtime go111 --trigger-http
