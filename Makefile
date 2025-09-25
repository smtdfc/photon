
DATE := $(shell date +%Y%m%d)
HASH := $(shell git rev-parse --short HEAD)

ALL_MODULES := $(shell grep 'directory' go.work | awk '{print $$2}')

VER ?= v0.0.0-$(DATE)$(HASH)
MODULE ?= all
PUSH ?= no  

.PHONY: tag
tag:
	@modules=$$( [ "$(MODULE)" = "all" ] && echo "$(ALL_MODULES)" || echo "$(MODULE)" ); \
	for mod in $$modules; do \
		if [ -d "$$mod" ]; then \
			cd $$mod; \
			if [ -n "$$(git status --porcelain)" ]; then \
				echo "Committing changes in $$mod"; \
				git add .; \
				git commit -m "chore: commit before tagging $(VER)"; \
			fi; \
			echo "Tagging module $$mod with $(VER)"; \
			git tag -f $(VER); \
			if [ "$(PUSH)" = "yes" ]; then \
				git push origin HEAD --tags; \
			fi; \
			cd - > /dev/null; \
		else \
			echo "Module '$$mod' do not exist !"; \
		fi; \
	done
	@echo "Tagging done."
