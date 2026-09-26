# Design-validation image only.
# This is NOT the future game client/server runtime image.
FROM python:3.14-slim

WORKDIR /workspace

COPY design ./design
COPY tools ./tools

CMD ["python3", "tools/validate_design.py"]
