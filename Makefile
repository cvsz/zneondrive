SHELL := /bin/bash

CONTROL := bash tools/zneondrive-control.sh
CLIENT_PACKAGE ?=
SERVER_PACKAGE ?=
SERVICE ?=
CONFIRM_RESET ?=
LOG_TAIL ?= 200

.PHONY: help doctor env-init deps-server deps-client control-panel \
	validate-design validate-docs validate-runtime validate-control check-json lint \
	test-reference go-test go-integration go-vet go-build test build security ci \
	server-install server-up server-down server-restart server-status server-health server-logs server-reset \
	runtime-up runtime-down db-shell redis-cli \
	client-generate client-build client-package-linux editor-build client-install client-play client-doctor \
	game-server-build game-server-package-linux package-all-linux game-server-install game-server-start game-server-stop game-server-status game-server-logs \
	full-install full-up full-down status clean

help:
	@printf '%s\n' \
	  'PROJECT: NEON DRIVE — Make control surface' \
	  '' \
	  'Bootstrap / control' \
	  '  make doctor                 Check host/runtime prerequisites' \
	  '  make env-init               Create .env + unique local server key' \
	  '  make deps-server            Install base Ubuntu/Debian server dependencies' \
	  '  make deps-client            Install base Linux client/dev dependencies (not UE itself)' \
	  '  make control-panel          Open interactive operator menu' \
	  '' \
	  'Service plane' \
	  '  make server-install         Pull/build/start PostgreSQL + Redis + Go API' \
	  '  make server-up|server-down|server-restart' \
	  '  make server-status|server-health' \
	  '  make server-logs [SERVICE=game-api]' \
	  '  make db-shell | make redis-cli' \
	  '  CONFIRM_RESET=YES make server-reset' \
	  '' \
	  'Player client' \
	  '  make client-install CLIENT_PACKAGE=/path/client.zip' \
	  '  make client-build           Build NeonDriveClient with UE_ROOT on Linux' \
  '  make client-package-linux   Cook/stage/archive a Linux player package' \
	  '  make client-play            Launch installed client; ZNEON_GAME_API_URL may override API' \
	  '  make client-doctor          Windows/Linux client guidance and host check' \
	  '' \
	  'Dedicated gameplay server' \
	  '  make game-server-package-linux  Cook/stage/archive Linux dedicated server' \
  '  make package-all-linux           Package Linux client + dedicated server' \
  '  make game-server-install SERVER_PACKAGE=/path/server.tar.gz' \
	  '  make game-server-start|game-server-stop|game-server-status|game-server-logs' \
	  '' \
	  'Full stack' \
	  '  make full-install CLIENT_PACKAGE=... SERVER_PACKAGE=...' \
	  '  make full-up | make full-down | make status' \
	  '' \
	  'Quality' \
	  '  make validate-design validate-docs validate-runtime validate-control' \
	  '  make go-test go-integration go-vet go-build test build security ci'

doctor:
	@$(CONTROL) doctor

env-init:
	@$(CONTROL) env-init

deps-server:
	@$(CONTROL) deps-server

deps-client:
	@$(CONTROL) deps-client

control-panel:
	@$(CONTROL) control-panel

validate-design:
	python3 tools/validate_design.py

validate-docs:
	python3 tools/validate_docs.py

validate-runtime:
	python3 tools/validate_runtime_v0_4.py
	python3 tools/validate_runtime_v0_5.py
	python3 tools/validate_runtime_v0_6.py

validate-control:
	python3 tools/validate_control_surface.py
	bash -n tools/zneondrive-control.sh

check-json:
	python3 -c 'import json,pathlib; [json.loads(p.read_text(encoding="utf-8")) for p in pathlib.Path("design").rglob("*.json")]; json.loads(pathlib.Path("game/NeonDrive.uproject").read_text(encoding="utf-8")); print("JSON OK")'

lint: check-json validate-design validate-docs validate-runtime validate-control

test-reference:
	PYTHONPATH=src python3 -m unittest discover -s tests -v

go-test:
	cd services/game-api && go test ./...

go-integration:
	@test -n "$${TEST_DATABASE_URL:-}" || { echo 'TEST_DATABASE_URL is required'; exit 2; }
	cd services/game-api && go test -tags=integration ./...

go-vet:
	cd services/game-api && go vet ./...

go-build:
	cd services/game-api && go build -trimpath ./cmd/server

test: lint test-reference go-test

build: go-build
	python3 -m compileall -q src
	@echo 'Reference + Go service build complete. Unreal builds require UE_ROOT / the self-hosted UE 5.8 runner.'

security:
	@echo 'Run repository CodeQL/Dependency Review workflows for authoritative scanner evidence.'
	@echo 'Runtime transport, anti-cheat, HA/DR and production security gates remain evidence-gated.'

server-install:
	@$(CONTROL) server-install

server-up:
	@$(CONTROL) server-up

server-down:
	@$(CONTROL) server-down

server-restart:
	@$(CONTROL) server-restart

server-status:
	@$(CONTROL) server-status

server-health:
	@$(CONTROL) server-health

server-logs:
	@LOG_TAIL=$(LOG_TAIL) $(CONTROL) server-logs $(SERVICE)

server-reset:
	@CONFIRM_RESET=$(CONFIRM_RESET) $(CONTROL) server-reset

runtime-up: server-up

runtime-down: server-down

db-shell:
	@$(CONTROL) db-shell

redis-cli:
	@$(CONTROL) redis-cli

client-generate:
	@$(CONTROL) client-generate

client-build:
	@$(CONTROL) client-build

client-package-linux:
	@$(CONTROL) client-package-linux

editor-build:
	@$(CONTROL) editor-build

client-install:
	@CLIENT_PACKAGE="$(CLIENT_PACKAGE)" $(CONTROL) client-install "$(CLIENT_PACKAGE)"

client-play:
	@$(CONTROL) client-play

client-doctor:
ifeq ($(OS),Windows_NT)
	powershell -ExecutionPolicy Bypass -File tools/install-client.ps1 -Mode doctor
else
	@$(CONTROL) doctor
endif

game-server-build:
	@$(CONTROL) game-server-build

game-server-package-linux:
	@$(CONTROL) game-server-package-linux

package-all-linux:
	@$(CONTROL) package-all-linux

game-server-install:
	@SERVER_PACKAGE="$(SERVER_PACKAGE)" $(CONTROL) game-server-install "$(SERVER_PACKAGE)"

game-server-start:
	@$(CONTROL) game-server-start

game-server-stop:
	@$(CONTROL) game-server-stop

game-server-status:
	@$(CONTROL) game-server-status

game-server-logs:
	@LOG_TAIL=$(LOG_TAIL) $(CONTROL) game-server-logs

full-install:
	@CLIENT_PACKAGE="$(CLIENT_PACKAGE)" SERVER_PACKAGE="$(SERVER_PACKAGE)" $(CONTROL) full-install "$(CLIENT_PACKAGE)"

full-up:
	@$(CONTROL) full-up

full-down:
	@$(CONTROL) full-down

status:
	@$(CONTROL) status

clean:
	rm -rf .runtime
	@echo 'Removed local control-panel runtime metadata only. Compose volumes and packaged builds were preserved.'

ci: lint test-reference go-test go-vet go-build
