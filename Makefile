SHELL := /bin/sh

.PHONY: help validate-design validate-docs validate-runtime check-json lint test test-reference build security ci runtime-up runtime-down go-test

help:
	@printf '%s\n' 'Targets: validate-design validate-docs validate-runtime check-json lint test-reference go-test test build security ci runtime-up runtime-down'

validate-design:
	python3 tools/validate_design.py

validate-docs:
	python3 tools/validate_docs.py

validate-runtime:
	python3 tools/validate_runtime_v0_4.py
	python3 tools/validate_runtime_v0_5.py
	python3 tools/validate_runtime_v0_6.py

check-json:
	python3 -c 'import json,pathlib; [json.loads(p.read_text(encoding="utf-8")) for p in pathlib.Path("design").rglob("*.json")]; print("JSON OK")'

lint: check-json validate-design validate-docs validate-runtime

test-reference:
	PYTHONPATH=src python3 -m unittest discover -s tests -v

go-test:
	cd services/game-api && go test ./...

test: validate-design validate-runtime test-reference

build:
	python3 -m compileall -q src
	@echo 'Reference contracts compile. Unreal source build requires the self-hosted unreal-5.8 runner.'

security:
	@echo 'Repository security workflows remain authoritative; runtime threat/anti-cheat gates are still open.'

runtime-up:
	docker compose up --build

runtime-down:
	docker compose down

ci: lint test build
