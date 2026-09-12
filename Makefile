SHELL := /bin/sh

.PHONY: help validate-design check-json lint test test-reference build security ci

help:
	@printf '%s\n' 'Targets: validate-design check-json lint test-reference test build security ci'

validate-design:
	python3 tools/validate_design.py

check-json:
	python3 -c 'import json,pathlib; [json.loads(p.read_text(encoding="utf-8")) for p in pathlib.Path("design").rglob("*.json")]; print("JSON OK")'

lint: check-json validate-design

test-reference:
	PYTHONPATH=src python3 -m unittest discover -s tests -v

test: validate-design test-reference

build:
	python3 -m compileall -q src
	@echo 'Reference contracts compile. Production runtime build remains undefined until technology ADRs are accepted.'

security:
	@echo 'Repository security workflows remain authoritative; add runtime scanners after stack selection.'

ci: lint test build
