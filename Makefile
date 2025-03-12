.PHONY: check fix

div = $(shell printf '=%.0s' {1..120})

DIR="."
check:
	@echo ${div}
	uv run --active ruff check $(DIR)
	uv run --active ruff format $(DIR) --check
	@echo "Done!"

fix:
	@echo ${div}
	uv run --active ruff format $(DIR) 
	@echo ${div}
	uv run --active ruff check $(DIR) --fix
	@echo "Done!"