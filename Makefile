NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

all: format lint test install

test:
	@echo -e "$(OK_COLOR)==> Testing$(NO_COLOR)"
	./tools/test.sh

format:
	@echo -e "$(OK_COLOR)==> Formating$(NO_COLOR)"
	./tools/format.sh

lint:
	@echo -e "$(OK_COLOR)==> Linting$(NO_COLOR)"
	./tools/lint.sh

install:
	@echo -e "$(OK_COLOR)==> Installing$(NO_COLOR)"
	./tools/install.sh
